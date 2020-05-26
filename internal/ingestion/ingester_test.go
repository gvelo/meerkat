// Copyright 2020 The Meerkat Authors
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

package ingestion

import (
	"bytes"
	"context"
	"fmt"
	"github.com/hashicorp/serf/serf"
	"meerkat/internal/ingestion/ingestionpb"
	"testing"
)

var s string

func init() {
	//s := `{"_ts":"2020-05-11T18:46:06.577Z","columnA":23}
    //      {"_ts":"2020-05-11T18:46:06.577Z","columnA":23}
    //      {"_ts":"2020-05-11T18:46:06.577Z","columnA":23}

	s = `{"field1":0}
{"field1":1}
{"field1":2}

`
	//for i := 0; i < 20; i++ {
	//	s = s + s
	//}

	fmt.Println(len(s))
}

type clusterMock struct {
}

func (c clusterMock) SetTag(name string, value string) error {
	panic("implement me")
}

func (c clusterMock) Members() []serf.Member {
	panic("implement me")
}

func (c clusterMock) LiveMembers() []serf.Member {
	return []serf.Member{{Name: "testmember1"}, {Name: "testmember2"}, {Name: "testmember3"}}
}

func (c clusterMock) Join() {
	panic("implement me")
}

func (c clusterMock) Shutdown() {
	panic("implement me")
}

func (c clusterMock) AddEventChan(ch chan serf.Event) {
	panic("implement me")
}

func (c clusterMock) RemoveEventChan(ch chan serf.Event) {
	panic("implement me")
}

func (c clusterMock) NodeName() string {
	panic("implement me")
}

type ingestRpcMock struct {
}

func (i ingestRpcMock) SendRequest(ctx context.Context, member string, request *ingestionpb.IngestionRequest) error {
	fmt.Println(request.Table.Partitions)
	return nil
}

type buffReg struct {
}

func (b buffReg) Add(table *ingestionpb.Table) {

}

func TestTest(t *testing.T) {

	fmt.Println(s)
	r := bytes.NewReader([]byte(s))

	ing := NewIngester(&ingestRpcMock{}, &clusterMock{}, &buffReg{})

	err := ing.Ingest(r, "testTable")


}
