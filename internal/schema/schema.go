// Copyright 2019 The Meerkat Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package schema

import (
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"meerkat/internal/cluster"
	"sort"
	"sync"
	"time"
)

//go:generate protoc  -I . -I ../../build/proto/   --plugin ../../build/protoc-gen-gogofaster --gogofaster_out=plugins=grpc:. ./schema.proto

const (
	eventHandlerID = "schema_event_handler"

	indexMapName  = "index"
	fieldsMapName = "fields"
	pAllocMapName = "alloc"

	TSFieldName = "_ts"
	IDFieldName = "_id"
)

type OpType int

const (
	OpIndexCreate OpType = iota
	OpIndexUpdate
	OpIndexDelete
)

var internalFields = map[string]bool{
	TSFieldName: true,
	IDFieldName: true,
}

type IndexInfo struct {
	Id             string         `json:"id"`
	Name           string         `json:"name"`
	Desc           string         `json:"desc"`
	Created        time.Time      `json:"created"`
	Updated        time.Time      `json:"updated"`
	Fields         []Field        `json:"fields"` // TODO(gvelo): remove the array and keep fields info indexed on maps.
	PartitionAlloc PartitionAlloc `json:"partitionAlloc"`
	fieldsByID     map[string]Field
	fieldsByName   map[string]Field
}

func (ii *IndexInfo) addField(f Field) {
	ii.fieldsByID[f.Id] = f
	ii.fieldsByName[f.Name] = f
	ii.Fields = ii.fields()
}

func (ii *IndexInfo) removeField(id string) {
	field, found := ii.fieldsByID[id]
	if found {
		delete(ii.fieldsByID, id)
		delete(ii.fieldsByName, field.Name)
		ii.Fields = ii.fields()
	}
}

func (ii *IndexInfo) FieldByName(name string) (Field, error) {

	field, ok := ii.fieldsByName[name]

	if !ok {
		return Field{}, &NotFoundError{Err: fmt.Sprintf("field with name %v cannot be found", name)}
	}

	return field, nil
}

func (ii *IndexInfo) fields() []Field {
	fields := make([]Field, len(ii.fieldsByID))
	i := 0
	for _, f := range ii.fieldsByID {
		fields[i] = f
		i++
	}
	return fields
}

//TODO(gvelo): remove this method and use custom json serialization.
func (ii *IndexInfo) init() {
	ii.fieldsByID = make(map[string]Field)
	ii.fieldsByName = make(map[string]Field)
	for _, f := range ii.Fields {
		ii.fieldsByID[f.Id] = f
		ii.fieldsByName[f.Name] = f
	}
}

//TODO(gvelo): reflection.deepCopy()
func (ii IndexInfo) copy() IndexInfo {
	f := make([]Field, len(ii.Fields))
	copy(f, ii.Fields)
	ii.Fields = f
	ii.init()
	return ii
}

func (ii *IndexInfo) validate() error {

	if ii.Name == "" {
		return &ValidationError{
			Err:   "index name cannot be empty",
			Field: "name",
		}
	}

	// TODO(gvelo) move this validation to the new json custom unmarshall.
	fieldsByName := make(map[string]Field)

	for _, field := range ii.Fields {

		err := field.validate()

		if err != nil {
			return err
		}

		_, ok := fieldsByName[field.Name]

		if ok {
			return &ValidationError{
				Err:   fmt.Sprintf("field with name %v already exist", field.Name),
				Field: "field.name",
			}
		}

		fieldsByName[field.Name] = field

	}

	return nil

}

func (ii *IndexInfo) addInternalFields(now time.Time) {

	ts := Field{
		Id:        TSFieldName,
		Name:      TSFieldName,
		Desc:      "TS Field",
		IndexId:   ii.Id,
		FieldType: FieldType_TIMESTAMP,
		Nullable:  false,
		Created:   now,
		Updated:   now,
	}

	id := Field{
		Id:        IDFieldName,
		Name:      IDFieldName,
		Desc:      "id",
		IndexId:   ii.Id,
		FieldType: FieldType_UUID,
		Nullable:  false,
		Created:   now,
		Updated:   now,
	}

	ii.addField(ts)
	ii.addField(id)

}

