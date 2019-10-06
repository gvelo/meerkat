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
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"meerkat/internal/cluster"
	"sync"
	"time"
)

//go:generate protoc  -I . -I ../../build/proto/   --plugin ../../build/protoc-gen-gogofaster --gogofaster_out=plugins=grpc:. ./schema.proto

const (
	eventHandlerID = "schema_event_handler"

	indexMapName  = "index"
	fieldsMapName = "fields"
	pAllocMapName = "alloc"
)

type IndexInfo struct {
	Id             string         `json:"id"`
	Name           string         `json:"name"`
	Desc           string         `json:"desc"`
	Created        time.Time      `json:"created"`
	Updated        time.Time      `json:"updated"`
	Fields         []Field        `json:"fields"`
	PartitionAlloc PartitionAlloc `json:"partitionAlloc"`
	fieldsByID     map[string]Field
}

// TODO(gvelo): this is quite confusing. we should move this mapping
//              to a new field by index cache.
func (ii *IndexInfo) addField(f Field) {
	ii.fieldsByID[f.Id] = f
	ii.Fields = ii.fields()
}

func (ii *IndexInfo) removeField(id string) {
	delete(ii.fieldsByID, id)
	ii.Fields = ii.fields()
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

func (ii *IndexInfo) init() {
	ii.fieldsByID = make(map[string]Field)
	for _, f := range ii.Fields {
		ii.fieldsByID[f.Id] = f
	}
}

func (ii IndexInfo) copy() IndexInfo {
	f := make([]Field, len(ii.Fields))
	copy(f, ii.Fields)
	ii.Fields = f
	return ii
}

func (ii *IndexInfo) validate() error {

	if ii.Name == "" {
		return &ValidationError{
			Err:   "index name cannot be empty",
			Field: "name",
		}
	}

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

func (f *Field) validate() error {

	if f.Name == "" {
		return &ValidationError{
			Err:   "field name cannot be empty",
			Field: "name",
		}
	}

	return nil

}

func (p *PartitionAlloc) validate() error {
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
	CreateIndex(index IndexInfo) (IndexInfo, error)
	UpdateIndex(index IndexInfo) (IndexInfo, error)
	DeleteIndex(id string) error

	AllFields(id string) ([]Field, error)
	Field(id string) (Field, error)
	FieldByName(name string) (Field, error)
	UpdateField(field Field) error
	DeleteField(id string) error
	CreateFields(id string, fields Field) (Field, error)

	Alloc(id string) (PartitionAlloc, error)
	UpdateAlloc(id string, parAlloc PartitionAlloc) error

	Shutdown()
}

type schema struct {
	catalog     cluster.Catalog
	indexCache  map[string]IndexInfo
	fieldCache  map[string]Field
	fieldByName map[string]Field
	pAllocCache map[string]PartitionAlloc
	indexByName map[string]IndexInfo
	log         zerolog.Logger
	catalogCh   chan []cluster.Entry
	mu          sync.Mutex
	done        chan struct{}
}

func (s *schema) AllIndex() []IndexInfo {

	s.mu.Lock()
	defer s.mu.Unlock()

	indexes := make([]IndexInfo, len(s.indexCache))

	i := 0

	for _, index := range s.indexCache {
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

func (s *schema) CreateIndex(indexInfo IndexInfo) (IndexInfo, error) {

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
		f.Id = uuid.New().String()
		f.IndexId = indexInfo.Id
		f.Created = now
		f.Updated = now
		indexInfo.Fields[i] = f
	}

	entries, err := createEntries(indexInfo, now)

	s.catalog.SetAll(entries)

	s.addToCache(indexInfo)

	return indexInfo.copy(), nil

}

func (s *schema) addToCache(indexInfo IndexInfo) {
	indexInfo.init()
	s.indexCache[indexInfo.Id] = indexInfo
	s.indexByName[indexInfo.Name] = indexInfo
	for _, f := range indexInfo.Fields {
		s.fieldCache[f.Id] = f
		s.fieldByName[f.Name] = f
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
			Err:   fmt.Sprintf("NumOfPartitions cannot be updated in index %s", index.Name),
			Field: "PartitionAlloc.NumOfPartitions",
		}
	}

	index.init()

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

	s.catalog.SetAll(entries)

	for _, f := range fieldsToRemove {
		delete(s.fieldCache, f.Id)
	}

	// remove the old index name
	delete(s.indexByName, indexInfo.Name)
	s.addToCache(index)

	return index.copy(), nil
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

	s.catalog.SetAll(entries)

	delete(s.indexCache, indexInfo.Id)
	delete(s.indexByName, indexInfo.Name)

	for _, f := range indexInfo.Fields {
		delete(s.fieldCache, f.Id)
	}

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

func (s *schema) FieldByName(name string) (Field, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, ok := s.fieldByName[name]

	if !ok {
		return Field{}, &NotFoundError{
			Err: fmt.Sprintf("field %v doesnt exist", name),
		}
	}

	return f, nil
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

	s.catalog.Set(fieldEntry)

	s.fieldCache[field.Id] = field
	s.indexCache[i.Id] = i
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

	s.catalog.Set(fieldEntry)

	delete(s.fieldCache, f.Id)
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

	s.catalog.Set(entry)

	s.fieldCache[field.Id] = field
	s.indexCache[i.Id] = i
	s.indexByName[i.Name] = i

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

	s.catalog.Set(entry)
	index.PartitionAlloc = pAlloc

	return nil

}

func NewSchema(catalog cluster.Catalog) (Schema, error) {

	s := &schema{
		catalog:     catalog,
		indexCache:  make(map[string]IndexInfo),
		fieldCache:  make(map[string]Field),
		pAllocCache: make(map[string]PartitionAlloc),
		indexByName: make(map[string]IndexInfo),
		log:         log.With().Str("component", "schema").Logger(),
		catalogCh:   make(chan []cluster.Entry),
		done:        make(chan struct{}),
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

	//TODO(gvelo): read from index and field buckets should be synchronized.

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
			s.syncDelta(delta)
		}
	}

}

func (s *schema) syncDelta(delta []cluster.Entry) {

	s.mu.Lock()
	defer s.mu.Unlock()

	// update all fields and Pallocs.
	for _, e := range delta {
		switch e.MapName {
		case fieldsMapName:
			s.updateField(e)
		case pAllocMapName:
			s.updatePAlloc(e)
		}
	}

	// update all index.
	for _, e := range delta {
		if e.MapName == indexMapName {
			s.updateIndex(e)
		}
	}

}

func (s *schema) updateField(e cluster.Entry) {

	var field Field

	err := proto.Unmarshal(e.Value, &field)

	if err != nil {
		s.log.Panic().Err(err).Msg("Cannot unmarshal field value")
	}

	if e.Deleted {
		delete(s.fieldCache, field.Id)
		index, ok := s.indexCache[field.IndexId]
		if ok {
			index.removeField(field.Id)
			s.indexCache[index.Id] = index
		}
		return
	}

	s.fieldCache[field.Id] = field
	index, ok := s.indexCache[field.IndexId]
	if ok {
		index.addField(field)
		s.indexCache[index.Id] = index
	}

}

func (s *schema) updatePAlloc(e cluster.Entry) {

	var pAlloc PartitionAlloc

	err := proto.Unmarshal(e.Value, &pAlloc)

	if err != nil {
		s.log.Panic().Err(err).Msg("Cannot unmarshal pAlloc value")
	}

	s.pAllocCache[e.Key] = pAlloc

	index, ok := s.indexCache[e.Key]
	if ok {
		index.PartitionAlloc = pAlloc
		s.indexCache[index.Id] = index
	}

}

func (s *schema) updateIndex(e cluster.Entry) {

	var index Index

	err := proto.Unmarshal(e.Value, &index)

	if err != nil {
		s.log.Panic().Err(err).Msg("cannot unmarshall index catalog entry")
	}

	// fields will be deleted on updateFields
	if e.Deleted {
		delete(s.indexCache, index.Id)
		delete(s.indexByName, index.Name)
		return
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

	s.indexCache[indexInfo.Id] = indexInfo
	s.indexByName[indexInfo.Name] = indexInfo

}
