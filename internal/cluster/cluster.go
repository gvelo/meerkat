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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/serf/serf"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"os"
	"path"
	"sync"
)

const (
	confFile    = "cluster.json"
	confPerm    = 0660
	tagHostname = "hostname"
	tagStatus   = "status"
)

const (
	created = "created"
	Joining = "joining"
	Ready   = "ready"
	Leaving = "leaving"
	Leave   = "leave"
	Fail    = "fail"
)

type Node interface {
	Id() string
	Addr() net.IP
	Tag(tagName string) string
	Status() string
	ClientConn() *grpc.ClientConn
}

type node struct {
	id         string
	addr       net.IP
	tags       map[string]string
	status     string
	clientConn *grpc.ClientConn
}

func (n *node) Id() string {
	return n.id
}

func (n *node) Addr() net.IP {
	return n.addr
}

func (n *node) Tag(tagName string) string {
	return n.tags[tagName]
}

func (n *node) Status() string {
	return n.status
}

func (n *node) ClientConn() *grpc.ClientConn {
	return n.clientConn
}

func (n *node) inStatus(statusList []string) bool {
	for _, status := range statusList {
		if status == n.status {
			return true
		}
	}
	return false
}

func (n *node) MarshalJSON() ([]byte, error) {

	j := map[string]interface{}{
		"id":     n.id,
		"addr":   n.addr,
		"status": n.status,
		"tags":   n.tags,
	}

	return json.Marshal(j)

}

// Cluster is used to track cluster node membership and to inform
// local state to other members node. It is a thin  wrapper around Serf.
// TODO(gvelo): segregate the Cluster interface on NodeRegistry and
// ClusterStateManager
type Cluster interface {

	// Join the cluster and transition to Joining state
	Join() error

	// Ready Set the node status to Ready
	Ready() error

	// Leaving set the node status to leaving
	Leaving() error

	// Leave the cluster gracefully and then shuts down the Serf instance,
	// stopping all network activity and background maintenance associated with
	// the instance.
	Leave() error

	// SetTag is used to dynamically update the tags associated with
	// the local node. This will propagate the change to the rest of
	// the cluster. Blocks until a the message is broadcast out.
	SetTag(name string, value string) error

	// Nodes return a list of nodes filtered by status.
	// if excludeLocalNode is true the returned list will not contain the local
	// node info.
	Nodes(statusFilter []string, excludeLocalNode bool) []Node

	Node(id string) Node

	// Return the local node id
	NodeId() string
}

// clusterConfig store local node info and last known nodes.
type clusterConfig struct {
	Id    string
	Nodes []net.IP
}

// NewCluster return a new cluster instance.
func NewCluster(port int, seeds []string, dbPath string) (Cluster, error) {

	c := &cluster{
		port:     port,
		confPath: path.Join(dbPath, confFile),
		seeds:    seeds,
		log:      log.With().Str("src", "cluster").Logger(),
		tags:     make(map[string]string),
		status:   created,
		serfChan: make(chan serf.Event, 1024),
		nodes:    make(map[string]*node),
	}

	c.log.Info().Msg("creating cluster")

	hostName, err := os.Hostname()

	if err != nil {
		c.log.Error().Err(err).Msg("unable to determine hostname")
		return nil, err
	}

	c.log.Info().Msgf("hostname %v", hostName)

	c.tags[tagStatus] = Joining
	c.tags[tagHostname] = hostName

	err = c.initConfig()

	if err != nil {
		return nil, err
	}

	err = c.initSerf()

	if err != nil {
		return nil, err
	}

	return c, nil

}

type cluster struct {
	port     int
	seeds    []string
	log      zerolog.Logger
	conf     clusterConfig
	confPath string
	serf     *serf.Serf
	hostname string
	tags     map[string]string
	statusMu sync.Mutex
	serfChan chan serf.Event
	status   string
	nodes    map[string]*node
	nodesMu  sync.Mutex
}

