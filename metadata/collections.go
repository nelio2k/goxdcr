// Copyright (c) 2013-2020 Couchbase, Inc.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
// except in compliance with the License. You may obtain a copy of the License at
//   http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the
// License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific language governing permissions
// and limitations under the License.

package metadata

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/couchbase/goxdcr/base"
	"github.com/golang/snappy"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type ManifestsDoc struct {
	collectionsManifests []*CollectionsManifest

	// When upserting or retrieving from metakv, compress them into a single byte slice
	// This is not to be used otherwise - use CollectionsManifests()
	CompressedCollectionsManifests []byte `json:"collection_manifests"`

	//revision number
	revision interface{}
}

func (m *ManifestsDoc) Revision() interface{} {
	return m.revision
}

func (m *ManifestsDoc) SetRevision(rev interface{}) {
	m.revision = rev
}

func (m *ManifestsDoc) CollectionsManifests() []*CollectionsManifest {
	return m.collectionsManifests
}

func (m *ManifestsDoc) SetCollectionsManifests(manifests []*CollectionsManifest) {
	m.collectionsManifests = manifests
}

func (m *ManifestsDoc) PreMarshal() error {
	serializedJson, err := json.Marshal(m.collectionsManifests)
	if err != nil {
		return err
	}
	m.CompressedCollectionsManifests = snappy.Encode(nil, serializedJson)
	return nil
}

func (m *ManifestsDoc) ClearCompressedData() {
	m.CompressedCollectionsManifests = nil
}

func (m *ManifestsDoc) PostUnmarshal() error {
	var serializedJson []byte
	serializedJson, err := snappy.Decode(serializedJson, m.CompressedCollectionsManifests)
	if err != nil {
		return err
	}

	err = json.Unmarshal(serializedJson, &(m.collectionsManifests))
	if err != nil {
		return err
	}

	return nil
}

type CollectionsManifestPair struct {
	Source *CollectionsManifest `json:"sourceManifest"`
	Target *CollectionsManifest `json:"targetManifest"`
}

func (c *CollectionsManifestPair) IsSameAs(other *CollectionsManifestPair) bool {
	return c.Source.IsSameAs(other.Source) && c.Target.IsSameAs(other.Target)
}

func NewCollectionsManifestPair(source, target *CollectionsManifest) *CollectionsManifestPair {
	return &CollectionsManifestPair{
		Source: source,
		Target: target,
	}
}

// Manifest structure representing the JSON returned from ns_server's endpoint
type CollectionsManifest struct {
	uid    uint64
	scopes ScopesMap
}

func UnitTestGenerateCollManifest(uid uint64, scopes ScopesMap) *CollectionsManifest {
	return &CollectionsManifest{uid, scopes}
}

func (c *CollectionsManifest) String() string {
	if c == nil {
		return ""
	}
	var output []string
	output = append(output, fmt.Sprintf("CollectionsManifest uid: %v { ", c.uid))
	output = append(output, c.scopes.String())
	output = append(output, " }")
	return strings.Join(output, " ")
}

func (c *CollectionsManifest) GetScopeAndCollectionName(collectionId uint32) (scopeName, collectionName string, err error) {
	if c == nil {
		err = base.ErrorInvalidInput
		return
	}
	for _, scope := range c.Scopes() {
		for _, collection := range scope.Collections {
			if collection.Uid == collectionId {
				scopeName = scope.Name
				collectionName = collection.Name
				return
			}
		}
	}
	err = base.ErrorNotFound
	return
}

func (c *CollectionsManifest) GetCollectionId(scopeName, collectionName string) (uint32, error) {
	if c == nil {
		return 0, base.ErrorInvalidInput
	}
	for _, scope := range c.Scopes() {
		if scopeName == scope.Name {
			for _, collection := range scope.Collections {
				if collection.Name == collectionName {
					return collection.Uid, nil
				}
			}
		}
	}
	return 0, base.ErrorNotFound
}

func (target *CollectionsManifest) ImplicitGetBackfillCollections(prevTarget, source *CollectionsManifest) (backfillsNeeded CollectionNamespaceMapping, err error) {
	backfillsNeeded = make(CollectionNamespaceMapping)

	// First, find any previously unmapped target collections that are now mapped
	oldSrcToTargetMapping, _, _ := source.ImplicitMap(prevTarget)
	srcToTargetMapping, _, _ := source.ImplicitMap(target)

	added, _ := oldSrcToTargetMapping.Diff(srcToTargetMapping)
	backfillsNeeded.Consolidate(added)

	// Then, find out if any target Collection UID changed from underneath the manifests
	// These target collections need to be backfilled that couldn't be discovered from the namespacemappings
	// because they were never "missing", but their collection UID changed
	_, changedTargets, _, err := target.Diff(prevTarget)
	if err != nil {
		return nil, err
	}
	backfillsNeededDueToTargetRecreation := srcToTargetMapping.GetSubsetBasedOnSpecifiedTargets(changedTargets)

	backfillsNeeded.Consolidate(backfillsNeededDueToTargetRecreation)
	return
}

// Actual obj stored in metakv
type collectionsMetaObj struct {
	// data for unmarshalling and parsing
	UidMeta_    string        `json:"uid"`
	ScopesMeta_ []interface{} `json:"scopes"`
}

func newCollectionsMetaObj() *collectionsMetaObj {
	return &collectionsMetaObj{}
}

