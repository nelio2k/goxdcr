package conflictlog

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/couchbase/goxdcr/base"
	"github.com/stretchr/testify/require"
)

func TestRules_Parse(t *testing.T) {
	testData := []struct {
		name            string
		jsonStr         string
		expectedMapping map[base.CollectionNamespace]Target
		shouldFail      bool
	}{
		{
			name:       "[negative] empty json",
			jsonStr:    `{}`,
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
				"collection": 10,
				"loggingRules": []
			}`,
			shouldFail: true,
		},
		{
			name: "[negative] empty source key",
			jsonStr: `{
				"bucket":"B1",
				"collection": 10,
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
			name: "[negative] target incomplete",
			jsonStr: `{
				"bucket":"B1",
				"collection": "S1.C1",
				"loggingRules": {
					"US.Ohio": {
						"bucket": "B2"
					}
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
			name: "[negative] invalid collection value",
			jsonStr: `{
				"bucket":"B1",
				"collection": "S1.C1",
				"loggingRules": {
					"US.Ohio": {
						"bucket": "B1",
						"collection": "S2"
					}
				}
			}`,
			shouldFail: true,
		},
		{
			name: "[positive] only default target",
			jsonStr: `{
				"bucket":"B1",
				"collection":"S1.C1"
			}`,
		},
		{
			name: "[positive] only default target but loggingRules=nil",
			jsonStr: `{
				"bucket":"B1",
				"collection":"S1.C1",
				"loggingRules": null
			}`,
		},
		{
			name: "[positive] only default target but loggingRules={}",
			jsonStr: `{
				"bucket":"B1",
				"collection":"S1.C1",
				"loggingRules": {}
			}`,
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
				{ScopeName: "US"}: NewTarget("B1", "S1", "C1"),
			},
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
				{ScopeName: "US", CollectionName: ""}:     NewTarget("B1", "S1", "C1"),
				{ScopeName: "India", CollectionName: ""}:  NewTarget("B1", "S1", "C1"),
			},
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

			for source, expTarget := range tt.expectedMapping {
				target, ok := rules.Mapping[source]

				if !ok {
					require.Fail(t, "source not found in mapping", source.String(), buf.String())
				} else {
					require.Equal(t, expTarget, target)
				}
			}
		})
	}

}