func (c *cluster) SetTag(name string, value string) error {

	c.statusMu.Lock()
	defer c.statusMu.Unlock()

	c.log.Info().Msgf("setting node tag [%v]=%v", name, value)

	if c.tags[name] == value {
		return nil
	}

	c.tags[name] = value

	err := c.serf.SetTags(c.tags)

	return err
}

func (c *cluster) Nodes(statusFilter []string, excludeLocalNode bool) []Node {

	c.nodesMu.Lock()
	defer c.nodesMu.Unlock()

	var result []Node

	for _, node := range c.nodes {

		if excludeLocalNode && node.id == c.conf.Id {
			continue
		}

		if len(statusFilter) == 0 || node.inStatus(statusFilter) {
			result = append(result, node)
		}

	}

	return result

}

func (c *cluster) Node(id string) Node {

	c.nodesMu.Lock()
	defer c.nodesMu.Unlock()

	return c.nodes[id]
}

func (c *cluster) Join() error {

	c.statusMu.Lock()
	defer c.statusMu.Unlock()

	if c.status != created {
		return fmt.Errorf("invalid state %v", c.status)
	}

	var nodes []string
	nodes = append(nodes, c.seeds...)
	for _, ip := range c.conf.Nodes {
		nodes = append(nodes, ip.String())
	}

	if len(nodes) > 0 {

		c.log.Info().Msgf("trying to join nodes %v", nodes)

		// TODO(gvelo) contact nodes in batch.
		n, err := c.serf.Join(nodes, true)

		if err != nil {
			c.log.Error().Msg("cannot contact any node")
		}

		c.log.Info().Msgf("%v nodes successfully contacted", n)

		c.status = Joining

		return nil

	}

	// there is no seeds or known hosts, we are the first node of a cluster

	c.log.Info().Msg("no nodes to contact")

	c.status = Joining
	return nil

}

func (c *cluster) Ready() error {

	c.statusMu.Lock()
	defer c.statusMu.Unlock()

	if c.status != Joining {
		return fmt.Errorf("cannot transition to ready state from %v", c.status)
	}

	c.tags[tagStatus] = Ready

	err := c.serf.SetTags(c.tags)

	if err != nil {
		return fmt.Errorf("cannot update tags: %v", err)
	}

	c.status = Ready

	return nil

}

func (c *cluster) Leaving() error {

	c.statusMu.Lock()
	defer c.statusMu.Unlock()

	if c.status != Ready {
		return fmt.Errorf("cannot transition to leaving state from %v", c.status)
	}

	c.tags[tagStatus] = Leaving

	err := c.serf.SetTags(c.tags)

	if err != nil {
		return fmt.Errorf("cannot update tags: %v", err)
	}

	c.status = Leaving

	return nil

}