// Used if there's an error getting collections manifest
func NewDefaultCollectionsManifest() CollectionsManifest {
	defaultManifest := CollectionsManifest{scopes: make(ScopesMap)}

	defaultScope := NewEmptyScope(base.DefaultScopeCollectionName, 0)
	defaultCollection := Collection{Uid: 0, Name: base.DefaultScopeCollectionName}

	defaultScope.Collections[base.DefaultScopeCollectionName] = defaultCollection

	defaultManifest.scopes[base.DefaultScopeCollectionName] = defaultScope

	return defaultManifest
}

func NewCollectionsManifestFromMap(manifestInfo map[string]interface{}) (CollectionsManifest, error) {
	metaObj := newCollectionsMetaObj()
	var manifest CollectionsManifest

	if uid, ok := manifestInfo["uid"].(string); ok {
		metaObj.UidMeta_ = uid
	} else {
		return manifest, fmt.Errorf("Uid is not float64, but %v instead", reflect.TypeOf(manifestInfo["uid"]))
	}

	if scopes, ok := manifestInfo["scopes"].([]interface{}); ok {
		metaObj.ScopesMeta_ = scopes
	} else {
		return manifest, base.ErrorInvalidInput
	}

	err := manifest.Load(metaObj)
	if err != nil {
		return manifest, err
	}
	return manifest, nil
}

func NewCollectionsManifestFromBytes(data []byte) (CollectionsManifest, error) {
	var manifest CollectionsManifest
	err := manifest.LoadBytes(data)
	return manifest, err
}

// For unit test
func TestNewCollectionsManifestFromBytesWithCustomUid(data []byte, uid uint64) (CollectionsManifest, error) {
	var manifest CollectionsManifest
	err := manifest.LoadBytes(data)
	manifest.uid = uid
	return manifest, err
}

// Does not clone temporary variables
func (c CollectionsManifest) Clone() CollectionsManifest {
	return CollectionsManifest{
		uid:    c.uid,
		scopes: c.scopes.Clone(),
	}
}

func (this *CollectionsManifest) IsSameAs(other *CollectionsManifest) bool {
	if this == nil || other == nil {
		return false
	}

	if this.Uid() != other.Uid() {
		return false
	}

	for scopeName, scope := range this.scopes {
		otherScope, ok := other.scopes[scopeName]
		if !ok {
			return false
		}
		if !scope.IsSameAs(otherScope) {
			return false
		}
	}

	return true
}

func (c *CollectionsManifest) Diff(older *CollectionsManifest) (added, modified, removed ScopesMap, err error) {
	if c == nil || older == nil {
		return ScopesMap{}, ScopesMap{}, ScopesMap{}, base.ErrorInvalidInput
	}

	if c.Uid() < older.Uid() {
		err = fmt.Errorf("Should compare against an older version of manifest")
		return
	} else if c.IsSameAs(older) {
		return
	}

	added = make(ScopesMap)
	modified = make(ScopesMap)
	removed = make(ScopesMap)

	// First, find things that exists in the current manifest that doesn't exist in the older one
	// or versions are different, which by contract means they are modified and newer (uid never go backwards)
	for scopeName, scope := range c.Scopes() {
		olderScope, exists := older.Scopes()[scopeName]
		if !exists {
			added[scopeName] = scope
		} else if exists && !scope.IsSameAs(olderScope) {
			collectionAddedYet := false
			collectionModifiedYet := false
			// At least one collection is different
			// If this scope is different because all the collections are removed, the "removed" portion
			// should catch it
			for collectionName, collection := range scope.Collections {
				olderCollection, exists := olderScope.Collections[collectionName]
				if !exists {
					if !collectionAddedYet {
						added[scopeName] = NewEmptyScope(scopeName, scope.Uid)
						collectionAddedYet = true
					}
					added[scopeName].Collections[collectionName] = collection
				} else if exists && !collection.IsSameAs(olderCollection) {
					if !collectionModifiedYet {
						modified[scopeName] = NewEmptyScope(scopeName, scope.Uid)
						collectionModifiedYet = true
					}
					modified[scopeName].Collections[collectionName] = collection
				}
			}
		}
	}

	// Then find things that don't exist in the new but exist in the old
	if c.Count(true /*includeDefaultScope*/, true /*includeDefaultCollection*/) != older.Count(true, true) {
		for olderScopeName, olderScope := range older.Scopes() {
			scope, exists := c.Scopes()[olderScopeName]
			if !exists {
				removed[olderScopeName] = olderScope
			} else if exists && olderScope.Count(true /*includeDefaultCollection*/) != scope.Count(true) {
				// At least one collection is removed, potentially including default collection
				removed[olderScopeName] = NewEmptyScope(olderScopeName, olderScope.Uid)
				for olderCollectionName, olderCollection := range olderScope.Collections {
					_, exists := scope.Collections[olderCollectionName]
					if !exists {
						removed[olderScopeName].Collections[olderCollectionName] = olderCollection
					}
				}
			}
		}
	}

	return
}

// Counts the total number of collections in this manifest, including or excluding default scope/collection
func (c *CollectionsManifest) Count(includeDefaultScope, includeDefaultCollection bool) int {
	return c.scopes.Count(includeDefaultScope, includeDefaultCollection)
}

// Given a metadata object, load into the more user-friendly collectionsManifest
func (c *CollectionsManifest) Load(collectionsMeta *collectionsMetaObj) error {
	var err error
	c.uid, err = strconv.ParseUint(collectionsMeta.UidMeta_, base.CollectionsUidBase, 64)
	c.scopes = make(ScopesMap)
	for _, oneScopeDetail := range collectionsMeta.ScopesMeta_ {
		scopeDetailMap, ok := oneScopeDetail.(map[string]interface{})
		if !ok {
			return base.ErrorInvalidInput
		}
		scopeName, ok := scopeDetailMap[base.NameKey].(string)
		if !ok {
			return base.ErrorInvalidInput
		}
		c.scopes[scopeName], err = NewScope(scopeName, scopeDetailMap)
		if err != nil {
			return base.ErrorInvalidInput
		}
	}
	return nil
}

