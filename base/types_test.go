/*
Copyright 2020-Present Couchbase, Inc.

Use of this software is governed by the Business Source License included in
the file licenses/BSL-Couchbase.txt.  As of the Change Date specified in that
file, in accordance with the Business Source License, use of this software will
be governed by the Apache License, Version 2.0, included in the file
licenses/APL2.txt.
*/

package base

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUleb128EncoderDecoder(t *testing.T) {
	fmt.Println("============== Test case start: TestUleb128EncoderDecoder =================")
	assert := assert.New(t)

	seed := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(seed)

	for i := 0; i < 50; i++ {
		input := generator.Uint32()
		testLeb, _, err := NewUleb128(input, nil, true)
		assert.Nil(err)

		verifyOutput := testLeb.ToUint32()
		assert.Equal(verifyOutput, input)
	}

	// Direct mem mapping test - for reading key with embedded CID
	var testByteSlice []byte = make([]byte, 1, 1)
	testByteSlice[0] = 0x09
	var testOut uint32 = Uleb128(testByteSlice).ToUint32()
	assert.Equal(uint32(9), testOut)

	fmt.Println("============== Test case end: TestUleb128EncoderDecoder =================")
}

func TestCollectionNamespaceFromString(t *testing.T) {
	fmt.Println("============== Test case start: TestCollectionNamespaceFromString =================")
	defer fmt.Println("============== Test case end: TestCollectionNamespaceFromString =================")

	assert := assert.New(t)
	namespace, err := NewCollectionNamespaceFromString("a123.123b")
	assert.Nil(err)
	assert.Equal("a123", namespace.ScopeName)
	assert.Equal("123b", namespace.CollectionName)

	_, err = NewCollectionNamespaceFromString("abcdef")
	assert.NotNil(err)
}

func TestExplicitMappingValidatorParseRule(t *testing.T) {
	fmt.Println("============== Test case start: TestExplicitMappingValidatorParseRule =================")
	defer fmt.Println("============== Test case end: TestExplicitMappingValidatorParseRule =================")
	assert := assert.New(t)

	validator := NewExplicitMappingValidator()

	key := "Scope"
	value := "Scope"
	assert.Equal(explicitRuleScopeToScope, validator.parseRule(key, value))
	assert.Equal(explicitRuleScopeToScope, validator.parseRule(key, nil))

	key = "Scope.collection"
	value = "scope2.collection2"
	assert.Equal(explicitRuleOneToOne, validator.parseRule(key, value))
	assert.Equal(explicitRuleOneToOne, validator.parseRule(key, nil))

	// Invalid names
	key = "#%(@&#FJ"
	value = "scope"
	assert.Equal(explicitRuleInvalidScopeName, validator.parseRule(key, value))

	// Too Long
	var longStrArr []string
	for i := 0; i < 256; i++ {
		longStrArr = append(longStrArr, "a")
	}
	longStr := strings.Join(longStrArr, "")
	key = longStr
	assert.Equal(explicitRuleStringTooLong, validator.parseRule(longStr, nil))

	value = longStr
	assert.Equal(explicitRuleStringTooLong, validator.parseRule(longStr, nil))

	key = "abc"
	assert.Equal(explicitRuleStringTooLong, validator.parseRule(longStr, nil))
}