func (c *cluster) Leave() error {

	c.statusMu.Lock()
	defer c.statusMu.Unlock()

	if c.status != Leaving {
		return fmt.Errorf("cannot transition to leave state from %v", c.status)
	}

	c.log.Info().Msg("leaving")

	err := c.serf.Leave()

	if err != nil {
		c.log.Error().Err(err).Msgf("error leaving cluster")
	} else {
		c.log.Info().Msg("cluster leave ok")
	}

	err = c.serf.Shutdown()

	if err != nil {
		c.log.Error().Err(err).Msg("error shuttingdown cluster")
		return err
	}

	c.log.Info().Msg("cluster shutdown ok")

	return nil

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

func (c *cluster) NodeId() string {
	return c.conf.Id
}

func (c *cluster) initSerf() error {

	serfConf := serf.DefaultConfig()
	serfConf.Init()
	serfConf.NodeName = c.conf.Id
	serfConf.EventCh = c.serfChan

	for tagName, tagValue := range c.tags {
		serfConf.Tags[tagName] = tagValue
	}

	if c.port != -1 {
		serfConf.MemberlistConfig.BindPort = c.port
	}

	go c.handleSerfEvents()

	var err error
	c.serf, err = serf.Create(serfConf)

	return err

}

func (c *cluster) newNode() error {

	c.log.Info().Msgf("creating new node")

	id := uuid.New()

	c.conf.Id = base64.RawURLEncoding.EncodeToString(id[:])

	c.log.Info().Msgf("new node created %v", c.conf.Id)

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

func (c *cluster) handleSerfEvents() {

	c.log.Info().Msg("starting serf event dispatcher")

	for serfEvent := range c.serfChan {

		c.log.Debug().Msgf("serf event received %v", serfEvent)

		if memberEvent, ok := serfEvent.(serf.MemberEvent); ok {

			switch memberEvent.EventType() {

			case serf.EventMemberJoin:
				c.addNodes(memberEvent.Members)
			case serf.EventMemberLeave:
				c.updateNodes(memberEvent.Members)
			case serf.EventMemberFailed:
				c.updateNodes(memberEvent.Members)
			case serf.EventMemberUpdate:
				c.updateNodes(memberEvent.Members)
			case serf.EventMemberReap:
				c.removeNodes(memberEvent.Members)
			default:
				c.log.Error().Msgf("unknown serf event %v", memberEvent.EventType())
			}

		}

	}
}

func (c *cluster) addNodes(members []serf.Member) {

	for _, member := range members {

		c.log.Debug().Msgf("adding node %v", member.Name)

		var conn *grpc.ClientConn
		var err error

		// create a grpc connection to every node except the local one
		if c.conf.Id != member.Name {

			conn, err = c.createGrpcConn(member.Addr)

			if err != nil {
				c.log.Error().Err(err).Msgf("cannot create grpc connection to node %v")
				continue
			}

		}

		node := createNode(member)
		node.clientConn = conn

		c.nodesMu.Lock()
		c.nodes[node.id] = node
		c.nodesMu.Unlock()

	}

}

func (c *cluster) updateNodes(members []serf.Member) {

	c.nodesMu.Lock()
	defer c.nodesMu.Unlock()

	for _, member := range members {

		oldNode, ok := c.nodes[member.Name]

		if !ok {
			c.log.Error().Msgf("update node %v not found")
			continue
		}

		c.log.Debug().Msgf("updating node %v state %v", member.Name, member.Status)

		node := createNode(member)
		node.clientConn = oldNode.clientConn
		c.nodes[node.id] = node

	}

}

func (c *cluster) removeNodes(members []serf.Member) {

	c.nodesMu.Lock()
	defer c.nodesMu.Unlock()

	for _, member := range members {

		node, ok := c.nodes[member.Name]

		if !ok {
			c.log.Error().Msgf("node %v not found")
			continue
		}

		c.log.Debug().Msgf("removing node %v state %v", member.Name, member.Status)

		err := node.clientConn.Close()

		if err != nil {
			c.log.Error().Err(err).Msgf("error closing client connection")
		}

		delete(c.nodes, node.id)

	}

}

func (c *cluster) createGrpcConn(addr net.IP) (*grpc.ClientConn, error) {

	// TODO(gvelo): add transport security.
	// TODO(gvelo): externalize grpc port
	target := addr.String() + ":9191"
	// grpc.Dial is a non-blocking call so it is safe to use it inside the
	// serf dispatching gorutine.
	return grpc.Dial(target, grpc.WithInsecure())

}

func createNode(member serf.Member) *node {

	node := &node{
		id:   member.Name,
		addr: member.Addr,
		tags: member.Tags,
	}

	switch member.Status {
	case serf.StatusAlive:
		node.status = member.Tags[tagStatus]
	case serf.StatusLeaving:
		node.status = Leaving
	case serf.StatusLeft:
		node.status = Leave
	case serf.StatusFailed:
		node.status = Fail
	default:
		node.status = "unknown status"
	}

	return node

}
