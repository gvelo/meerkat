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
	"github.com/hashicorp/serf/serf"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"sync"
)

type ConnRegistry interface {
	Range(f func(member string, conn *grpc.ClientConn) bool)
	RangeWithFilter(members []string, fn func(member string, conn *grpc.ClientConn, ok bool) bool)
}

type connRegistry struct {
	log         zerolog.Logger
	connections sync.Map
	chUpdates   chan serf.Event
}

func NewConnRegistry(cluster Cluster) (ConnRegistry, error) {

	cc := &connRegistry{
		log:       log.With().Str("component", "connRegistry").Logger(),
		chUpdates: make(chan serf.Event, 1024),
	}

	cc.log.Info().Msg("creating connection registry")

	cluster.AddEventChan(cc.chUpdates)

	for _, m := range cluster.LiveMembers() {
		err := cc.addConnection(m)
		if err != nil {
			return nil, err
		}
	}

	go cc.updateConn()

	return cc, nil

}

func (cr *connRegistry) addConnection(m serf.Member) error {

	cr.log.Debug().Msgf("adding connection to member %v on %v", m.Name, m.Addr)

	_, ok := cr.connections.Load(m.Name)

	if ok {
		cr.log.Debug().Msgf("connection to member %v on %v already exist", m.Name, m.Addr)
		return nil
	}

	// TODO(gvelo): add transport security.
	// TODO(gvelo): externalize grpc port
	addr := m.Addr.String() + ":9191"
	c, err := grpc.Dial(addr, grpc.WithInsecure())

	if err != nil {
		return err
	}

	cr.connections.Store(m.Name, c)

	return nil
}

func (cr *connRegistry) addConnections(members []serf.Member) {
	for _, m := range members {
		err := cr.addConnection(m)
		if err != nil {
			cr.log.Error().Err(err).Msgf("error adding connection for member %v on addr %v", m.Name, m.Addr)
		}
	}
}

func (cr *connRegistry) removeConnections(members []serf.Member) {
	for _, m := range members {
		cr.log.Debug().Msgf("removing connection for member %v on %v", m.Name, m.Addr)
		v, ok := cr.connections.Load(m.Name)
		if !ok {
			continue
		}
		conn := v.(grpc.ClientConn)
		err := conn.Close()
		if err != nil {
			cr.log.Error().Err(err).Msgf("error closing connection to %v", m.Name)
		}
		cr.connections.Delete(m.Name)
	}
}

func (cr *connRegistry) updateConn() {

	// TODO: review this logic when member readiness status
	//       be added to signal that the bootstrap has ended.

	for e := range cr.chUpdates {

		mEvent, ok := e.(serf.MemberEvent)

		if !ok {
			continue
		}

		switch mEvent.EventType() {
		case serf.EventMemberJoin:
			cr.addConnections(mEvent.Members)
		case serf.EventMemberLeave:
			cr.removeConnections(mEvent.Members)
		case serf.EventMemberUpdate:
			// TODO(gvelo): should we handle EventMemberUpdate ?
			//cr.addConnections(mEvent.Members)
		case serf.EventMemberReap:
			cr.removeConnections(mEvent.Members)
		default:
			cr.log.Error().Msgf("Unknown serf event %v", e.EventType())
		}

	}
}

func (cr *connRegistry) Range(f func(member string, conn *grpc.ClientConn) bool) {
	cr.connections.Range(func(m, conn interface{}) bool {
		return f(m.(string), conn.(*grpc.ClientConn))
	})
}

func (cr *connRegistry) RangeWithFilter(members []string, fn func(member string, conn *grpc.ClientConn, ok bool) bool) {

	for _, m := range members {

		c, ok := cr.connections.Load(m)

		var clientConn *grpc.ClientConn

		if ok {
			clientConn = c.(*grpc.ClientConn)
		}

		if !fn(m, clientConn, ok) {
			break
		}

	}

}
