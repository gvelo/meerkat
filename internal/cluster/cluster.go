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
	"encoding/json"
	"github.com/google/uuid"
	"github.com/hashicorp/serf/serf"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net"
	"os"
	"path"
	"sync"
)

const (
	confFile = "cluster.json"
	confPerm = 0660

	tagBoostrapping = "bootstrapping"
	tagHostname     = "hostname"
)

// Cluster is used to track cluster node membership and to inform
// local state to other members node. It is a thin  wraper around Serf.
type Cluster interface {
	// SetTag is used to dynamically update the tags associated with
	// the local node. This will propagate the change to the rest of
	// the cluster. Blocks until a the message is broadcast out.
	SetTag(name string, value string) error

	// Members returns a point-in-time snapshot of the members of
	// this cluster. This includes failed, left and non-ready
	// (not yet bootstrapped) nodes.
	Members() []serf.Member

	// LiveMembers returns a point-in-time snapshot of the members of
	// this cluster. This includes only nodes ready to receive requests.
	LiveMembers() []serf.Member

	// Start trying to join the cluster. Return error if no node could
	// be contacted.
	Start() error

	// Shutdow first leave the cluster gracefully and then shuts down
	// the Serf instance, stopping all networkvactivity and background
	// maintenance associated with the instance.
	Shutdown()

	// Ready set the ready state on the local node. A call to this method signal
	// that the bootstrapping has finished and the node is ready to receive
	// requests from other nodes in the cluster.
	Ready() error

	// AddEventChan add a channel to receive serf events. Care must be taken
	// that this channel doesn't block either by processing the event quick
	// enough or providing a buffered channel with enought capacity,
	// otherwise is can block states update withing Serf itself.
	AddEventChan(ch chan serf.Event)

	// RemoveEventChan removes a previously added channel.
	RemoveEventChan(ch chan serf.Event)
}

// clusterConfig store local node info and last known nodes.
type clusterConfig struct {
	Name  string
	Nodes []net.IP
}

// NewCluster return a new cluster instance.
func NewCluster(port int, seeds []string, dbPath string) Cluster {

	cl := &cluster{
		port:      port,
		confPath:  path.Join(dbPath, confFile),
		seeds:     seeds,
		log:       log.With().Str("component", "cluster").Logger(),
		tags:      make(map[string]string),
		eventChan: make(map[chan serf.Event]chan serf.Event),
	}

	return cl

}

type cluster struct {
	port      int
	seeds     []string
	log       zerolog.Logger
	conf      clusterConfig
	confPath  string
	serf      *serf.Serf
	hostname  string
	tags      map[string]string
	mu        sync.Mutex
	eventChan map[chan serf.Event]chan serf.Event
	serfChan  chan serf.Event
}

func (c *cluster) AddEventChan(ch chan serf.Event) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.eventChan[ch] = ch
}

func (c *cluster) RemoveEventChan(ch chan serf.Event) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.eventChan, ch)
}

func (c *cluster) SetTag(name string, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.log.Info().Msgf("set node tag [%v]=%v", name, value)

	err := c.serf.SetTags(c.tags)

	return err
}

func (c *cluster) Members() []serf.Member {
	return c.serf.Members()
}

func (c *cluster) LiveMembers() []serf.Member {

	members := c.serf.Members()

	live := members[:0]

	for _, m := range members {
		if m.Status == serf.StatusAlive && m.Tags[tagBoostrapping] == "F" {
			live = append(live, m)
		}
	}

	return live
}

func (c *cluster) join() error {

	// TODO(gvelo): add retry if we are unable to join any nodes.

	// first, try to join using the provided seeds.
	if len(c.seeds) > 0 {

		c.log.Info().Msgf("trying to join seeds %v", c.seeds)

		n, err := c.serf.Join(c.seeds, true)

		if err == nil {
			c.log.Info().Msgf("%v nodes successfully joined", n)
			return nil
		}

	}

	// TODO(gvelo): if we are unable to contact the seeds try
	// to contact the last known nodes in batch on N nodes.

	if len(c.conf.Nodes) > 0 {
		c.log.Info().Msgf("trying to join nodes %v", c.conf.Nodes)
	}

	c.log.Info().Msg("no nodes available")

	return nil

}

