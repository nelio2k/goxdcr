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
	mcc "github.com/couchbase/gomemcached/client"
	"github.com/couchbase/goxdcr/base"
	utilities "github.com/couchbase/goxdcr/utils"
	"github.com/couchbaselabs/gojsonsm"
)

type FilterIface interface {
	FilterUprEvent(uprEvent *mcc.UprEvent) bool
}

type Filter struct {
	id                       string
	filterExpression         string
	filterExpressionInternal string
	matcher                  *gojsonsm.Matcher
	utils                    utilities.UtilsIface
	dp                       utilities.DataPoolIface
}

func NewFilter(id string, filterExpression string, utils utilities.UtilsIface) (*Filter, error) {
	dpPtr := utilities.NewDataPool()
	if dpPtr == nil {
		return nil, base.ErrorNoDataPool
	}

	filter := &Filter{
		id:                       id,
		filterExpression:         filterExpression,
		filterExpressionInternal: base.ReplaceKeyWordsForExpression(filterExpression),
		utils: utils,
		dp:    dpPtr,
	}

	if len(filterExpression) == 0 {
		return nil, base.ErrorInvalidInput
	}

	expression, err := gojsonsm.ParseSimpleExpression(filter.filterExpressionInternal)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse expression: %v", err.Error())
	}
	var trans gojsonsm.Transformer
	matchDef := trans.Transform([]gojsonsm.Expression{expression})
	filter.matcher = gojsonsm.NewMatcher(matchDef)

	if filter.matcher == nil {
		return nil, base.ErrorNoMatcher
	}

	return filter, nil
}

func (filter *Filter) FilterUprEvent(uprEvent *mcc.UprEvent) bool {
	sliceToBeFiltered, err, releaseFunc := filter.utils.ProcessUprEventForFiltering(uprEvent, filter.dp)
	if err != nil {
		return false
	}
	defer releaseFunc()
	filterSuccess := filter.FilterByteSlice(sliceToBeFiltered)
	return filterSuccess
}

func (filter *Filter) FilterByteSlice(slice []byte) bool {
	defer filter.matcher.Reset()
	match, _ := filter.matcher.Match(slice)
	return match
}
