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

type broadcastResult struct {
	NodeId string
	err    error
	hash   int64
}

type logRpc interface {
	broadcast(entry *LogEntry) ([]*broadcastResult, error)
	fetchSnapshot(nodeId) ([]*LogEntry, error)
}

type logVersionUpdater struct {
}

type logAntientropy struct {
}

type logServer struct {

}

type Observer interface {
	OnAdd(entry []*LogEntry, hash int64)
}

type Log interface {
	Register(observer Observer)
	Deregister(observer Observer)
	Hash() int64
	AddOperation(op isLogEntry_Op)
	AddAll([]*LogEntry)
	Add([]*LogEntry)
	Log() []*LogEntry
}
