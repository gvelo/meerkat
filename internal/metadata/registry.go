// Copyright 2021 The Meerkat Authors
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

package metadata

//go:generate protoc  -I . -I ../../build/proto/ -I ../../internal/storage/ --plugin ../../build/protoc-gen-gogofaster --gogofaster_out=plugins=grpc:. ./metadata.proto

// type NodeList []string
// type PartitionMap []NodeList
//
// type DatabaseIngestionConfig struct {
// 	AllowDynamicTableCreation bool
// 	TableAliases              map[string]string
// }
//
// type DatabaseQueryConfig struct {
// 	TableAliases map[string]string
// }
//
// type DatabaseDescriptor struct {
// 	Name            string
// 	NumOfPartitions int
// 	PartitionMap    PartitionMap
// 	IngestionConf   *DatabaseIngestionConfig
// 	QueryConf       *DatabaseQueryConfig
// 	Version         int
// }
//
// type ColumnMapping struct {
// 	Name       string
// 	ColumnType storage.ColumnType
// 	Nullable   bool
// 	Indexed    bool
// }
//
// type TableIngestionConfig struct {
// 	AllowDynamicColumnCreation bool
// 	ColumnMapping              map[string]*ColumnMapping
// 	ColumnAliases              map[string]string
// }
//
// type TableQueryConfig struct {
// 	ColumnAliases map[string]string
// }
//
// type Table struct {
// 	Name            string                `json:"name"`
// 	NumOfPartitions int                   `json:"numOfPartitions"`
// 	PartitionMap    PartitionMap          `json:"partitionMap"`
// 	QueryConfig     *TableQueryConfig     `json:"queryConfig"`
// 	IngestionConfig *TableIngestionConfig `json:"ingestionConfig"`
// 	Version         int                   `json:"version"`
// }
//
// type Metadata interface {
// 	GetDbs() []*Database
// 	GetDb(dbName string) *Database
// 	GetDbByAlias(dbName string) *Database
// 	GetTables(dbName string) []*Table
// 	GetTablesByAlias(dbAlias string) []*Table
// 	GetTable(dbName string, tableName string) *Table
// 	GetTableByAlias(dbName string, tableName string) *Table
// 	Mutate(mutation Mutation)
// 	Log() Log
// }
//
// type Operation int
//
// const (
// 	opCreate Operation = iota
// 	opUpdate
// 	opDelete
// )
//
// type Mutation interface {
// 	Id() string
// 	SrcNode() string
// 	Operation() Operation
// }
//
// type CreateDbMutation struct {
// 	Name            string
// 	NumOfPartitions int
// 	PartitionMap    PartitionMap
// 	IngestionConfig *DatabaseIngestionConfig
// 	QueryConfig     *DatabaseQueryConfig
// }
//
// type db struct {
// 	Name string
// }
//
// type metadata struct {
// 	dbAliases map[string]string
// }
