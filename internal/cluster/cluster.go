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
)

type Cluster interface {
	SetTag(name string, value string) error
	Members() []serf.Member
	LiveMembers() []serf.Member
	Start() error
	Shutdown()
	Ready() error
	AddEventChan(ch chan serf.Event)
	RemoveEventChan(ch chan serf.Event)
}

type clusterConfig struct {
	Name  string
	Nodes []net.IP
}

func NewCluster(seeds []string, dbPath string) Cluster {

	cl := &cluster{
		confPath:  path.Join(dbPath, confFile),
		seeds:     seeds,
		log:       log.With().Str("component", "cluster").Logger(),
		tags:      make(map[string]string),
		eventChan: make(map[chan serf.Event]chan serf.Event),
	}

	return cl

}

type cluster struct {
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
		if m.Status == serf.StatusAlive && m.Tags["ready"] == "true" {
			live = append(live, m)
		}
	}

	return live
}

func (c *cluster) join() error {

	c.log.Info().Msgf("trying to join the Serf cluster")

	// first, try to join using the provided seeds.

	n, err := c.serf.Join(c.seeds, true)

	if err == nil {
		c.log.Info().Msgf("%v nodes successfully joined", n)
		return nil
	}

	// TODO(gvelo): if we are unable to contact the seeds try
	// to contact the last known nodes in batch on N nodes.

	return err

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
	err := c.SetTag("ready", "true")
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
	serfConf.Tags["hostname"] = hostName
	serfConf.Tags["bootstrapping"] = "true"
	serfConf.EventCh = c.serfChan

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