func (c *CollectionsManifest) LoadBytes(data []byte) error {
	collectionsMeta := newCollectionsMetaObj()
	err := json.Unmarshal(data, collectionsMeta)
	if err != nil {
		return err
	}

	return c.Load(collectionsMeta)
}

// Implements the marshaller interface
func (c *CollectionsManifest) MarshalJSON() ([]byte, error) {
	collectionsMeta := newCollectionsMetaObj()
	collectionsMeta.UidMeta_ = fmt.Sprintf("%x", c.uid)

	// marshal scopes in order of names - this will ensure equality between two identical manifests
	var scopeNames []string
	for name, _ := range c.scopes {
		scopeNames = append(scopeNames, name)
	}
	scopeNames = base.SortStringList(scopeNames)

	for _, scopeName := range scopeNames {
		scope := c.scopes[scopeName]
		collectionsMeta.ScopesMeta_ = append(collectionsMeta.ScopesMeta_, scope.toScopeDetail())
	}

	outBytes, _ := json.Marshal(collectionsMeta)
	return outBytes, nil
}

// Implements the marshaller interface
func (c *CollectionsManifest) UnmarshalJSON(b []byte) error {
	collectionsMeta := newCollectionsMetaObj()

	err := json.Unmarshal(b, collectionsMeta)
	if err != nil {
		return err
	}

	return c.Load(collectionsMeta)
}

// Returns a sha256 representing the manifest data
func (c *CollectionsManifest) Sha256() (result [sha256.Size]byte, err error) {
	var jsonBytes []byte
	jsonBytes, err = c.MarshalJSON()
	if err != nil {
		return
	}

	result = sha256.Sum256(jsonBytes)
	return
}

func (c *CollectionsManifest) Uid() uint64 {
	return c.uid
}

func (c *CollectionsManifest) Scopes() ScopesMap {
	return c.scopes
}

//type CollectionsMap map[string]Collection
// Diff by name between two manifests
func (sourceManifest *CollectionsManifest) ImplicitMap(targetManifest *CollectionsManifest) (successfulMapping CollectionNamespaceMapping, unmappedSources CollectionsMap, unmappedTargets CollectionsMap) {
	if sourceManifest == nil || targetManifest == nil {
		return
	}

	successfulMapping = make(CollectionNamespaceMapping)
	unmappedSources = make(CollectionsMap)
	unmappedTargets = make(CollectionsMap)

	for _, sourceScope := range sourceManifest.Scopes() {
		_, exists := targetManifest.Scopes()[sourceScope.Name]
		if !exists {
			// Whole scope does not exist on target
			for _, collection := range sourceScope.Collections {
				unmappedSources[collection.Name] = collection
			}
		} else {
			// Some collections may or may not exist on target
			for _, sourceCollection := range sourceScope.Collections {
				_, err := targetManifest.GetCollectionId(sourceScope.Name, sourceCollection.Name)
				if err == nil {
					implicitNamespace := &base.CollectionNamespace{sourceScope.Name, sourceCollection.Name}
					successfulMapping.AddSingleMapping(implicitNamespace, implicitNamespace)
				} else {
					unmappedSources[sourceCollection.Name] = sourceCollection
				}
			}
		}
	}

	for _, targetScope := range targetManifest.Scopes() {
		_, exists := sourceManifest.Scopes()[targetScope.Name]
		if !exists {
			// Whole scope does not exist on source
			for _, collection := range targetScope.Collections {
				unmappedTargets[collection.Name] = collection
			}
		} else {
			// Some may or maynot exist on source
			for _, targetCollection := range targetScope.Collections {
				_, err := sourceManifest.GetCollectionId(targetScope.Name, targetCollection.Name)
				if err != nil {
					unmappedTargets[targetCollection.Name] = targetCollection
				}
			}
		}
	}
	return
}

type CollectionsPtrList []*Collection

func (c CollectionsPtrList) Len() int           { return len(c) }
func (c CollectionsPtrList) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c CollectionsPtrList) Less(i, j int) bool { return c[i].Uid < c[j].Uid }

func SortCollectionsPtrList(list CollectionsPtrList) CollectionsPtrList {
	sort.Sort(list)
	return list
}

func (c CollectionsPtrList) IsSameAs(other CollectionsPtrList) bool {
	if len(c) != len(other) {
		return false
	}

	// Lists are logically "equal" if they have the same items but in diff order
	aList := SortCollectionsPtrList(c)
	bList := SortCollectionsPtrList(other)

	for i, col := range aList {
		if !col.IsSameAs(*(bList[i])) {
			return false
		}
	}
	return true
}

type CollectionsList []Collection

func (c CollectionsList) Len() int           { return len(c) }
func (c CollectionsList) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c CollectionsList) Less(i, j int) bool { return c[i].Uid < c[j].Uid }

func SortCollectionsList(list CollectionsList) CollectionsList {
	sort.Sort(list)
	return list
}

func (c CollectionsList) IsSameAs(other CollectionsList) bool {
	if len(c) != len(other) {
		return false
	}

	// Lists are logically "equal" if they have the same items but in diff order
	aList := SortCollectionsList(c)
	bList := SortCollectionsList(other)

	for i, col := range aList {
		if !col.IsSameAs(bList[i]) {
			return false
		}
	}
	return true
}

