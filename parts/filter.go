// Copyright (c) 2018-2019 Couchbase, Inc.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
// except in compliance with the License. You may obtain a copy of the License at
//   http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the
// License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific language governing permissions
// and limitations under the License.

package parts

import (
	"fmt"
	"github.com/couchbase/gojsonsm"
	mc "github.com/couchbase/gomemcached"
	mcc "github.com/couchbase/gomemcached/client"
	"github.com/couchbase/goxdcr/base"
	utilities "github.com/couchbase/goxdcr/utils"
)

type FilterIface interface {
	FilterUprEvent(uprEvent *mcc.UprEvent) bool
}

type DpSlicesHolder struct {
	slicesHolderCh chan [][]byte
}

func NewDpSlicesHolder() *DpSlicesHolder {
	return &DpSlicesHolder{
		slicesHolderCh: make(chan [][]byte, 20), // Up to 20 concurrent source nozzles
	}
}

func (dps *DpSlicesHolder) GetSliceOfSlices() [][]byte {
	for {
		select {
		case aSlice := <-dps.slicesHolderCh:
			return aSlice
		default:
			return make([][]byte, 0, 2)
		}
	}
}

func (dps *DpSlicesHolder) PutSliceOfSlices(doneSlice [][]byte) {
	select {
	case dps.slicesHolderCh <- doneSlice:
		return
	default:
		// Holder is full, don't block - just let it go to waste
		return
	}
}

type Filter struct {
	id                       string
	filterExpressionInternal string
	matcher                  gojsonsm.Matcher
	utils                    utilities.UtilsIface
	dp                       utilities.DataPoolIface
	flags                    base.FilterFlagType
	slicesHolder             *DpSlicesHolder
}

func NewFilter(id string, filterExpression string, utils utilities.UtilsIface) (*Filter, error) {
	dpPtr := utilities.NewDataPool()
	if dpPtr == nil {
		return nil, base.ErrorNoDataPool
	}

	if len(filterExpression) == 0 {
		return nil, base.ErrorInvalidInput
	}

	filter := &Filter{
		id:                       id,
		filterExpressionInternal: base.ReplaceKeyWordsForExpression(filterExpression),
		utils:                    utils,
		dp:                       dpPtr,
		slicesHolder:             NewDpSlicesHolder(),
	}

	matcher, err := base.GoJsonsmGetFilterExprMatcher(filter.filterExpressionInternal)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse expression: %v Err: %v", filter.filterExpressionInternal, err.Error())
	}
	filter.matcher = matcher

	if filter.matcher == nil {
		return nil, base.ErrorNoMatcher
	}

	if !base.FilterContainsXattrExpression(filterExpression) {
		filter.flags |= base.FilterFlagSkipXattr
	}

	if !base.FilterContainsKeyExpression(filterExpression) {
		filter.flags |= base.FilterFlagSkipKey
	} else if base.FilterOnlyContainsKeyExpression(filterExpression) {
		filter.flags |= base.FilterFlagKeyOnly
	}

	return filter, nil
}

// Returns:
// 1. bool - Whether or not it was a match
// 2. err code
// 3. If err is not nil, additional description
// 4. Total bytes of failed datapool gets - which means len of []byte alloc (garbage)
func (filter *Filter) FilterUprEvent(uprEvent *mcc.UprEvent) (bool, error, string, int64) {
	var err error
	if uprEvent == nil {
		return false, base.ErrorInvalidInput, "UprEvent is nil", 0
	}

	if uprEvent.Opcode == mc.UPR_DELETION || uprEvent.Opcode == mc.UPR_EXPIRATION {
		// For now, pass through
		return true, nil, "", 0
	}

	slicesToBeReleased := filter.slicesHolder.GetSliceOfSlices()
	// defer is last in first out - this is going to be the first defer on the stack, which will execute after releaseFunc() below has finished recycle datapools
	defer filter.slicesHolder.PutSliceOfSlices(slicesToBeReleased)

	sliceToBeFiltered, err, errDesc, releaseFunc, failedDpCnt := filter.utils.ProcessUprEventForFiltering(uprEvent, filter.dp, filter.flags, &slicesToBeReleased)
	if releaseFunc != nil {
		// releaseFunc() will do cleanup of slicesToBeReleased above
		defer releaseFunc()
	}

	if err != nil {
		if err == base.FilterForcePassThrough {
			return true, nil, "", failedDpCnt
		} else {
			return false, err, errDesc, failedDpCnt
		}
	}
	matched, err := filter.FilterByteSlice(sliceToBeFiltered)
	if err != nil {
		errDesc = fmt.Sprintf("gojsonsm filter returned err for document %v%v%v", base.UdTagBegin, string(uprEvent.Key), base.UdTagEnd)
	}

	return matched, err, errDesc, failedDpCnt
}

func (filter *Filter) matchWrapper(slice []byte, errPtr *error) (matched bool) {
	defer func() {
		if r := recover(); r != nil {
			*errPtr = fmt.Errorf("Error from matcher: %v", r)
		}
	}()
	matched, *errPtr = filter.matcher.Match(slice)
	return matched
}

func (filter *Filter) FilterByteSlice(slice []byte) (matched bool, err error) {
	defer filter.matcher.Reset()
	matched = filter.matchWrapper(slice, &err)
	return
}