func TestExplicitMappingValidatorRules(t *testing.T) {
	fmt.Println("============== Test case start: TestExplicitMappingValidatorRules =================")
	defer fmt.Println("============== Test case end: TestExplicitMappingValidatorRules =================")
	assert := assert.New(t)

	validator := NewExplicitMappingValidator()
	// We only disallow _system scope
	key := "_invalidScopeName"
	value := "validTargetScopeName"
	assert.Nil(validator.ValidateKV(key, value))

	key = "validScopeName"
	value = "%invalidScopeName"
	assert.NotNil(validator.ValidateKV(key, value))

	// Positive test cases
	key = "Scope"
	value = "TargetScope"
	assert.Nil(validator.ValidateKV(key, value))

	key = "Scope2"
	value = "TargetScope2"
	assert.Nil(validator.ValidateKV(key, value))

	// 1-N is disallowed
	key = "Scope"
	value = "TargetScope2"
	assert.NotNil(validator.ValidateKV(key, value))

	// N-1 is not allowed
	key = "Scope2Test"
	value = "TargetScope2"
	assert.NotNil(validator.ValidateKV(key, value))

	key = "AnotherScope.AnotherCollection"
	value = "AnotherTargetScope.AnotherTargetCollection"
	assert.Nil(validator.ValidateKV(key, value))

	key = "AnotherScope2.AnotherCollection2"
	value = "AnotherTargetScope2.AnotherTargetCollection2"
	assert.Nil(validator.ValidateKV(key, value))

	// 1-N is not allowed
	key = "AnotherScope.AnotherCollection"
	value = "AnotherTargetScope.AnotherTargetCollection2"
	assert.NotNil(validator.ValidateKV(key, value))

	// N-1 is not allowed
	key = "AnotherScope2.AnotherCollection3"
	value = "AnotherTargetScope2.AnotherTargetCollection2"
	assert.NotNil(validator.ValidateKV(key, value))

	// Adding non-duplicating blacklist rules
	key = "Scope3"
	assert.Nil(validator.ValidateKV(key, nil))

	key = "Scope.Collection"
	assert.Nil(validator.ValidateKV(key, nil))

	// Adding duplicating blacklist rules
	key = "Scope3.Collection3"
	assert.NotNil(validator.ValidateKV(key, nil))

	key = "Scope"
	assert.NotNil(validator.ValidateKV(key, nil))

	// Test complex mapping - one specific collection will have special mapping, everything else implicit under scope
	// 1. ScopeRedundant.ColRedundant -> ScopeTRedundant.ColTRedundant
	// 2. ScopeRedundant -> ScopeTRedundant
	key = "ScopeRedundant"
	value = "ScopeTRedundant"
	assert.Nil(validator.ValidateKV(key, value))

	// This is not redundant
	key = "ScopeRedundant.ColRedundant"
	value = "ScopeTRedundant.ColTRedundant"
	assert.Nil(validator.ValidateKV(key, value))

	// This is redundant
	key = "ScopeRedundant.ColRedundant"
	value = "ScopeTRedundant.ColRedundant"
	assert.NotNil(validator.ValidateKV(key, value))

	// ** Converse of above
	// Not Redundant
	key = "ScopeRedundant2.ColRedundant2"
	value = "ScopeTRedundant2.ColTRedundant2"
	assert.Nil(validator.ValidateKV(key, value))

	// Combining both should be fine
	key = "ScopeRedundant2"
	value = "ScopeTRedundant2"
	assert.Nil(validator.ValidateKV(key, value))

	// System scope cannot be mapped
	key = "_system"
	value = "scope"
	assert.Equal(ErrorSystemScopeMapped, validator.ValidateKV(key, value))

	// Cannot map to system scope
	key = "scope"
	value = "_system"
	assert.Equal(ErrorSystemScopeMapped, validator.ValidateKV(key, value))

	// System collection cannot be mapped
	key = "_system._mobile"
	value = "scope.collection"
	assert.Equal(ErrorSystemScopeMapped, validator.ValidateKV(key, value))

	// Cannot map to a system collection
	key = "scope.collection"
	value = "_system._mobile"
	assert.Equal(ErrorSystemScopeMapped, validator.ValidateKV(key, value))
}

func TestWrappedFlags(t *testing.T) {
	assert := assert.New(t)
	fmt.Println("============== Test case start: TestWrappedFlags =================")
	defer fmt.Println("============== Test case end: TestWrappedFlags =================")

	wrappedUpr := WrappedUprEvent{}
	assert.False(wrappedUpr.Flags.CollectionDNE())

	wrappedUpr.Flags.SetCollectionDNE()
	assert.True(wrappedUpr.Flags.CollectionDNE())
}

func TestDefaultNs(t *testing.T) {
	assert := assert.New(t)
	fmt.Println("============== Test case start: TestDefaultNS =================")
	defer fmt.Println("============== Test case end: TestDefaultNS =================")

	validator := NewExplicitMappingValidator()
	assert.Equal(explicitRuleScopeToScope, validator.parseRule("_default", "_default"))
	assert.Equal(explicitRuleScopeToScope, validator.parseRule("_default", nil))
	assert.Equal(explicitRuleOneToOne, validator.parseRule("_default._default", nil))
	assert.Equal(explicitRuleOneToOne, validator.parseRule("_default._default", "_default._default"))
	validator = NewExplicitMappingValidator()
	assert.Equal(explicitRuleOneToOne, validator.parseRule("_default.testCol", nil))
	assert.Equal(explicitRuleOneToOne, validator.parseRule("_default.testCol", "_default._default"))
	assert.Equal(explicitRuleOneToOne, validator.parseRule("testScope.testCol", nil))
	assert.Equal(explicitRuleOneToOne, validator.parseRule("testScope.testCol", "_default.nonDefCol"))

}