func (f *Field) validate() error {

	if f.Name == "" {
		return &ValidationError{
			Err:   "field name cannot be empty",
			Field: "name",
		}
	}

	return nil

}

func (ft FieldType) MarshalJSON() ([]byte, error) {
	v := FieldType_name[int32(ft)]
	return json.Marshal(v)
}

func (ft *FieldType) UnmarshalJSON(b []byte) error {

	var s string

	err := json.Unmarshal(b, &s)

	if err != nil {
		return err
	}

	v, found := FieldType_value[s]

	if !found {
		return fmt.Errorf("invalid FieldType value [%v]", s)
	}

	*ft = FieldType(v)

	return nil
}

func (p *PartitionAlloc) validate() error {
	// TODO: validate partition allocation.
	return nil
}

type ValidationError struct {
	Err   string
	Field string
}

func (e *ValidationError) Error() string {
	return e.Err
}

type NotFoundError struct {
	Err string
}

func (e *NotFoundError) Error() string {
	return e.Err
}

type Schema interface {
	AllIndex() []IndexInfo
	Index(id string) (IndexInfo, error)
	IndexByName(name string) (IndexInfo, error)
	CreateIndex(index IndexInfo) (IndexInfo, error)
	UpdateIndex(index IndexInfo) (IndexInfo, error)
	DeleteIndex(id string) error

	AllFields(id string) ([]Field, error)
	Field(id string) (Field, error)
	UpdateField(field Field) error
	DeleteField(id string) error
	CreateFields(id string, fields Field) (Field, error)

	Alloc(id string) (PartitionAlloc, error)
	UpdateAlloc(id string, parAlloc PartitionAlloc) error

	AddEventHandler(id string, h chan IndexUpdateEvent)
	RemoveEventHandler(id string)

	Shutdown()
}

type IndexUpdateEvent struct {
	OpType    OpType
	IndexInfo IndexInfo
}

type schema struct {
	catalog       cluster.Catalog
	indexCache    map[string]IndexInfo
	fieldCache    map[string]Field
	pAllocCache   map[string]PartitionAlloc
	indexByName   map[string]IndexInfo
	log           zerolog.Logger
	catalogCh     chan []cluster.Entry
	mu            sync.Mutex
	done          chan struct{}
	eventHandlers map[string]chan IndexUpdateEvent
}

func (s *schema) AllIndex() []IndexInfo {

	s.mu.Lock()
	defer s.mu.Unlock()

	indexes := make([]IndexInfo, len(s.indexCache))

	i := 0

	for _, index := range s.indexCache {
		// TODO(gvelo) instead of copy-on-read use copy-on-write
		indexes[i] = index.copy()
		i++
	}

	return indexes

}

func (s *schema) Index(id string) (IndexInfo, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	indexInfo, ok := s.indexCache[id]

	if !ok {
		return IndexInfo{}, &NotFoundError{Err: fmt.Sprintf("index %v cannot be found", id)}
	}

	return indexInfo.copy(), nil
}

func (s *schema) IndexByName(name string) (IndexInfo, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	indexInfo, ok := s.indexByName[name]

	if !ok {
		return IndexInfo{}, &NotFoundError{Err: fmt.Sprintf("index %v cannot be found", name)}
	}

	return indexInfo.copy(), nil
}