func (c *cluster) Shutdown() {

	c.log.Info().Msg("shuttingdown cluster")

	err := c.serf.Leave()

	if err != nil {
		c.log.Error().Err(err).Msgf("error leaving cluster")
	} else {
		c.log.Info().Msg("cluster leave ok")
	}

	err = c.serf.Shutdown()

	if err != nil {
		c.log.Error().Err(err).Msg("error shuttingdown cluster")
	} else {
		c.log.Info().Msg("cluster shutdown ok")
	}

}

func (c *cluster) Ready() error {
	err := c.SetTag(tagBoostrapping, "F")
	return err
}

func (c *cluster) Start() error {

	c.log.Info().Msg("starting cluster")

	err := c.initConfig()

	if err != nil {
		return err
	}

	err = c.initSerf()

	if err != nil {
		return err
	}

	err = c.join()

	return err

}

func (c *cluster) initConfig() error {

	err := c.loadConfig()

	if os.IsNotExist(err) {
		c.log.Info().Msg("cannot find node config, creating new node.")
		err = c.newNode()
		if err != nil {
			return err
		}
		return nil
	}

	return err

}

func (c *cluster) initSerf() error {

	hostName, err := os.Hostname()

	if err != nil {
		c.log.Error().Err(err).Msg("unable to determine hostname")
		return err
	}

	c.log.Info().Msgf("hostname %v", hostName)

	c.serfChan = make(chan serf.Event, 1024)
	go c.dispatchEvents()

	serfConf := serf.DefaultConfig()
	serfConf.Init()
	serfConf.NodeName = c.conf.Name
	serfConf.Tags[tagHostname] = hostName
	serfConf.Tags[tagBoostrapping] = "T"
	serfConf.EventCh = c.serfChan

	if c.port != -1 {
		serfConf.MemberlistConfig.BindPort = c.port
	}

	// TODO(gvelo): redirect serf log to zerolog
	// TODO(gvelo): save last known members to cluster config file

	c.serf, err = serf.Create(serfConf)

	return err

}

func (c *cluster) newNode() error {

	c.log.Info().Msgf("creating new node")

	c.conf.Name = uuid.New().String()
	c.conf.Nodes = make([]net.IP, 0)

	c.log.Info().Msgf("new node created nodename=%v", c.conf.Name)

	return c.saveConfig()

}

func (c *cluster) loadConfig() error {

	c.log.Info().Msgf("loading cluster config from %v", c.confPath)

	b, err := ioutil.ReadFile(c.confPath)

	if err != nil {
		return err
	}

	conf := clusterConfig{}

	err = json.Unmarshal(b, &conf)

	if err != nil {
		return err
	}

	c.conf = conf

	return nil

}

func (c *cluster) saveConfig() error {

	c.log.Info().Msgf("saving cluster configuration")

	b, err := json.Marshal(c.conf)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.confPath, b, confPerm)

	if err != nil {
		return err
	}

	return nil

}

func (c *cluster) dispatchEvents() {

	c.log.Info().Msg("start dispatching serf events")

	for e := range c.serfChan {

		c.log.Debug().Msgf("serf event received %v", e)

		// filter events from the local member
		if me, ok := e.(serf.MemberEvent); ok {
			f := me.Members[:0]
			for _, m := range me.Members {
				if m.Name != c.conf.Name {
					f = append(f, m)
				}
			}
			// drop event if it only contain local member
			if len(f) == 0 {
				continue
			}
		}

		c.log.Debug().Msgf("dispatching serf event %v", e)

		for c := range c.eventChan {
			c <- e
		}

	}
}