func TestValidRemoteClusterName(t *testing.T) {
	assert := assert.New(t)
	fmt.Println("============== Test case start: TestValidRemoteClusterName =================")
	defer fmt.Println("============== Test case end: TestValidRemoteClusterName =================")

	errMap := make(ErrorMap)

	ValidateRemoteClusterName("abc.be_fd.com", errMap)
	assert.Equal(0, len(errMap))

	ValidateRemoteClusterName("abc", errMap)
	assert.Equal(0, len(errMap))

	ValidateRemoteClusterName("12.23.34.45", errMap)
	assert.Equal(0, len(errMap))

	ValidateRemoteClusterName("endwithaPeriod.com.", errMap)
	assert.NotEqual(0, len(errMap))
}

func TestValidRules(t *testing.T) {
	assert := assert.New(t)
	fmt.Println("============== Test case start: TestValidRules =================")
	defer fmt.Println("============== Test case end: TestValidRules =================")

	validator := NewExplicitMappingValidator()
	assert.Nil(validator.ValidateKV("scope1", DefaultScopeCollectionName))
	assert.Nil(validator.ValidateKV("scope1.collection1", "scope1.collection1"))
	assert.Nil(validator.ValidateKV("scope1.collection2", "scope1.collection2"))

	validator = NewExplicitMappingValidator()
	assert.Nil(validator.ValidateKV("scope1.collection1", "scope1.collection1"))
	assert.Nil(validator.ValidateKV("scope1.collection2", "scope1.collection2"))
	assert.Nil(validator.ValidateKV("scope1", nil))

	validator = NewExplicitMappingValidator()
	assert.Nil(validator.ValidateKV("scope1", nil))
	assert.Nil(validator.ValidateKV("scope1.collection1", "scope1.collection1"))
	assert.Nil(validator.ValidateKV("scope1.collection2", "scope1.collection2"))
}

func TestRulesRedundancyCheck(t *testing.T) {
	assert := assert.New(t)
	fmt.Println("============== Test case start: TestRulesRedundancyCheck =================")
	defer fmt.Println("============== Test case end: TestRulesRedundancyCheck =================")

	validator := NewExplicitMappingValidator()
	assert.Nil(validator.ValidateKV("testScope.testCol", "testScope2.testCol"))
	assert.NotNil(validator.ValidateKV("testScope", "testScope2"))

	validator = NewExplicitMappingValidator()
	assert.Nil(validator.ValidateKV("testScope", "testScope2"))
	assert.NotNil(validator.ValidateKV("testScope.testCol", "testScope2.testCol"))

	validator = NewExplicitMappingValidator()
	assert.NotNil(validator.ValidateKV("scope1", ""))
	assert.NotNil(validator.ValidateKV("", "scope1"))
}

func TestRuleRedundancyNilCheck(t *testing.T) {
	assert := assert.New(t)
	fmt.Println("============== Test case start: TestRulesRedundancyCheck =================")
	defer fmt.Println("============== Test case end: TestRulesRedundancyCheck =================")

	validator := NewExplicitMappingValidator()
	assert.Nil(validator.ValidateKV("testScope", nil))
	assert.NotNil(validator.ValidateKV("testScope.testCol", nil))

	validator = NewExplicitMappingValidator()
	assert.Nil(validator.ValidateKV("testScope.testCol", nil))
	assert.NotNil(validator.ValidateKV("testScope", nil))
}

func TestRuleNameTooLong(t *testing.T) {
	assert := assert.New(t)
	fmt.Println("============== Test case start: TestRuleNameTooLong =================")
	defer fmt.Println("============== Test case end: TestRuleNameTooLong =================")

	longStr := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	validator := NewExplicitMappingValidator()
	assert.NotNil(validator.ValidateKV(longStr, longStr))
	ns2 := fmt.Sprintf("%v%v%v", longStr, ScopeCollectionDelimiter, longStr)
	assert.NotNil(validator.ValidateKV(ns2, ns2))
}