func (s *schema) CreateIndex(indexInfo IndexInfo) (IndexInfo, error) {

	// TODO(gvelo): storing more entries than the catalog emit channel capacity will
	//  deadlock since the gorutine that read from catalog use the same lock.
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	indexInfo = indexInfo.copy()

	if indexInfo.Id != "" {
		return IndexInfo{}, &ValidationError{
			Err:   "index id should not be set for new index",
			Field: "id",
		}
	}

	dup, ok := s.indexByName[indexInfo.Name]

	if ok && dup.Id != indexInfo.Id {
		return IndexInfo{}, &ValidationError{
			Err:   fmt.Sprintf("index with name [%v] already exist", indexInfo.Name),
			Field: "name",
		}
	}

	err := indexInfo.validate()

	if err != nil {
		return IndexInfo{}, err
	}

	indexInfo.Id = uuid.New().String()
	indexInfo.Created = now
	indexInfo.Updated = now

	for i, f := range indexInfo.Fields {
		if f.Id != "" {
			return IndexInfo{}, &ValidationError{
				Err:   "field id should not be set for new index",
				Field: "id",
			}
		}
		// TODO(gvelo) could we use a binary representation of the UUID ? ( []byte ?)
		f.Id = uuid.New().String()
		f.IndexId = indexInfo.Id
		f.Created = now
		f.Updated = now
		indexInfo.Fields[i] = f
	}

	// add internal fields
	indexInfo.init()
	indexInfo.addInternalFields(now)

	entries, err := createEntries(indexInfo, now)

	s.addToCache(indexInfo)

	indexInfo = indexInfo.copy()

	s.emitEvent(IndexUpdateEvent{OpIndexCreate, indexInfo})

	s.catalog.SetAll(entries)

	return indexInfo, nil

}

func (s *schema) addToCache(indexInfo IndexInfo) {
	indexInfo.init()
	s.indexCache[indexInfo.Id] = indexInfo
	s.indexByName[indexInfo.Name] = indexInfo
	for _, f := range indexInfo.Fields {
		s.fieldCache[f.Id] = f
	}
}

func (s *schema) UpdateIndex(index IndexInfo) (IndexInfo, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	indexInfo, ok := s.indexCache[index.Id]

	if !ok {
		return IndexInfo{}, &NotFoundError{
			Err: fmt.Sprintf("index %v does't exist", index.Id),
		}
	}

	err := index.validate()

	if err != nil {
		return IndexInfo{}, err
	}

	dup, ok := s.indexByName[index.Name]

	if ok && dup.Id != index.Id {
		return IndexInfo{}, &ValidationError{
			Err:   fmt.Sprintf("index with name [%v] already exist", index.Name),
			Field: "name",
		}
	}

	if indexInfo.PartitionAlloc.NumOfPartitions != index.PartitionAlloc.NumOfPartitions {
		return IndexInfo{}, &ValidationError{
			Err:   "NumOfPartitions cannot be updated",
			Field: "PartitionAlloc.NumOfPartitions",
		}
	}

	index.init()

	// TODO(gvelo) validate internal fields
	// delete fields
	now := time.Now()

	var fieldsToRemove []Field

	for id, f := range indexInfo.fieldsByID {
		_, ok := index.fieldsByID[id]
		if !ok {
			f.Updated = now
			fieldsToRemove = append(fieldsToRemove, f)
		}
	}

	delEntries, err := createFieldEntries(fieldsToRemove, now)

	if err != nil {
		return IndexInfo{}, err
	}

	for i := range delEntries {
		delEntries[i].Deleted = true
	}

	// add new fields

	for i, f := range index.Fields {
		f.Updated = now
		f.IndexId = indexInfo.Id
		if f.Id != "" {
			_, ok := s.fieldCache[f.Id]
			if !ok {
				return IndexInfo{}, &NotFoundError{
					Err: fmt.Sprintf("field %v doesn't exist", f.Id),
				}
			}
			continue
		}
		f.Id = uuid.New().String()
		index.Fields[i] = f
	}

	index.Updated = now

	entries, err := createEntries(index, now)

	if err != nil {
		return IndexInfo{}, err
	}

	entries = append(entries, delEntries...)

	for _, f := range fieldsToRemove {
		delete(s.fieldCache, f.Id)
	}

	// remove the old index name
	delete(s.indexByName, indexInfo.Name)
	s.addToCache(index)

	index = index.copy()

	s.emitEvent(IndexUpdateEvent{OpIndexUpdate, index})

	s.catalog.SetAll(entries)

	return index, nil
}