// Used to support marshalling CollectionToCollectionsMapping
type c2cMarshalObj struct {
	SourceCollections []*Collection `json:Source`
	// keys are integers of the index above written as strings
	IndirectTargetMap map[string][]*Collection `json:Map`
}

func newc2cMarshalObj() *c2cMarshalObj {
	return &c2cMarshalObj{
		IndirectTargetMap: make(map[string][]*Collection),
	}
}

type ScopesMap map[string]Scope

func (s ScopesMap) Len() int {
	return len(s)
}

func (s *ScopesMap) String() string {
	if s == nil {
		return "(None)"
	}
	var output []string
	for _, scope := range *s {
		output = append(output, scope.String())
	}
	return strings.Join(output, " ")
}

func (s *ScopesMap) Count(includeDefaultScope, includeDefaultCollection bool) int {
	var count int
	for scopeName, scope := range *s {
		if scopeName == base.DefaultScopeCollectionName {
			if includeDefaultScope {
				count += scope.Count(includeDefaultCollection)
			}
		} else {
			count += scope.Count(includeDefaultCollection)
		}
	}
	return count

}

func (s ScopesMap) Clone() ScopesMap {
	clone := make(ScopesMap)
	for k, v := range s {
		clone[k] = v.Clone()
	}
	return clone
}

func (s ScopesMap) GetCollection(id uint32) (col Collection, found bool) {
	for _, scope := range s {
		for _, collection := range scope.Collections {
			if collection.Uid == id {
				col = collection
				found = true
				return
			}
		}
	}
	return
}

func (s ScopesMap) GetCollectionByNames(scopeName, collectionName string) (col Collection, found bool) {
	scope, found := s[scopeName]
	if !found {
		return
	}
	col, found = scope.Collections[collectionName]
	return
}

type Scope struct {
	Uid         uint32 `json:"Uid"`
	Name        string `json:"Name"`
	Collections CollectionsMap
}

func (s *Scope) String() string {
	if s == nil {
		return ""
	}
	var output []string
	output = append(output, fmt.Sprintf("Scope Uid: %v Name: %v { ", s.Uid, s.Name))
	output = append(output, s.Collections.String())
	output = append(output, " }")
	return strings.Join(output, " ")
}

func NewScope(name string, scopeDetail map[string]interface{}) (Scope, error) {
	uidString, ok := scopeDetail[base.UIDKey].(string)
	if !ok {
		return Scope{}, fmt.Errorf("Uid is not float64, but %v instead", reflect.TypeOf(scopeDetail[base.UIDKey]))
	}
	uid64, err := strconv.ParseUint(uidString, base.CollectionsUidBase, 64)
	if err != nil {
		return Scope{}, err
	}
	uid32 := uint32(uid64)

	collectionsDetail, ok := scopeDetail[base.CollectionsKey].([]interface{})
	if !ok {
		return Scope{}, base.ErrorInvalidInput
	}

	collectionsMap, err := NewCollectionsMap(collectionsDetail)
	if err != nil {
		return Scope{}, err
	}

	return Scope{
		Name:        name,
		Uid:         uid32,
		Collections: collectionsMap,
	}, nil
}

func NewEmptyScope(name string, uid uint32) Scope {
	return Scope{
		Name:        name,
		Uid:         uid,
		Collections: make(CollectionsMap),
	}
}

func (this *Scope) IsSameAs(other Scope) bool {
	if this.Uid != other.Uid {
		return false
	}
	if this.Name != other.Name {
		return false
	}
	if !this.Collections.IsSameAs(other.Collections) {
		return false
	}
	return true
}

func (s *Scope) toScopeDetail() map[string]interface{} {
	detailMap := make(map[string]interface{})
	detailMap[base.NameKey] = s.Name
	detailMap[base.UIDKey] = fmt.Sprintf("%x", s.Uid)
	detailMap[base.CollectionsKey] = s.Collections.toCollectionsDetail()

	return detailMap
}

func (s *Scope) Count(includeDefaultCollection bool) int {
	return s.Collections.Count(includeDefaultCollection)
}

func (s Scope) Clone() Scope {
	return Scope{
		Uid:         s.Uid,
		Name:        s.Name,
		Collections: s.Collections.Clone(),
	}
}

type Collection struct {
	Uid  uint32 `json:"Uid"`
	Name string `json:"Name"`
}

func (c *Collection) String() string {
	if c == nil {
		return ""
	}
	return fmt.Sprintf("Collection { Uid: %v Name: %v }", c.Uid, c.Name)
}

func NewCollectionsMap(collectionsList []interface{}) (map[string]Collection, error) {
	collectionMap := make(CollectionsMap)

	for _, detail := range collectionsList {
		colDetail, ok := detail.(map[string]interface{})
		if !ok {
			return nil, base.ErrorInvalidInput
		}

		name, ok := colDetail[base.NameKey].(string)
		if !ok {
			return nil, base.ErrorInvalidInput
		}

		uidStr, ok := colDetail[base.UIDKey].(string)
		if !ok {
			return nil, base.ErrorInvalidInput
		}
		uid64, err := strconv.ParseUint(uidStr, base.CollectionsUidBase, 64)
		if err != nil {
			return nil, err
		}
		uid32 := uint32(uid64)

		collectionMap[name] = Collection{
			Uid:  uid32,
			Name: name,
		}
	}
	return collectionMap, nil
}

func (this *Collection) IsSameAs(other Collection) bool {
	return this.Uid == other.Uid && this.Name == other.Name
}