func TestArrayXattrFieldIterator(t *testing.T) {
	type fields struct {
		name    string
		xattr   string
		entries []string
		err     bool
	}
	tests := []fields{
		{
			name:    "empty array",
			xattr:   "[]",
			entries: []string{},
		},
		{
			name:    "one string entry",
			xattr:   `["foo"]`,
			entries: []string{`"foo"`},
		},
		{
			name:    "one JSON obj entry",
			xattr:   `[{"foo":"bar"}]`,
			entries: []string{`{"foo":"bar"}`},
		},
		{
			name:    "multiple string entries",
			xattr:   `["foo","bar","lorem","epsum"]`,
			entries: []string{`"foo"`, `"bar"`, `"lorem"`, `"epsum"`},
		},
		{
			name:    "multiple json obj entries",
			xattr:   `[{"foo":"foo1"},{"bar":"bar1"},{"lorem":"lorem1"},{"epsum":"epsum1"}]`,
			entries: []string{`{"foo":"foo1"}`, `{"bar":"bar1"}`, `{"lorem":"lorem1"}`, `{"epsum":"epsum1"}`},
		},
		{
			name:    "multiple string and json obj entries - I",
			xattr:   `["foo",{"bar":"bar1"},"lorem",{"epsum":"epsum1"}]`,
			entries: []string{`"foo"`, `{"bar":"bar1"}`, `"lorem"`, `{"epsum":"epsum1"}`},
		},
		{
			name:    "multiple string and json obj entries - II",
			xattr:   `[{"bar":"bar1"},"foo",{"epsum":"epsum1"},"lorem"]`,
			entries: []string{`{"bar":"bar1"}`, `"foo"`, `{"epsum":"epsum1"}`, `"lorem"`},
		},
		{
			name:    "one VV deltas entry",
			xattr:   `["NqiIe0LekFPLeX4JvTO6Iw@0x00008cd6ac059a16"]`,
			entries: []string{`"NqiIe0LekFPLeX4JvTO6Iw@0x00008cd6ac059a16"`},
		},
		{
			name:    "multiple VV deltas entry",
			xattr:   `["NqiIe0LekFPLeX4JvTO6Iw@0x00008cd6ac059a16","LhRPsa7CpjEvP5zeXTXEBA@0x0a","LhRPsa7CpjEvP5zsdsxEBA@0xffff"]`,
			entries: []string{`"NqiIe0LekFPLeX4JvTO6Iw@0x00008cd6ac059a16"`, `"LhRPsa7CpjEvP5zeXTXEBA@0x0a"`, `"LhRPsa7CpjEvP5zsdsxEBA@0xffff"`},
		},
		{
			name:    "invalid array with whitespaces",
			xattr:   `["foo", "bar", "lorem", "epsum"]`,
			entries: []string{`"foo"`, `"bar"`, `"lorem"`, `"epsum"`},
			err:     true,
		},
		{
			name:    "invalid entries",
			xattr:   `["foo",1,"lorem","epsum"]`,
			entries: []string{`"foo"`, `1`, `"lorem"`, `"epsum"`},
			err:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it, err := NewArrayXattrFieldIterator([]byte(tt.xattr))
			assert.Nil(t, err, nil)

			i := 0
			for it.HasNext() {
				val, err := it.Next()
				if err != nil {
					assert.True(t, tt.err)
					break
				} else {
					assert.Nil(t, err, nil)
					assert.Equal(t, tt.entries[i], string(val))
				}
				i++
			}
		})
	}
}