func (s *schema) DeleteIndex(id string) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	indexInfo, ok := s.indexCache[id]

	if !ok {
		return &NotFoundError{
			Err: fmt.Sprintf("index %v does't exist", id),
		}
	}

	now := time.Now()

	indexInfo.Updated = now

	for i := range indexInfo.Fields {
		indexInfo.Fields[i].Updated = now
	}

	entries, err := createEntries(indexInfo, now)

	if err != nil {
		return err
	}

	for i := range entries {
		entries[i].Deleted = true
	}

	delete(s.indexCache, indexInfo.Id)
	delete(s.indexByName, indexInfo.Name)

	for _, f := range indexInfo.Fields {
		delete(s.fieldCache, f.Id)
	}

	s.emitEvent(IndexUpdateEvent{OpIndexDelete, indexInfo.copy()})

	s.catalog.SetAll(entries)

	return nil

}

func createEntries(indexInfo IndexInfo, t time.Time) ([]cluster.Entry, error) {

	var entries []cluster.Entry

	indexEntry, err := createIndexEntry(indexInfo, t)

	if err != nil {
		return nil, err
	}

	fieldEntries, err := createFieldEntries(indexInfo.Fields, t)

	if err != nil {
		return nil, err
	}

	pAllocEntry, err := createPAllocEntry(indexInfo.PartitionAlloc, indexInfo.Id, t)

	entries = append(entries, indexEntry)
	entries = append(entries, fieldEntries...)
	entries = append(entries, pAllocEntry)

	return entries, nil
}

func createIndexEntry(indexInfo IndexInfo, t time.Time) (cluster.Entry, error) {

	index := &Index{
		Id:      indexInfo.Id,
		Name:    indexInfo.Name,
		Desc:    indexInfo.Desc,
		Created: indexInfo.Created,
		Updated: indexInfo.Updated,
	}

	indexBytes, err := proto.Marshal(index)

	if err != nil {
		return cluster.Entry{}, err
	}

	indexEntry := cluster.Entry{
		MapName: indexMapName,
		Key:     index.Id,
		Value:   indexBytes,
		Time:    t,
	}

	return indexEntry, nil

}

func createFieldEntries(fields []Field, t time.Time) ([]cluster.Entry, error) {

	var entries []cluster.Entry

	for _, field := range fields {

		fieldBytes, err := proto.Marshal(&field)

		if err != nil {
			return nil, err
		}

		fieldEntry := cluster.Entry{
			MapName: fieldsMapName,
			Key:     field.Id,
			Value:   fieldBytes,
			Time:    t,
		}

		entries = append(entries, fieldEntry)

	}

	return entries, nil
}

func createPAllocEntry(pAlloc PartitionAlloc, id string, t time.Time) (cluster.Entry, error) {

	pAllocBytes, err := proto.Marshal(&pAlloc)

	if err != nil {
		return cluster.Entry{}, err
	}

	partitionEntry := cluster.Entry{
		MapName: pAllocMapName,
		Key:     id,
		Value:   pAllocBytes,
		Time:    t,
	}

	return partitionEntry, nil

}

func (s *schema) AllFields(id string) ([]Field, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	index, ok := s.indexCache[id]

	if !ok {
		return nil, &NotFoundError{
			Err: fmt.Sprintf("index %v doesnt exist", id),
		}
	}

	r := make([]Field, len(index.Fields))
	copy(r, index.Fields)

	return r, nil

}

func (s *schema) Field(id string) (Field, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	f, ok := s.fieldCache[id]

	if !ok {
		return Field{}, &NotFoundError{
			Err: fmt.Sprintf("field %v doesnt exist", id),
		}
	}

	return f, nil

}

func (s *schema) UpdateField(field Field) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	f, ok := s.fieldCache[field.Id]

	if !ok {
		return &NotFoundError{
			Err: fmt.Sprintf("field %v doesn't exist", field.Id),
		}
	}

	i, ok := s.indexCache[f.IndexId]

	if !ok {
		return &NotFoundError{
			Err: fmt.Sprintf("field %v doesn't exist", field.Id),
		}
	}

	field.IndexId = f.IndexId
	field.Created = f.Created
	field.Updated = time.Now()

	i = i.copy()
	i.addField(field)

	err := i.validate()

	if err != nil {
		return err
	}

	fieldBytes, err := proto.Marshal(&field)

	if err != nil {
		return err
	}

	fieldEntry := cluster.Entry{
		MapName: fieldsMapName,
		Key:     field.Id,
		Value:   fieldBytes,
		Time:    field.Updated,
	}

	s.fieldCache[field.Id] = field
	s.indexCache[i.Id] = i

	s.emitEvent(IndexUpdateEvent{OpIndexUpdate, i})

	s.catalog.Set(fieldEntry)

	return nil

}