func (c Collection) Clone() Collection {
	return Collection{
		Uid:  c.Uid,
		Name: c.Name,
	}
}

type CollectionsMap map[string]Collection

func (c *CollectionsMap) String() string {
	var output []string
	if c == nil {
		return ""
	}
	for _, col := range *c {
		output = append(output, col.String())
	}
	return strings.Join(output, " ")
}

func (c *CollectionsMap) toCollectionsDetail() (detailList []interface{}) {
	// Output in sorted collection name ordering
	var colNames []string
	for name, _ := range *c {
		colNames = append(colNames, name)
	}
	colNames = base.SortStringList(colNames)

	for _, colName := range colNames {
		collection := (*c)[colName]
		detail := make(map[string]interface{})
		detail[base.NameKey] = colName
		detail[base.UIDKey] = fmt.Sprintf("%x", collection.Uid)
		detailList = append(detailList, detail)
	}
	return
}

func (this *CollectionsMap) IsSameAs(other CollectionsMap) bool {
	for colName, collection := range *this {
		otherCollection, ok := other[colName]
		if !ok {
			return false
		}
		if !collection.IsSameAs(otherCollection) {
			return false
		}
	}
	return true
}

func (c *CollectionsMap) Count(includeDefaultCollection bool) int {
	var count int
	for colName, _ := range *c {
		if colName == base.DefaultScopeCollectionName {
			if includeDefaultCollection {
				count++
			}
		} else {
			count++
		}
	}
	return count
}

func (c CollectionsMap) Clone() CollectionsMap {
	clone := make(CollectionsMap)
	for k, v := range c {
		clone[k] = v.Clone()
	}
	return clone
}

// Diff by name
// Modified is if the names are the same but collection IDs are different
func (c CollectionsMap) Diff(older CollectionsMap) (added, removed, modified CollectionsMap) {
	added = make(CollectionsMap)
	removed = make(CollectionsMap)
	modified = make(CollectionsMap)

	for collectionName, collection := range c {
		olderCol, exists := older[collectionName]
		if !exists {
			added[collectionName] = collection
		} else if olderCol.Uid != collection.Uid {
			modified[collectionName] = collection
		}
	}

	for olderColName, olderCol := range older {
		_, exists := c[olderColName]
		if !exists {
			removed[olderColName] = olderCol
		}
	}

	return
}

type ManifestsList []*CollectionsManifest

func (c ManifestsList) Len() int           { return len(c) }
func (c ManifestsList) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ManifestsList) Less(i, j int) bool { return c[i].uid < c[j].uid }

func (c ManifestsList) Sort() {
	sort.Sort(c)
}

func (c ManifestsList) String() string {
	var output []string
	output = append(output, "Manifests List:")
	for i := 0; i < len(c); i++ {
		if c[i] != nil {
			output = append(output, fmt.Sprintf("%v", c[i].String()))
		}
	}
	return strings.Join(output, "\n")
}

// Remember to sort before calling Sha if needed
func (c ManifestsList) Sha256() (result [sha256.Size]byte, err error) {
	if len(c) == 0 {
		return
	}

	runningHash := sha256.New()
	var emptyManifest CollectionsManifest
	emptyManifestBytes, _ := emptyManifest.Sha256()
	for i := 0; i < len(c); i++ {
		if c[i] == nil {
			runningHash.Write(emptyManifestBytes[0:len(emptyManifestBytes)])
		} else {
			var oneBytes [sha256.Size]byte
			oneBytes, err = c[i].Sha256()
			if err != nil {
				return
			}
			runningHash.Write(oneBytes[0:len(oneBytes)])
		}
	}

	tempSlice := runningHash.Sum(nil)
	if len(tempSlice) > sha256.Size {
		err = fmt.Errorf("Invalid sha256 hash - list: %v", c.String())
	}
	copy(result[:], tempSlice[:sha256.Size])
	return
}

type CollectionNamespaceList []*base.CollectionNamespace

func (c CollectionNamespaceList) Len() int      { return len(c) }
func (c CollectionNamespaceList) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c CollectionNamespaceList) Less(i, j int) bool {
	return (*(c[i])).LessThan(*(c[j]))
}

func SortCollectionsNamespaceList(list CollectionNamespaceList) CollectionNamespaceList {
	sort.Sort(list)
	return list
}

func (c CollectionNamespaceList) String() string {
	var buffer bytes.Buffer
	for _, j := range c {
		buffer.WriteString(fmt.Sprintf("|Scope: %v Collection: %v| ", j.ScopeName, j.CollectionName))
	}
	return buffer.String()
}

func (c CollectionNamespaceList) IsSame(other CollectionNamespaceList) bool {
	if len(c) != len(other) {
		return false
	}

	aList := SortCollectionsNamespaceList(c)
	bList := SortCollectionsNamespaceList(other)

	for i, col := range aList {
		if *col != *bList[i] {
			return false
		}
	}

	return true
}

func (c CollectionNamespaceList) Clone() (other CollectionNamespaceList) {
	for _, j := range c {
		ns := &base.CollectionNamespace{}
		*ns = *j
		other = append(other, ns)
	}
	return
}

func (c CollectionNamespaceList) Contains(namespace *base.CollectionNamespace) bool {
	if namespace == nil {
		return false
	}

	for _, j := range c {
		if *j == *namespace {
			return true
		}
	}
	return false
}

// Caller should have called IsSame() before doing consolidate
func (c *CollectionNamespaceList) Consolidate(other CollectionNamespaceList) {
	aMissingAction := func(item *base.CollectionNamespace) {
		*c = append(*c, item)
	}

	c.diffOrConsolidate(other, aMissingAction, nil /*bMissingAction*/)
}

