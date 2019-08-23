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

import (
	"context"
	"fmt"
	"github.com/hashicorp/serf/serf"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"sync"
)

type Transport interface {
	SendDelta(ctx context.Context, delta []Entry)
	GetSnapShot(ctx context.Context, members []string, out chan SnapshotResult)
}

type SnapshotResult struct {
	member   string
	err      error
	Snapshot []Entry
}

type transport struct {
	log         zerolog.Logger
	connections sync.Map
	chUpdates   chan serf.Event
}

func NewTransport(c Cluster) (Transport, error) {

	t := &transport{
		log:       log.With().Str("component", "transport").Logger(),
		chUpdates: make(chan serf.Event, 128),
	}

	t.log.Info().Msg("creating catalog transport")

	c.AddEventChan(t.chUpdates)

	for _, m := range c.LiveMembers() {
		err := t.addConnection(m)
		if err != nil {
			return nil, err
		}
	}

	go t.updateConn()

	return t, nil

}

func (t *transport) addConnection(m serf.Member) error {

	t.log.Debug().Msgf("connection to member %v on %v", m.Name, m.Addr)

	_, ok := t.connections.Load(m.Name)

	if ok {
		return nil
	}

	//TODO(gvelo): add transport security.
	c, err := grpc.Dial(m.Addr.String(), grpc.WithInsecure())

	if err != nil {
		return err
	}

	t.connections.Store(m.Name, c)

	return nil
}

func (t *transport) addConnections(members []serf.Member) {
	for _, m := range members {
		err := t.addConnection(m)
		if err != nil {
			t.log.Error().Err(err).Msgf("error adding connection for member %v on addr %v", m.Name, m.Addr)
		}
	}
}

func (t *transport) removeConnections(members []serf.Member) {
	for _, m := range members {
		t.log.Debug().Msgf("removing connection for member %v on %v", m.Name, m.Addr)
		v, ok := t.connections.Load(m.Name)
		if !ok {
			continue
		}
		conn := v.(grpc.ClientConn)
		err := conn.Close()
		if err != nil {
			t.log.Error().Err(err).Msgf("error closing connection to %v", m.Name)
		}
		t.connections.Delete(m.Name)
	}
}

func (t *transport) updateConn() {

	// TODO: review this logic when member readiness status
	// be added to signal that the bootstrap has ended.

	for e := range t.chUpdates {

		mEvent, ok := e.(serf.MemberEvent)

		if !ok {
			continue
		}

		switch mEvent.EventType() {
		case serf.EventMemberJoin:
			t.addConnections(mEvent.Members)
		case serf.EventMemberLeave:
			t.removeConnections(mEvent.Members)
		case serf.EventMemberUpdate:
			t.addConnections(mEvent.Members)
		case serf.EventMemberReap:
			t.removeConnections(mEvent.Members)
		default:
			t.log.Error().Msgf("Unknown serf event %v", e.EventType())
		}

	}
}

func (t *transport) SendDelta(ctx context.Context, delta []Entry) {

	t.connections.Range(func(key, value interface{}) bool {

		m := key.(string)
		cc := value.(*grpc.ClientConn)

		go func(m string, cc *grpc.ClientConn) {

			cl := NewCatalogClient(cc)

			request := &AddRequest{
				Entries: delta,
			}

			t.log.Debug().Str("member", m).Msg("sending catalog delta to member")

			_, err := cl.Add(ctx, request)

			if err != nil {
				t.log.Error().Str("member", m).Err(err).Msg("error sending catalog delta")
			}

		}(m, cc)

		// keep iterating
		return true

	})

}

func (t *transport) GetSnapShot(ctx context.Context, members []string, out chan SnapshotResult) {

	for _, m := range members {

		v, ok := t.connections.Load(m)

		if !ok {
			r := SnapshotResult{
				member: m,
				err:    fmt.Errorf("member %v doesn't exist", m),
			}
			out <- r
			continue
		}

		cc := v.(*grpc.ClientConn)

		go func(m string, cc *grpc.ClientConn) {

			cl := NewCatalogClient(cc)
			r, err := cl.SnapShot(ctx, &SnapshotRequest{})

			if err != nil {
				out <- SnapshotResult{member: m, err: err}
				return
			}

			out <- SnapshotResult{
				member:   m,
				Snapshot: r.Entries,
			}

		}(m, cc)
		// TODO(gvelo): should we close the channel ?
	}
}