func (s *schema) DeleteField(id string) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	f, ok := s.fieldCache[id]

	if !ok {
		return &NotFoundError{
			Err: fmt.Sprintf("field %v doesn't exist", id),
		}
	}

	i, ok := s.indexCache[f.IndexId]

	if !ok {
		return &NotFoundError{
			Err: fmt.Sprintf("field %v doesn't exist", id),
		}
	}

	f.Updated = time.Now()

	i.removeField(id)

	//TODO(gvelo):use createEntries()

	fieldBytes, err := proto.Marshal(&f)

	if err != nil {
		return err
	}

	fieldEntry := cluster.Entry{
		MapName: fieldsMapName,
		Key:     f.Id,
		Value:   fieldBytes,
		Time:    f.Updated,
		Deleted: true,
	}

	delete(s.fieldCache, f.Id)

	s.emitEvent(IndexUpdateEvent{OpIndexUpdate, i.copy()})

	s.catalog.Set(fieldEntry)

	return nil

}

func (s *schema) CreateFields(id string, field Field) (Field, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	i, ok := s.indexCache[id]

	if !ok {
		return Field{}, &NotFoundError{
			Err: fmt.Sprintf("index %v doesn't exist", id),
		}
	}

	now := time.Now()

	field.Id = uuid.New().String()
	field.IndexId = id
	field.Created = now
	field.Updated = now

	i = i.copy()
	i.addField(field)

	err := i.validate()

	if err != nil {
		return Field{}, err
	}

	// TODO:(gvelo): use createEntries() ?

	fieldBytes, err := proto.Marshal(&field)

	if err != nil {
		return Field{}, err
	}

	entry := cluster.Entry{
		MapName: fieldsMapName,
		Key:     field.Id,
		Time:    field.Created,
		Value:   fieldBytes,
	}

	s.fieldCache[field.Id] = field
	s.indexCache[i.Id] = i
	s.indexByName[i.Name] = i

	s.emitEvent(IndexUpdateEvent{OpIndexUpdate, i})

	s.catalog.Set(entry)

	return field, nil
}

func (s *schema) Alloc(id string) (PartitionAlloc, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	index, ok := s.indexCache[id]

	if !ok {
		return PartitionAlloc{}, &NotFoundError{
			Err: fmt.Sprintf("index %v doesn't exist", id),
		}
	}

	return index.PartitionAlloc, nil

}

func (s *schema) UpdateAlloc(id string, pAlloc PartitionAlloc) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	index, ok := s.indexCache[id]

	if !ok {
		return &NotFoundError{
			Err: fmt.Sprintf("index %v doesn't exist", id),
		}
	}

	err := pAlloc.validate()

	if err != nil {
		return err
	}

	if index.PartitionAlloc.NumOfPartitions != pAlloc.NumOfPartitions {
		return &ValidationError{
			Err:   "NumOfPartitions cannot be updated",
			Field: "PartitionAlloc.NumOfPartitions",
		}
	}

	now := time.Now()

	pAlloc.Updated = now

	// TODO(gvelo): move to createEntry()

	pAllocBytes, err := proto.Marshal(&pAlloc)

	if err != nil {
		return err
	}

	entry := cluster.Entry{
		MapName: pAllocMapName,
		Key:     id,
		Value:   pAllocBytes,
		Time:    now,
	}

	index.PartitionAlloc = pAlloc
	s.indexCache[index.Id] = index

	s.emitEvent(IndexUpdateEvent{OpIndexUpdate, index.copy()})

	s.catalog.Set(entry)

	return nil

}

