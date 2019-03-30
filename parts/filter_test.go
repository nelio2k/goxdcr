// +build !pcre

package parts

import (
	"encoding/json"
	"fmt"
	mcc "github.com/couchbase/gomemcached/client"
	"github.com/couchbase/goxdcr/base"
	utilities "github.com/couchbase/goxdcr/utils"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

var filterId string = "testFilter"
var realUtil *utilities.Utilities = utilities.NewUtilities()
var dp utilities.DataPoolIface = utilities.NewDataPool()

func TestFilter(t *testing.T) {
	fmt.Println("============== Test case start: TestFilterCreation =================")
	assert := assert.New(t)

	filter, err := NewFilter(filterId, "REGEXP_CONTAINS(META().id, \"^abc\")", realUtil)

	assert.Nil(err)
	assert.NotNil(filter)
	assert.Equal(2, len(base.ReservedWordsReplaceMap))
	assert.True(filter.flags&base.FilterFlagSkipKey == 0)
	assert.True(filter.flags&base.FilterFlagSkipXattr > 0)

	fmt.Println("============== Test case end: TestFilterCreation =================")
}

func TestFilterBool(t *testing.T) {
	fmt.Println("============== Test case start: TestFilterBool =================")
	assert := assert.New(t)

	udMarsh, err := ioutil.ReadFile("./testdata/MB-33010.json")
	assert.Nil(err)

	filter, err := NewFilter(filterId, "bool=false", realUtil)
	assert.Nil(err)
	assert.NotNil(filter)

	matchResult, err := filter.matcher.Match(udMarsh)
	assert.Nil(err)
	assert.True(matchResult)

	fmt.Println("============== Test case end: TestFilterBool =================")
}

func RetrieveUprFile(fileName string) (*mcc.UprEvent, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var uprEvent mcc.UprEvent
	err = json.Unmarshal(data, &uprEvent)
	if err != nil {
		return nil, err
	}

	return &uprEvent, nil
}

func TestFilterBool2(t *testing.T) {
	fmt.Println("============== Test case start: TestFilterBool2 =================")
	assert := assert.New(t)

	uprEvent, err := RetrieveUprFile("./testdata/MB-33010.json")
	assert.Nil(err)
	assert.NotNil(uprEvent)

	filter, err := NewFilter(filterId, "bool=false AND NOT int > dfdfk", realUtil)
	assert.Nil(err)
	assert.NotNil(filter)

	slices := make([][]byte, 0, 2)
	dataSlice, err, _, releaseFunc, _ := realUtil.ProcessUprEventForFiltering(uprEvent, dp, 0 /*flags*/, &slices)
	assert.Nil(err)
	assert.NotNil(releaseFunc)

	matchResult, err := filter.matcher.Match(dataSlice)
	assert.Nil(err)
	assert.False(matchResult)

	fmt.Println("============== Test case end: TestFilterBool2 =================")
}

func TestFilterPerf(t *testing.T) {
	fmt.Println("============== Test case start: TestFilterPerf =================")
	assert := assert.New(t)

	uprEvent, err := RetrieveUprFile("./testdata/perfData.bin")
	assert.Nil(err)
	assert.NotNil(uprEvent)

	filter, err := NewFilter(filterId, "city EXISTS", realUtil)
	assert.Nil(err)
	assert.NotNil(filter)
	assert.Equal(base.FilterFlagType(3), filter.flags)

	slices := make([][]byte, 0, 2)
	dataSlice, err, _, releaseFunc, _ := realUtil.ProcessUprEventForFiltering(uprEvent, dp, filter.flags, &slices)
	assert.Nil(err)
	assert.NotNil(releaseFunc)

	matchResult, err := filter.matcher.Match(dataSlice)
	assert.Nil(err)
	assert.True(matchResult)

	fmt.Println("============== Test case end: TestFilterPerf =================")
}

func TestFilterPerfKeyOnly(t *testing.T) {
	fmt.Println("============== Test case start: TestFilterPerfKeyOnly =================")
	assert := assert.New(t)

	uprEvent, err := RetrieveUprFile("./testdata/perfData.bin")
	assert.Nil(err)
	assert.NotNil(uprEvent)

	filter, err := NewFilter(filterId, "META().id = \"26dcc0-000000001586\"", realUtil)
	assert.Nil(err)
	assert.NotNil(filter)
	assert.True(filter.flags&base.FilterFlagKeyOnly > 0)

	slices := make([][]byte, 0, 2)
	dataSlice, err, _, releaseFunc, _ := realUtil.ProcessUprEventForFiltering(uprEvent, dp, filter.flags, &slices)
	assert.Nil(err)
	assert.NotNil(releaseFunc)

	matchResult, err := filter.matcher.Match(dataSlice)
	assert.Nil(err)
	assert.True(matchResult)

	fmt.Println("============== Test case end: TestFilterPerfKeyOnly =================")
}

func TestKeyPanic(t *testing.T) {
	fmt.Println("============== Test case start: TestKeyPanic =================")
	assert := assert.New(t)

	uprEvent, err := RetrieveUprFile("./testdata/edgyMB-33583.json")
	assert.Nil(err)
	assert.NotNil(uprEvent)

	filter, err := NewFilter(filterId, "REGEXP_CONTAINS(META().id, \"C1-key-1\")", realUtil)
	assert.Nil(err)
	assert.NotNil(filter)
	assert.True(filter.flags&base.FilterFlagKeyOnly > 0)

	slices := make([][]byte, 0, 2)
	dataSlice, err, _, releaseFunc, _ := realUtil.ProcessUprEventForFiltering(uprEvent, dp, filter.flags, &slices)
	assert.Nil(err)
	assert.NotNil(releaseFunc)

	matchResult, err := filter.matcher.Match(dataSlice)
	assert.Nil(err)
	assert.False(matchResult)

	fmt.Println("============== Test case end: TestKeyPanic =================")
}

func TestDpHolder(t *testing.T) {
	fmt.Println("============== Test case start: TestDpHolder =================")
	assert := assert.New(t)
	holder := NewDpSlicesHolder()
	assert.Equal(0, len(holder.slicesHolderCh))

	slices := holder.GetSliceOfSlices()
	assert.Equal(0, len(holder.slicesHolderCh))

	holder.PutSliceOfSlices(slices)
	assert.Equal(1, len(holder.slicesHolderCh))

	slices = holder.GetSliceOfSlices()
	assert.Equal(0, len(holder.slicesHolderCh))

	slices2 := holder.GetSliceOfSlices()
	holder.PutSliceOfSlices(slices)
	holder.PutSliceOfSlices(slices2)
	assert.Equal(2, len(holder.slicesHolderCh))

	// Get > 20 slices
	var holderOfHolder [][][]byte
	for i := 0; i < 21; i++ {
		holderOfHolder = append(holderOfHolder, holder.GetSliceOfSlices())
	}

	assert.Equal(0, len(holder.slicesHolderCh))
	assert.Equal(21, len(holderOfHolder))

	for i := 0; i < 21; i++ {
		// should not block
		holder.PutSliceOfSlices(holderOfHolder[i])
	}

	// Last one was tossed away
	assert.Equal(20, len(holder.slicesHolderCh))
	fmt.Println("============== Test case end: TestDpHolder =================")
}