func Test_xtocIterator(t *testing.T) {
	tests := []struct {
		name string
		body []byte
		len  int
		list []string
	}{
		{
			name: "nil",
			body: nil,
			list: []string{},
		},
		{
			name: "empty",
			body: []byte{},
			len:  0,
			list: []string{},
		},
		{
			name: "empty list",
			body: []byte(`[]`),
			len:  0,
			list: []string{},
		},
		{
			name: "empty list with whitespace",
			body: []byte(`   [   ] `),
			len:  0,
			list: []string{},
		},
		{
			name: "one entry",
			list: []string{"foo"},
			len:  1,
			body: []byte(fmt.Sprintf(`["%s"]`, "foo")),
		},
		{
			name: "one entry with whitespaces",
			list: []string{"foo"},
			len:  1,
			body: []byte(fmt.Sprintf(` [  "%s"   ]   `, "foo")),
		},
		{
			name: "one entry with escape quote",
			list: []string{`fo\"o\"`},
			len:  1,
			body: []byte(fmt.Sprintf(`["%s"]`, `fo\"o\"`)),
		},
		{
			name: "one entry with whitespaces and escape quote",
			list: []string{`fo\"o\"`},
			len:  1,
			body: []byte(fmt.Sprintf(` [  "%s"   ]   `, `fo\"o\"`)),
		},
		{
			name: "multiple entries",
			list: []string{"foo", "bar", "foo1", "bar1"},
			len:  4,
			body: []byte(fmt.Sprintf(`["%s","%s","%s","%s"]`, "foo", "bar", "foo1", "bar1")),
		},
		{
			name: "multiple entries with whitespaces",
			list: []string{"foo", "bar", "foo1", "bar1"},
			len:  4,
			body: []byte(fmt.Sprintf(`  [  "%s","%s",  "%s",  "%s"  ]   `, "foo", "bar", "foo1", "bar1")),
		},

		{
			name: "multiple entries 1",
			list: []string{`fo\"o\"`, "bar", "foo1"},
			len:  3,
			body: []byte(fmt.Sprintf(`["%s", "%s","%s"]`, `fo\"o\"`, "bar", "foo1")),
		},
		{
			name: "multiple entries with whitespaces 1",
			list: []string{`fo\"o\"`, "bar", "foo1"},
			len:  3,
			body: []byte(fmt.Sprintf(`  [  "%s", "%s",  "%s"  ]   `, `fo\"o\"`, "bar", "foo1")),
		},

		{
			name: "multiple entries 2",
			list: []string{`foo`, `ba\"r\"`, "foo1"},
			len:  3,
			body: []byte(fmt.Sprintf(`["%s","%s","%s"]`, `foo`, `ba\"r\"`, "foo1")),
		},
		{
			name: "multiple entries with whitespaces 2",
			list: []string{`foo`, `ba\"r\"`, "foo1"},
			len:  3,
			body: []byte(fmt.Sprintf(`  [  "%s", "%s",  "%s"  ]   `, `foo`, `ba\"r\"`, "foo1")),
		},
		{
			name: "multiple entries 3",
			list: []string{`foo`, "bar", `foo\"1\"`},
			len:  3,
			body: []byte(fmt.Sprintf(`["%s","%s","%s"]`, `foo`, "bar", `foo\"1\"`)),
		},
		{
			name: "multiple entries with whitespaces 3",
			list: []string{`foo`, "bar", `foo\"1\"`},
			len:  3,
			body: []byte(fmt.Sprintf(`  [  "%s","%s",  "%s"  ]   `, `foo`, "bar", `foo\"1\"`)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			li, err := NewXTOCIterator(tt.body)
			assert.Nil(t, err)

			length, err := li.Len()
			assert.Nil(t, err)
			assert.Equal(t, length, tt.len)

			var res []string
			for li.HasNext() {
				s, err := li.Next()
				assert.Nil(t, err)

				res = append(res, string(s))
			}

			assert.Equal(t, len(res), len(tt.list))

			for i := 0; i < len(tt.list); i++ {
				assert.Equal(t, res[i], tt.list[i])
			}
		})
	}
}

func TestKvVBMapType_HasSameNumberOfVBs(t *testing.T) {
	type args struct {
		other KvVBMapType
	}
	tests := []struct {
		name string
		k    *KvVBMapType
		args args
		want bool
	}{
		{
			"Same number of VBs",
			&KvVBMapType{
				"a": []uint16{
					0, 1, 2,
				},
			},
			args{KvVBMapType{
				"a": []uint16{
					0, 1, 2,
				},
			}},
			true,
		},
		{
			"One nil map",
			nil,
			args{KvVBMapType{
				"a": []uint16{
					0, 1, 2,
				},
			}},
			false,
		},
		{
			"Both nil maps",
			nil,
			args{nil},
			true,
		},
		{
			"Different number of VBs",
			&KvVBMapType{
				"a": []uint16{
					0, 1, 2, 3,
				},
			},
			args{KvVBMapType{
				"a": []uint16{
					0, 1, 2,
				},
			}},
			false,
		},
		{
			"Test argument is nil",
			&KvVBMapType{
				"a": []uint16{
					0, 1, 2,
				},
			},
			args{nil},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.k.HasSameNumberOfVBs(tt.args.other), "HasSameNumberOfVBs(%v)", tt.args.other)
		})
	}
}