func NewSchema(catalog cluster.Catalog) (Schema, error) {

	s := &schema{
		catalog:       catalog,
		indexCache:    make(map[string]IndexInfo),
		fieldCache:    make(map[string]Field),
		pAllocCache:   make(map[string]PartitionAlloc),
		indexByName:   make(map[string]IndexInfo),
		log:           log.With().Str("component", "schema").Logger(),
		catalogCh:     make(chan []cluster.Entry, 1024),
		done:          make(chan struct{}),
		eventHandlers: make(map[string]chan IndexUpdateEvent),
	}

	catalog.AddEventHandler(eventHandlerID, s.catalogCh)

	err := s.load()

	if err != nil {
		catalog.RemoveEventHandler(eventHandlerID)
		return nil, err
	}

	go s.sync()

	return s, nil

}

func (s *schema) Shutdown() {
	s.catalog.RemoveEventHandler(eventHandlerID)
	close(s.done)
}

func (s *schema) load() error {

	//TODO(gvelo): read from index and field buckets inside a transaction.

	fieldsByIdx := make(map[string][]Field)

	fieldEntries := s.catalog.GetAll(fieldsMapName)

	// load fields
	for _, fe := range fieldEntries {

		var field Field

		err := proto.Unmarshal(fe.Value, &field)

		if err != nil {
			s.log.Panic().Err(err).Msg("error unmarshaling field info from catalog")
		}

		fieldsByIdx[field.IndexId] = append(fieldsByIdx[field.IndexId], field)
		s.fieldCache[field.Id] = field

	}

	// load index
	idxEntries := s.catalog.GetAll(indexMapName)

	for _, e := range idxEntries {

		var index Index

		err := proto.Unmarshal(e.Value, &index)

		if err != nil {
			return err
		}

		indexInfo := IndexInfo{
			Id:      index.Id,
			Name:    index.Name,
			Desc:    index.Desc,
			Fields:  fieldsByIdx[index.Id],
			Created: index.Created,
			Updated: index.Updated,
		}

		allocEntry, ok := s.catalog.Get(pAllocMapName, index.Id)

		if ok {
			err := proto.Unmarshal(allocEntry.Value, &indexInfo.PartitionAlloc)
			if err != nil {
				s.log.Panic().Err(err).Msg("error unmarshaling Partition alloc info")
			}
		}

		indexInfo.init()
		s.indexCache[e.Key] = indexInfo
		s.indexByName[indexInfo.Name] = indexInfo
	}
	return nil
}

func (s *schema) sync() {

	s.log.Info().Msg("starting schema synchronization")

	for {
		select {
		case <-s.done:
			s.log.Info().Msg("stopping schema synchronization")
			return
		case delta := <-s.catalogCh:
			events := s.syncDelta(delta)
			s.emitEvents(events)
		}
	}

}

func (s *schema) syncDelta(delta []cluster.Entry) map[string]IndexUpdateEvent {

	s.mu.Lock()
	defer s.mu.Unlock()

	events := make(map[string]IndexUpdateEvent)

	// update all fields and Pallocs.
	for _, e := range delta {

		var index IndexInfo

		switch e.MapName {
		case fieldsMapName:
			index = s.updateField(e)
		case pAllocMapName:
			index = s.updatePAlloc(e)
		}

		if index.Id != "" {
			events[index.Id] = IndexUpdateEvent{
				OpType:    OpIndexUpdate,
				IndexInfo: index.copy(),
			}
		}

	}

	// update all index.
	for _, e := range delta {

		if e.MapName == indexMapName {

			index, op := s.updateIndex(e)

			if index.Id != "" {
				events[index.Id] = IndexUpdateEvent{
					OpType:    op,
					IndexInfo: index.copy(),
				}
			}

		}

	}

	return events

}