func (c CollectionNamespaceList) Diff(other CollectionNamespaceList) (added, removed CollectionNamespaceList) {
	aMissingAction := func(item *base.CollectionNamespace) {
		added = append(added, item)
	}

	bMissingAction := func(item *base.CollectionNamespace) {
		removed = append(removed, item)
	}

	c.diffOrConsolidate(other, aMissingAction, bMissingAction)
	return
}

func (c *CollectionNamespaceList) diffOrConsolidate(other CollectionNamespaceList, aMissingAction, bMissingAction func(item *base.CollectionNamespace)) {
	aList := SortCollectionsNamespaceList(*c)
	bList := SortCollectionsNamespaceList(other)

	var aIdx int
	var bIdx int

	// Note - c == aList in this case

	for aIdx < len(aList) && bIdx < len(bList) {
		if aList[aIdx].IsSameAs(*(bList[bIdx])) {
			aIdx++
			bIdx++
		} else if aList[aIdx].LessThan(*(bList[bIdx])) {
			// Blist does not have something aList have.
			if bMissingAction != nil {
				bMissingAction(aList[aIdx])
			}
			aIdx++
		} else {
			// Blist[bIdx] < aList[aIdx]
			// B list has something aList does not have
			if aMissingAction != nil {
				aMissingAction(bList[bIdx])
			}
			bIdx++
		}
	}

	for bIdx < len(bList) {
		// The rest is all missing from A list
		if aMissingAction != nil {
			aMissingAction(bList[bIdx])
		}
		bIdx++
	}
}

type collectionNsMetaObj struct {
	SourceCollections CollectionNamespaceList `json:SourceCollections`
	// keys are integers of the index above
	IndirectTargetMap map[uint64]CollectionNamespaceList `json:IndirectTargetMap`
}

func newCollectionNsMetaObj() *collectionNsMetaObj {
	return &collectionNsMetaObj{
		IndirectTargetMap: make(map[uint64]CollectionNamespaceList),
	}
}

// This is used for namespace mapping that transcends over manifest lifecycles
// Need to use pointers because of golang hash map support of indexable type
// This means rest needs to do some gymanistics, instead of just simply checking for pointers
type CollectionNamespaceMapping map[*base.CollectionNamespace]CollectionNamespaceList

func (c *CollectionNamespaceMapping) MarshalJSON() ([]byte, error) {
	metaObj := newCollectionNsMetaObj()

	var unsortedKeys []*base.CollectionNamespace
	for k, _ := range *c {
		unsortedKeys = append(unsortedKeys, k)
	}
	sortedKeys := base.SortCollectionNamespacePtrList(unsortedKeys)

	for i, k := range sortedKeys {
		metaObj.SourceCollections = append(metaObj.SourceCollections, k)
		metaObj.IndirectTargetMap[uint64(i)] = (*c)[k]
	}

	return json.Marshal(metaObj)
}

func (c *CollectionNamespaceMapping) UnmarshalJSON(b []byte) error {
	if c == nil {
		return base.ErrorInvalidInput
	}

	metaObj := newCollectionNsMetaObj()

	err := json.Unmarshal(b, metaObj)
	if err != nil {
		return err
	}

	if (*c) == nil {
		(*c) = make(CollectionNamespaceMapping)
	}

	var i uint64
	for i = 0; i < uint64(len(metaObj.SourceCollections)); i++ {
		sourceCol := metaObj.SourceCollections[i]
		targetCols, ok := metaObj.IndirectTargetMap[i]
		if !ok {
			return fmt.Errorf("Unable to unmarshal CollectionNamespaceMapping raw: %v", metaObj)
		}
		(*c)[sourceCol] = targetCols
	}
	return nil
}

func (c *CollectionNamespaceMapping) String() string {
	var buffer bytes.Buffer
	for src, tgtList := range *c {
		buffer.WriteString(fmt.Sprintf("SOURCE ||Scope: %v Collection: %v|| -> TARGET(s) %v\n", src.ScopeName, src.CollectionName, CollectionNamespaceList(tgtList).String()))
	}
	return buffer.String()
}

func (c CollectionNamespaceMapping) Clone() (clone CollectionNamespaceMapping) {
	clone = make(CollectionNamespaceMapping)
	for k, v := range c {
		srcClone := &base.CollectionNamespace{}
		*srcClone = *k
		clone[srcClone] = v.Clone()
	}
	return
}

// The input "src" does not have to match the actual key pointer in the map, just the right matching values
// Returns the srcPtr for referring to the exact tgtList
func (c *CollectionNamespaceMapping) Get(src *base.CollectionNamespace) (srcPtr *base.CollectionNamespace, tgt CollectionNamespaceList, exists bool) {
	if src == nil {
		return
	}

	for k, v := range *c {
		if *k == *src {
			// found
			tgt = v
			srcPtr = k
			exists = true
			return
		}
	}
	return
}

func (c *CollectionNamespaceMapping) AddSingleMapping(src, tgt *base.CollectionNamespace) (alreadyExists bool) {
	if src == nil || tgt == nil {
		return
	}

	srcPtr, tgtList, found := c.Get(src)

	if !found {
		// Just use these as entries
		var newList CollectionNamespaceList
		newList = append(newList, tgt)
		(*c)[src] = newList
	} else {
		// See if tgt already exists in the current list
		if tgtList.Contains(tgt) {
			alreadyExists = true
			return
		}

		(*c)[srcPtr] = append((*c)[srcPtr], tgt)
	}
	return
}

