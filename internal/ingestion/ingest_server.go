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
	"context"
	"fmt"
	"meerkat/internal/indexbuffer"
	"meerkat/internal/ingestion/ingestionpb"
)

func NewIngestionServer(bufReg indexbuffer.BufferRegistry) ingestionpb.IngesterServer {
	return &ingestServer{
		bufReg: bufReg,
	}
}

type ingestServer struct {
	bufReg indexbuffer.BufferRegistry
}

func (i ingestServer) Ingest(ctx context.Context, request *ingestionpb.IngestionRequest) (*ingestionpb.IngestResponse, error) {

	fmt.Println("------------ ingest server")

	i.bufReg.Add(request.Tables)
	return &ingestionpb.IngestResponse{},nil
}
