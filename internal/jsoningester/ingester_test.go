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

package jsoningester

import (
	"bytes"
	"context"
	"github.com/hashicorp/serf/serf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"meerkat/internal/ingestion"
	"testing"
)

type clusterMock struct {
	mock.Mock
}

func (c *clusterMock) SetTag(name string, value string) error {
	args := c.Called(name, value)
	return args.Error(0)
}

func (c *clusterMock) Members() []serf.Member {
	args := c.Called()
	return args.Get(0).([]serf.Member)
}

func (c *clusterMock) LiveMembers() []serf.Member {
	args := c.Called()
	return args.Get(0).([]serf.Member)
}

func (c *clusterMock) Join() {
	c.Called()
}

func (c *clusterMock) Shutdown() {
	c.Called()
}

func (c *clusterMock) AddEventChan(ch chan serf.Event) {
	c.Called(ch)
}

func (c *clusterMock) RemoveEventChan(ch chan serf.Event) {
	c.Called(ch)
}

func (c *clusterMock) NodeName() string {
	args := c.Called()
	return args.String(0)
}

type ingestRpcMock struct {
	mock.Mock
}

func (i *ingestRpcMock) SendRequest(ctx context.Context, member string, request *ingestion.IngestionRequest) error {
	args := i.Called(ctx, member, request)
	return args.Error(0)
}

type buffRegMock struct {
	mock.Mock
}

func (b *buffRegMock) Add(table *ingestion.Table) {
	b.Called(table)
}

func TestIngestion(t *testing.T) {

	s := `{"_ts":"2020-05-11T18:46:06.577Z","columnA":23,"columnB":{ "columnF": 23}}
          {"_ts":"2020-05-11T18:46:06.672Z","columnA":24}
          {"_ts":"2020-05-11T18:46:07.443Z","columnA":25}`

	buffRegMock := &buffRegMock{}
	clusterMock := &clusterMock{}
	ingestRpcMock := &ingestRpcMock{}

	clusterMock.On("LiveMembers").Return([]serf.Member{{Name: "testmember1"}, {Name: "testmember2"}, {Name: "testmember3"}})
	ingestRpcMock.On("SendRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	buffRegMock.On("Add", mock.Anything)

	r := bytes.NewReader([]byte(s))

	ingester := NewIngester(ingestRpcMock, clusterMock, buffRegMock)
	err := ingester.Ingest(r, "testTable")

	assert.Equal(t, 0, len(err))

	clusterMock.AssertExpectations(t)
	ingestRpcMock.AssertExpectations(t)
	buffRegMock.AssertExpectations(t)

}

func TestIngestionError(t *testing.T) {

	s := `{"_ts":"2020-05-11T18:46:06.577Z","columnA":23}
          {"_ts":"2020-05-11T18:46:06.577Z,"columnA":23}
          {"_ts":"2020-05-11T18:46:06.672Z","columnA":24}
          {"_ts":"2020-05-11T18:46:07.443Z","columnA":25}`

	buffRegMock := &buffRegMock{}
	clusterMock := &clusterMock{}
	ingestRpcMock := &ingestRpcMock{}

	clusterMock.On("LiveMembers").Return([]serf.Member{{Name: "testmember1"}, {Name: "testmember2"}, {Name: "testmember3"}})
	ingestRpcMock.On("SendRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	buffRegMock.On("Add", mock.Anything)

	r := bytes.NewReader([]byte(s))

	ingester := NewIngester(ingestRpcMock, clusterMock, buffRegMock)
	err := ingester.Ingest(r, "testTable")

	assert.Equal(t, 1, len(err))

	clusterMock.AssertExpectations(t)
	ingestRpcMock.AssertExpectations(t)
	buffRegMock.AssertExpectations(t)

}