// Given a scope and collection, see if it exists as one of the targets in the mapping
func (c *CollectionNamespaceMapping) TargetNamespaceExists(checkNamespace *base.CollectionNamespace) bool {
	if checkNamespace == nil {
		return false
	}
	for _, tgtList := range *c {
		if tgtList.Contains(checkNamespace) {
			return true
		}
	}
	return false
}

// Given a collection namespace mapping of source to target, and given a set of "ScopesMap",
// return a subset of collection namespace mapping of the original that contain the specified target scopesmap
func (c *CollectionNamespaceMapping) GetSubsetBasedOnSpecifiedTargets(targetScopeCollections ScopesMap) (retMap CollectionNamespaceMapping) {
	retMap = make(CollectionNamespaceMapping)

	for src, tgtList := range *c {
		for _, tgt := range tgtList {
			_, found := targetScopeCollections.GetCollectionByNames(tgt.ScopeName, tgt.CollectionName)
			if found {
				retMap.AddSingleMapping(src, tgt)
			}
		}
	}
	return
}

func (c CollectionNamespaceMapping) IsSame(other CollectionNamespaceMapping) bool {
	for src, tgtList := range c {
		_, otherTgtList, exists := other.Get(src)
		if !exists {
			return false
		}
		if !tgtList.IsSame(otherTgtList) {
			return false
		}
	}
	return true
}

func (c *CollectionNamespaceMapping) Delete(subset CollectionNamespaceMapping) (result CollectionNamespaceMapping) {
	// Instead of deleting, just make a brand new map
	result = make(CollectionNamespaceMapping)

	// Subtract B from A
	for aSrc, aTgtList := range *c {
		_, bTgtList, exists := subset.Get(aSrc)
		if !exists {
			// No need to delete
			result[aSrc] = aTgtList
			continue
		}
		if aTgtList.IsSame(bTgtList) {
			// The whole thing is deleted
			continue
		}
		// Gather the subset list
		var newList CollectionNamespaceList
		for _, ns := range aTgtList {
			if bTgtList.Contains(ns) {
				continue
			}
			newList = append(newList, ns)
		}
		result[aSrc] = newList
	}

	return
}

func (c *CollectionNamespaceMapping) Consolidate(other CollectionNamespaceMapping) {
	for otherSrc, otherTgtList := range other {
		srcPtr, tgtList, exists := c.Get(otherSrc)
		if !exists {
			(*c)[otherSrc] = otherTgtList.Clone()
		} else if !tgtList.IsSame(otherTgtList) {
			tgtList.Consolidate(otherTgtList)
			(*c)[srcPtr] = tgtList
		}
	}
}

func (c *CollectionNamespaceMapping) Diff(other CollectionNamespaceMapping) (added, removed CollectionNamespaceMapping) {
	added = make(CollectionNamespaceMapping)
	removed = make(CollectionNamespaceMapping)
	// First, populated "removed"
	for src, tgtList := range *c {
		_, oTgtList, exists := other.Get(src)
		if !exists {
			removed[src] = tgtList
		} else if !tgtList.IsSame(oTgtList) {
			listAdded, listRemoved := tgtList.Diff(oTgtList)
			for _, addedNamespace := range listAdded {
				added.AddSingleMapping(src, addedNamespace)
			}
			for _, removedNamespace := range listRemoved {
				removed.AddSingleMapping(src, removedNamespace)
			}
		}
	}

	// Then populate added
	for oSrc, oTgtList := range other {
		_, _, exists := c.Get(oSrc)
		if !exists {
			added[oSrc] = oTgtList
			// No else - any potential intersections would have been captured above
		}
	}
	return
}

// Json marshaller will serialize the map by key, but not necessarily the values, which is ordered list
// Because the lists may not be ordered, we need to calculate sha256 with lists ordered
func (c *CollectionNamespaceMapping) Sha256() (result [sha256.Size]byte, err error) {
	if c == nil {
		err = base.ErrorInvalidInput
		return
	}

	// Simpler to just create a temporary map with ordered list for sha calculation
	tempMap := make(CollectionNamespaceMapping)

	for k, v := range *c {
		tempMap[k] = SortCollectionsNamespaceList(v)
	}

	marshalledJson, err := tempMap.MarshalJSON()
	if err != nil {
		return
	}

	result = sha256.Sum256(marshalledJson)
	return
}

func (c *CollectionNamespaceMapping) ToSnappyCompressed() ([]byte, error) {
	marshalledJson, err := c.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return snappy.Encode(nil, marshalledJson), nil
}

func NewCollectionNamespaceMappingFromSnappyData(data []byte) (*CollectionNamespaceMapping, error) {
	marshalledJson, err := snappy.Decode(nil, data)
	if err != nil {
		return nil, err
	}

	newMap := make(CollectionNamespaceMapping)
	err = newMap.UnmarshalJSON(marshalledJson)

	return &newMap, err
}

type CollectionNsMappingsDoc struct {
	NsMappingRecords CompressedColNamespaceMappingList `json:"NsMappingRecords"`

	// internal id of repl spec - for detection of repl spec deletion and recreation event
	SpecInternalId string `json:"specInternalId"`

	//revision number
	revision interface{}
}

func (b *CollectionNsMappingsDoc) Size() int {
	if b == nil {
		return 0
	}

	return len(b.SpecInternalId) + b.NsMappingRecords.Size()
}