func (s *schema) updateField(e cluster.Entry) IndexInfo {

	var field Field

	err := proto.Unmarshal(e.Value, &field)

	if err != nil {
		s.log.Panic().Err(err).Msg("Cannot unmarshal field value")
	}

	if e.Deleted {

		delete(s.fieldCache, field.Id)

		// index map and field map are eventually consistent
		// so there is a chance of a missing index.
		index, ok := s.indexCache[field.IndexId]

		if ok {

			_, ok := index.fieldsByID[field.Id]

			// the field has been already removed.
			if !ok {
				return IndexInfo{}
			}

			index.removeField(field.Id)
			s.indexCache[index.Id] = index
			return index
		}

		return IndexInfo{}
	}

	s.fieldCache[field.Id] = field

	index, ok := s.indexCache[field.IndexId]

	if ok {

		f, ok := index.fieldsByID[field.Id]

		if ok && f.Updated.Equal(field.Updated) {
			return IndexInfo{}
		}

		index.addField(field)
		s.indexCache[index.Id] = index
		return index

	}

	return IndexInfo{}

}

func (s *schema) updatePAlloc(e cluster.Entry) IndexInfo {

	var pAlloc PartitionAlloc

	err := proto.Unmarshal(e.Value, &pAlloc)

	if err != nil {
		s.log.Panic().Err(err).Msg("Cannot unmarshal pAlloc value")
	}

	// TODO(gvelo): handle alloc delete properly.

	s.pAllocCache[e.Key] = pAlloc

	index, ok := s.indexCache[e.Key]

	if ok {

		if index.PartitionAlloc.Updated.Equal(pAlloc.Updated) {
			return IndexInfo{}
		}

		index.PartitionAlloc = pAlloc
		s.indexCache[index.Id] = index
		return index
	}

	return IndexInfo{}

}

func (s *schema) updateIndex(e cluster.Entry) (IndexInfo, OpType) {

	var index Index

	err := proto.Unmarshal(e.Value, &index)

	if err != nil {
		s.log.Panic().Err(err).Msg("cannot unmarshall index catalog entry")
	}

	if e.Deleted {

		i, found := s.indexCache[index.Id]

		if found {
			delete(s.indexCache, index.Id)
			delete(s.indexByName, index.Name)
			return i, OpIndexDelete
		}

		return IndexInfo{}, -1

	}

	indexInfo := IndexInfo{
		Id:             index.Id,
		Name:           index.Name,
		Desc:           index.Desc,
		Created:        index.Created,
		Updated:        index.Updated,
		PartitionAlloc: s.pAllocCache[e.Key],
	}

	for _, f := range s.fieldCache {
		if f.IndexId == indexInfo.Id {
			indexInfo.Fields = append(indexInfo.Fields, f)
		}
	}

	indexInfo.init()

	op := OpIndexCreate

	if i, found := s.indexCache[indexInfo.Id]; found {

		if i.Updated.Equal(index.Updated) {
			return indexInfo, -1
		}

		op = OpIndexUpdate
	}

	s.indexCache[indexInfo.Id] = indexInfo
	s.indexByName[indexInfo.Name] = indexInfo

	return indexInfo, op

}

func (s *schema) AddEventHandler(id string, h chan IndexUpdateEvent) {

	s.mu.Lock()
	defer s.mu.Unlock()

	s.eventHandlers[id] = h

}

func (s *schema) RemoveEventHandler(id string) {

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.eventHandlers, id)

}

func (s *schema) emitEvent(event IndexUpdateEvent) {

	for id, h := range s.eventHandlers {
		select {
		case h <- event:
			s.log.Debug().Msgf("dispatching event to %v", id)
		default:
			s.log.Error().Msgf("dispatcher blocks on event handler channel [%v]", id)
			h <- event
		}
	}

}

func (s *schema) emitEvents(events map[string]IndexUpdateEvent) {

	s.mu.Lock()
	defer s.mu.Unlock()

	var sortedEvt []IndexUpdateEvent

	for _, event := range events {
		sortedEvt = append(sortedEvt, event)
	}

	// dispatch event in order (create/update/delete)
	sort.Slice(sortedEvt, func(i, j int) bool {
		return sortedEvt[i].OpType < sortedEvt[j].OpType
	})

	for _, se := range sortedEvt {
		s.emitEvent(se)
	}

}
