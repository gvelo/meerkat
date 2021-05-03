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

 package cluster
//
// import (
// 	"context"
// 	"github.com/rs/zerolog"
// 	"github.com/rs/zerolog/log"
// 	"google.golang.org/grpc"
// )
//
// type CatalogRPC interface {
// 	SendDelta(ctx context.Context, delta []Entry)
// 	GetSnapShot(ctx context.Context, members []string) []SnapshotResult
// }
//
// type SnapshotResult struct {
// 	member   string
// 	err      error
// 	Snapshot []Entry
// }
//
// type catalogRPC struct {
// 	log zerolog.Logger
// 	cl  Manager
// }
//
// func NewCatalogRPC(cl Manager) (CatalogRPC, error) {
//
// 	t := &catalogRPC{
// 		log: log.With().Str("component", "catalogRPC").Logger(),
// 		cl:  cl,
// 	}
//
// 	t.log.Info().Msg("creating catalogRPC")
//
// 	return t, nil
//
// }
//
// func (t *catalogRPC) SendDelta(ctx context.Context, delta []Entry) {
//
// 	nodes := t.cl.Nodes([]string{Ready, Joining}, true)
//
// 	for _, node := range nodes {
//
// 		go func(nodeId string, cc *grpc.ClientConn) {
//
// 			catalogClient := NewCatalogClient(cc)
//
// 			request := &AddRequest{
// 				Entries: delta,
// 			}
//
// 			t.log.Debug().Str("member", nodeId).Msg("sending catalog delta to member")
//
// 			_, err := catalogClient.Add(ctx, request)
//
// 			if err != nil {
// 				t.log.Error().Str("member", nodeId).Err(err).Msg("error sending catalog delta")
// 			}
//
// 		}(node.Id(), node.ClientConn())
//
// 	}
//
// }
//
// func (t *catalogRPC) GetSnapShot(ctx context.Context, nodeIds []string) []SnapshotResult {
//
// 	resultCh := make(chan SnapshotResult)
// 	memberCount := len(nodeIds)
// 	result := make([]SnapshotResult, memberCount)[:0]
//
// 	for _, nodeId := range nodeIds {
//
// 		node := t.cl.Node(nodeId)
//
// 		if node.Status() == Joining || node.Status() == Ready {
//
// 			go func(nodeId string, cc *grpc.ClientConn) {
//
// 				cl := NewCatalogClient(cc)
//
// 				r, err := cl.SnapShot(ctx, &SnapshotRequest{})
//
// 				if err != nil {
// 					resultCh <- SnapshotResult{member: nodeId, err: err}
// 					return
// 				}
//
// 				resultCh <- SnapshotResult{
// 					member:   nodeId,
// 					Snapshot: r.Entries,
// 				}
//
// 			}(node.Id(), node.ClientConn())
//
// 		}
//
// 	}
//
// 	// TODO(gvelo): unsafe, use a wg
// 	for i := 0; i < memberCount; i++ {
// 		r := <-resultCh
// 		result = append(result, r)
// 	}
//
// 	return result
//
// }
