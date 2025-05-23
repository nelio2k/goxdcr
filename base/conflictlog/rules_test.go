/*
Copyright 2024-Present Couchbase, Inc.

Use of this software is governed by the Business Source License included in
the file licenses/BSL-Couchbase.txt.  As of the Change Date specified in that
file, in accordance with the Business Source License, use of this software will
be governed by the Apache License, Version 2.0, included in the file
licenses/APL2.txt.
*/

package conflictlog

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/couchbase/goxdcr/v8/base"
	"github.com/stretchr/testify/require"
)

func TestRules_Parse(t *testing.T) {
	testData := []struct {
		name string
		// jsonStr is the json in string form which simulates the input
		// to the update settings. The test setup will parse this first.
		// and any error in parsing is a test fail as it is not the objective
		// of the test
		jsonStr string

		// expectedMapping is the what the final mapping should look like
		// after parsing. A nil is an accepted value.
		expectedMapping map[base.CollectionNamespace]Target
		// expectedTarget is the expected fallback target value.
		expectedTarget Target

		// shouldFail=true implies we expect a failure for the input
		shouldFail bool
	}{
		{
			name:       "[negative] empty json",
			jsonStr:    `{}`,
			shouldFail: true,
		},
		{
			name: "[negative] invalid bucket value type",
			jsonStr: `{
				"bucket": 10,
				"collection": "S1.C1"
			}`,
			shouldFail: true,
		},
		{
			name: "[negative] invalid collection value type",
			jsonStr: `{
				"bucket":"B1",
				"collection": 10
			}`,
			shouldFail: true,
		},
		{
			name: "[negative] invalid loggingRules value type",
			jsonStr: `{
				"bucket":"B1",
				"collection": "S1.C1",
				"loggingRules": []
			}`,
			shouldFail: true,
		},
		{
			name: "[negative] empty source key",
			jsonStr: `{
				"bucket":"B1",
				"collection": "S1.C1",
				"loggingRules": {
					"": null
				}
			}`,
			shouldFail: true,
		},
		{
			name: "[negative] source key ending with dot",
			jsonStr: `{
				"bucket":"B1",
				"collection": "S1.C1",
				"loggingRules": {
					"US.": null
				}
			}`,
			shouldFail: true,
		},
		{
			name: "[negative] source only scope but has space at the end",
			jsonStr: `{
				"bucket":"B1",
				"collection": "S1.C1",
				"loggingRules": {
					"US ": null
				}
			}`,
			shouldFail: true,
		},
		{
			name: "[negative] source scope and collection but has space at the end",
			jsonStr: `{
				"bucket":"B1",
				"collection": "S1.C1",
				"loggingRules": {
					"US.Ohio ": null
				}
			}`,
			shouldFail: true,
		},
		{
			name: "[negative] invalid type for target values",
			jsonStr: `{
				"bucket":"B1",
				"collection": "S1.C1",
				"loggingRules": {
					"US.Ohio": {
						"bucket": 999,
						"collection": 1111
					}
				}
			}`,
			shouldFail: true,
		},
		{
			name: "[positive] scope and collection incomplete, should default to _default",
			jsonStr: `{
				"bucket":"B1",
				"collection": "S1.C1",
				"loggingRules": {
					"US.Ohio": {
						"bucket": "B2"
					}
				}
			}`,
			expectedMapping: map[base.CollectionNamespace]Target{
				{ScopeName: "US", CollectionName: "Ohio"}: NewTarget("B2", "_default", "_default"),
			},
			expectedTarget: NewTarget("B1", "S1", "C1"),
		},
		{
			name: "[negative] just collection incomplete, incomplete target",
			jsonStr: `{
				"bucket":"B1",
				"collection": "S1.C1",
				"loggingRules": {
					"US.Ohio": {
						"bucket": "B2",
						"collection": "S2"
					}
				}
			}`,
			shouldFail: true,
		},
		{
			name: "[negative] fallback target collection missing, incomplete target",
			jsonStr: `{
				"bucket":"B1",
				"collection": "S1",
				"loggingRules": {
					"US.Ohio": {
						"bucket": "B2",
						"collection": "S2.C1"
					}
				}
			}`,
			shouldFail: true,
		},
		{
			name: "[positive] fallback target scope and collection missing, should default to _default._default",
			jsonStr: `{
				"bucket":"B1",
				"loggingRules": {
					"US.Ohio": {
						"bucket": "B2",
						"collection": "S2.C1"
					}
				}
			}`,
			expectedMapping: map[base.CollectionNamespace]Target{
				{ScopeName: "US", CollectionName: "Ohio"}: NewTarget("B2", "S2", "C1"),
			},
			expectedTarget: NewTarget("B1", "_default", "_default"),
		},
		{
			name: "[positive] only default target",
			jsonStr: `{
				"bucket":"B1",
				"collection":"S1.C1"
			}`,
			expectedTarget: NewTarget("B1", "S1", "C1"),
		},
		{
			name: "[positive] only default target but loggingRules=nil",
			jsonStr: `{
					"bucket":"B1",
					"collection":"S1.C1",
					"loggingRules": null
				}`,
			expectedTarget: NewTarget("B1", "S1", "C1"),
		},
		{
			name: "[positive] only default target but loggingRules={}",
			jsonStr: `{
					"bucket":"B1",
					"collection":"S1.C1",
					"loggingRules": {}
				}`,
			expectedTarget: NewTarget("B1", "S1", "C1"),
		},
		{
			name: "[positive] only scope in source key",
			jsonStr: `{
					"bucket":"B1",
					"collection": "S1.C1",
					"loggingRules": {
						"US": null
					}
				}`,
			expectedMapping: map[base.CollectionNamespace]Target{
				{ScopeName: "US"}: BlacklistTarget(),
			},
			expectedTarget: NewTarget("B1", "S1", "C1"),
		},
		{
			name: "[positive] scope.collection in source key",
			jsonStr: `{
					"bucket":"B1",
					"collection": "S1.C1",
					"loggingRules": {
						"US": null,
						"India": {},
						"US.Ohio": {
							"bucket": "B2",
							"collection": "S2.C2"
						}
					}
				}`,
			expectedMapping: map[base.CollectionNamespace]Target{
				{ScopeName: "US", CollectionName: "Ohio"}: NewTarget("B2", "S2", "C2"),
				{ScopeName: "US", CollectionName: ""}:     BlacklistTarget(),
				{ScopeName: "India", CollectionName: ""}:  NewTarget("B1", "S1", "C1"),
			},
			expectedTarget: NewTarget("B1", "S1", "C1"),
		},
		{
			name: "[negative] target collection missing, incomplete target",
			jsonStr: `{
					"bucket":"B1",
					"collection": "S1",
					"loggingRules": {
						"US": null,
						"India": {},
						"US.Ohio": {
							"bucket": "B2",
							"collection": "S2"
						}
					}
				}`,
			shouldFail: true,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			j := map[string]interface{}{}
			err := json.Unmarshal([]byte(tt.jsonStr), &j)
			require.Nil(t, err)

			rules, err := ParseRules(j)
			if tt.shouldFail {
				require.NotNil(t, err)
				return
			} else {
				require.Nil(t, err)
			}

			err = rules.Validate()
			require.Nil(t, err)

			var buf strings.Builder
			for source, target := range rules.Mapping {
				buf.WriteString(fmt.Sprintf("%s.%s", source.ScopeName, source.CollectionName))
				buf.WriteString(fmt.Sprintf("=>%s, ", target.String()))
			}

			require.Equal(t, tt.expectedTarget.Bucket, rules.Target.Bucket)
			require.Equal(t, tt.expectedTarget.NS.ScopeName, rules.Target.NS.ScopeName)
			require.Equal(t, tt.expectedTarget.NS.CollectionName, rules.Target.NS.CollectionName)

			for source, expTarget := range tt.expectedMapping {
				target, ok := rules.Mapping[source]

				if !ok {
					require.FailNow(t, "source not found in mapping", source.String(), buf.String())
				} else {
					require.Equal(t, expTarget, target)
				}
			}
		})
	}
}