func (b *CollectionNsMappingsDoc) ToShaMap() (ShaToCollectionNamespaceMap, error) {
	if b == nil {
		return nil, base.ErrorInvalidInput
	}

	errorMap := make(base.ErrorMap)
	shaMap := make(ShaToCollectionNamespaceMap)

	for _, oneRecord := range b.NsMappingRecords {
		if oneRecord == nil {
			continue
		}

		serializedMap, err := snappy.Decode(nil, oneRecord.CompressedMapping)
		if err != nil {
			errorMap[oneRecord.Sha256Digest] = fmt.Errorf("Snappy decompress failed %v", err)
			continue
		}
		actualMap := make(CollectionNamespaceMapping)
		err = json.Unmarshal(serializedMap, &actualMap)
		if err != nil {
			errorMap[oneRecord.Sha256Digest] = fmt.Errorf("Unmarshalling failed %v", err)
			continue
		}
		// Sanity check
		checkSha, err := actualMap.Sha256()
		if err != nil {
			errorMap[oneRecord.Sha256Digest] = fmt.Errorf("Validing SHA failed %v", err)
			continue
		}
		checkShaString := fmt.Sprintf("%x", checkSha[:])
		if checkShaString != oneRecord.Sha256Digest {
			errorMap[oneRecord.Sha256Digest] = fmt.Errorf("SHA validation mismatch %v", checkShaString)
			continue
		}

		shaMap[oneRecord.Sha256Digest] = &actualMap
	}

	var err error
	if len(errorMap) > 0 {
		errStr := base.FlattenErrorMap(errorMap)
		err = fmt.Errorf(errStr)
	}
	return shaMap, err
}

// Will overwrite the existing records with the incoming map
func (b *CollectionNsMappingsDoc) LoadShaMap(shaMap ShaToCollectionNamespaceMap) error {
	if b == nil {
		return base.ErrorInvalidInput
	}

	errorMap := make(base.ErrorMap)
	b.NsMappingRecords = b.NsMappingRecords[:0]

	for sha, colNsMap := range shaMap {
		if colNsMap == nil {
			continue
		}
		// do sanity check - TODO remove before CC is shipped
		checkSha, err := colNsMap.Sha256()
		if err != nil {
			panic(err)
			continue
		}
		checkShaStr := fmt.Sprintf("%x", checkSha[:])
		if sha != checkShaStr {
			errorMap[sha] = fmt.Errorf("Before storing, sanity check failed - expected %v got %v", sha, checkShaStr)
		}

		compressedMapping, err := colNsMap.ToSnappyCompressed()
		if err != nil {
			errorMap[sha] = err
			continue
		}

		oneRecord := &CompressedColNamespaceMapping{compressedMapping, sha}
		b.NsMappingRecords.SortedInsert(oneRecord)
	}

	if len(errorMap) > 0 {
		return fmt.Errorf("Error LoadingShaMap - sha -> err: %v", base.FlattenErrorMap(errorMap))
	} else {
		return nil
	}
}

type ShaToCollectionNamespaceMap map[string]*CollectionNamespaceMapping

func (s *ShaToCollectionNamespaceMap) Clone() (newMap ShaToCollectionNamespaceMap) {
	if s == nil {
		return
	}

	newMap = make(ShaToCollectionNamespaceMap)

	for k, v := range *s {
		clonedVal := v.Clone()
		newMap[k] = &clonedVal
	}
	return
}

func (s *ShaToCollectionNamespaceMap) String() string {
	if s == nil {
		return "<nil>"
	}

	var output []string
	for k, v := range *s {
		output = append(output, fmt.Sprintf("Sha256Digest: %v Map:", k))
		if v != nil {
			output = append(output, v.String())
		}
	}

	return strings.Join(output, "\n")
}

func (s ShaToCollectionNamespaceMap) Diff(older ShaToCollectionNamespaceMap) (added, removed ShaToCollectionNamespaceMap) {
	if len(older) == 0 && len(s) > 0 {
		added = s
		return
	} else if len(older) > 0 && len(s) == 0 {
		removed = s
		return
	}

	added = make(ShaToCollectionNamespaceMap)
	removed = make(ShaToCollectionNamespaceMap)

	for k, v := range s {
		if _, exists := older[k]; !exists {
			added[k] = v
		}
	}

	for k, v := range older {
		if _, exists := s[k]; !exists {
			removed[k] = v
		}
	}

	return
}

type CompressedColNamespaceMapping struct {
	// Snappy compressed byte slice of CollectionNamespaceMapping
	CompressedMapping []byte `json:compressedMapping`
	Sha256Digest      string `json:string`
}

func (c *CompressedColNamespaceMapping) String() string {
	return fmt.Sprintf("Sha: %v Bytes: %v", c.Sha256Digest, fmt.Sprintf("%x", c.CompressedMapping[:]))
}

func (c *CompressedColNamespaceMapping) Size() int {
	if c == nil {
		return 0
	}
	return len(c.CompressedMapping) + len(c.Sha256Digest)
}

type CompressedColNamespaceMappingList []*CompressedColNamespaceMapping

func (c *CompressedColNamespaceMappingList) Size() int {
	if c == nil {
		return 0
	}

	var totalSize int
	for _, j := range *c {
		totalSize += j.Size()
	}
	return totalSize
}

func (c *CompressedColNamespaceMappingList) SortedInsert(elem *CompressedColNamespaceMapping) {
	if c == nil {
		return
	}

	if len(*c) == 0 {
		*c = append(*c, elem)
		return
	}

	var i int
	// First, find where this should be
	for i = 0; i < len(*c); i++ {
		if (*c)[i].Sha256Digest > elem.Sha256Digest {
			break
		}
	}

	*c = append(*c, nil)
	copy((*c)[i+1:], (*c)[i:])
	(*c)[i] = elem
}
