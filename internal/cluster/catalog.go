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

//go:generate protoc  -I . -I ../../build/proto/   --plugin ../../build/protoc-gen-gogofaster --gogofaster_out=plugins=grpc:. ./catalog.proto

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/twmb/murmur3"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc"
	"path"
	"sync"
	"time"
)

type Catalog interface {
	Set(entry Entry) bool
	SetAll(entries []Entry) []Entry
	Get(mapName string, key string) (Entry, bool)
	GetAll(mapName string) []Entry
	MergeSnapshot(catalog []Entry)
	SnapShot() []Entry
	Hash() string
	AddEventHandler(id string, h chan []Entry)
	RemoveEventHandler(id string)
}

type catalog struct {
	db            *bolt.DB
	mu            sync.Mutex
	log           zerolog.Logger
	maps          map[string]map[string]Entry
	replicaChan   chan []Entry
	hash          string
	eventHandlers map[string]chan []Entry
	replicator    *catalogReplicator
}

func (c *catalog) SnapShot() []Entry {

	var entries []Entry

	_ = c.db.View(func(tx *bolt.Tx) error {
		entries = c.entries(tx)
		return nil
	})

	return entries

}

func (c *catalog) entries(tx *bolt.Tx) []Entry {

	var entries []Entry

	_ = tx.ForEach(func(name []byte, bucket *bolt.Bucket) error {
		_ = bucket.ForEach(func(k, v []byte) error {
			var e Entry
			err := proto.Unmarshal(v, &e)
			if err != nil {
				c.log.Panic().Err(err).Msg("error unmarshalling value from db")
			}
			entries = append(entries, e)
			return nil
		})
		return nil
	})

	return entries
}

func (c *catalog) Set(entry Entry) bool {

	entries := []Entry{entry}
	delta := c.SetAll(entries)

	return len(delta) != 0

}

func (c *catalog) merge(tx *bolt.Tx, entry Entry) bool {

	bucket, err := tx.CreateBucketIfNotExists([]byte(entry.MapName))

	if err != nil {
		c.log.Panic().Err(err)
	}

	b := bucket.Get([]byte(entry.Key))

	if b != nil {

		var e Entry

		err := proto.Unmarshal(b, &e)

		if err != nil {
			c.log.Panic().Err(err).Msg("error unmarshalling catalog entry")
			return false
		}

		if e.Deleted {
			return false
		}

		if e.Time.After(entry.Time) {
			return false
		}

	}

	b, err = proto.Marshal(&entry)

	if err != nil {
		c.log.Panic().Err(err).Msg("error serializing catalog entry")
		return false
	}

	err = bucket.Put([]byte(entry.Key), b)

	if err != nil {
		c.log.Panic().Err(err).Msg("error storing catalog entry")
		return false
	}

	return true

}

func (c *catalog) mergeAll(tx *bolt.Tx, entries []Entry) []Entry {

	var delta []Entry

	for _, e := range entries {
		if c.merge(tx, e) {
			delta = append(delta, e)
		}
	}

	if len(delta) != 0 {
		c.hashCatalog(tx)
		c.emit(delta)
	}

	return delta

}

func (c *catalog) SetAll(entries []Entry) []Entry {

	var delta []Entry

	_ = c.db.Update(func(tx *bolt.Tx) error {

		delta = c.mergeAll(tx, entries)

		if len(delta) > 0 {
			c.replicate(delta)
		}
		return nil
	})

	return delta

}

func (c *catalog) Get(mapName string, key string) (Entry, bool) {

	var entry Entry
	var found bool

	_ = c.db.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(mapName))

		if bucket == nil {
			return nil
		}

		b := bucket.Get([]byte(key))

		if b != nil {

			err := proto.Unmarshal(b, &entry)

			if err != nil {
				c.log.Panic().Err(err).Msg("error unmarshalling entry")
			}

			if !entry.Deleted {
				found = true
			}

			return nil
		}

		return nil

	})

	return entry, found

}

func (c *catalog) GetAll(mapName string) []Entry {

	var entries []Entry

	_ = c.db.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(mapName))

		if bucket == nil {
			return nil
		}

		_ = bucket.ForEach(func(k, v []byte) error {

			var e Entry

			err := proto.Unmarshal(v, &e)

			if err != nil {
				c.log.Panic().Err(err).Msg("error unmarshalling entry")
			}

			if !e.Deleted {
				entries = append(entries, e)
			}

			return nil

		})

		return nil

	})

	return entries

}

func (c *catalog) MergeSnapshot(entries []Entry) {

	_ = c.db.Update(func(tx *bolt.Tx) error {
		_ = c.mergeAll(tx, entries)
		return nil
	})

}

func (c *catalog) Hash() string {
	return c.hash
}

func (c *catalog) AddEventHandler(id string, ch chan []Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.eventHandlers[id] = ch

}

func (c *catalog) RemoveEventHandler(id string) {

	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.eventHandlers, id)
}

func (c *catalog) replicate(delta []Entry) {
	select {
	case c.replicaChan <- delta:
		c.log.Debug().Msg("delta replica submitted")
	default:
		c.log.Error().Msg("replica channel blocked.")
		c.replicaChan <- delta
	}
}

func (c *catalog) hashCatalog(tx *bolt.Tx) {

	snapshot := c.entries(tx)

	h := murmur3.New128()

	for _, e := range snapshot {
		_, _ = h.Write([]byte(e.MapName))
		_, _ = h.Write([]byte(e.Key))
		_ = binary.Write(h, binary.LittleEndian, e.Time.UnixNano())
		_ = binary.Write(h, binary.LittleEndian, e.Deleted)
	}

	h1, h2 := h.Sum128()

	c.hash = fmt.Sprintf("%x%x", h1, h2)

}

func (c *catalog) emit(entries []Entry) {
	for id, h := range c.eventHandlers {
		select {
		case h <- entries:
			c.log.Debug().Msgf("dispatching event to %v", id)
		default:
			c.log.Error().Msgf("dispatcher blocks on event handler channel [%v]", id)
			h <- entries
		}
	}
}

func (c *catalog) Close() error {
	return c.db.Close()
}

func createCatalogServer(c *catalog) CatalogServer {
	return &catalogServer{
		catalog: c,
		log:     log.With().Str("component", "catalogServer").Logger(),
	}
}

type catalogServer struct {
	catalog *catalog
	log     zerolog.Logger
}

func (cs *catalogServer) Add(ctx context.Context, addRequest *AddRequest) (*AddResponse, error) {
	cs.log.Debug().Msgf("delta snapshot received %v", addRequest.Entries)
	cs.catalog.MergeSnapshot(addRequest.Entries)
	return &AddResponse{}, nil
}

func (cs *catalogServer) SnapShot(ctx context.Context, snapshotRequest *SnapshotRequest) (*SnapshotResponse, error) {
	cs.log.Debug().Msg("snapshot request received")
	r := &SnapshotResponse{
		Entries: cs.catalog.SnapShot(),
	}
	return r, nil
}

func NewCatalog(grpcSrv *grpc.Server, path string, cluster Cluster, catalogRPC CatalogRPC) (Catalog, error) {

	c, err := createCatalog(path)

	if err != nil {
		return nil, err
	}

	err = cluster.SetTag("catalog", c.Hash())

	if err != nil {
		return nil, err
	}

	cs := createCatalogServer(c)

	RegisterCatalogServer(grpcSrv, cs)

	c.replicator = createCatalogReplicator(cluster, c, c.replicaChan, catalogRPC)

	c.replicator.start()

	return c, nil

}

func createCatalog(dbPath string) (*catalog, error) {

	dbName := path.Join(dbPath, "catalog.db")

	db, err := bolt.Open(dbName, 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		return nil, err
	}

	c := &catalog{
		db:            db,
		log:           log.With().Str("component", "catalog").Logger(),
		maps:          make(map[string]map[string]Entry),
		replicaChan:   make(chan []Entry, 1024),
		eventHandlers: make(map[string]chan []Entry),
	}

	_ = c.db.View(func(tx *bolt.Tx) error {
		c.hashCatalog(tx)
		return nil
	})

	return c, nil
}

func createCatalogReplicator(cluster Cluster, catalog Catalog, deltaCh chan []Entry, catalogRPC CatalogRPC) *catalogReplicator {
	return &catalogReplicator{
		cluster:    cluster,
		catalog:    catalog,
		deltaCh:    deltaCh,
		catalogRPC: catalogRPC,
		log:        log.With().Str("component", "catalogReplicator").Logger(),
	}
}

//TODO(gvelo): split in broadcast & sync.
type catalogReplicator struct {
	cluster    Cluster
	catalog    Catalog
	deltaCh    chan []Entry
	catalogRPC CatalogRPC
	log        zerolog.Logger
}

func (cr *catalogReplicator) start() {
	cr.log.Info().Msg("starting catalog replicator")
	// TODO(gvelo): stop gracefully.
	ticker := time.NewTicker(10 * time.Second)
	go cr.broacast()
	go func() {
		for _ = range ticker.C {
			cr.sync()
		}
	}()
}

func (cr *catalogReplicator) broacast() {
	for delta := range cr.deltaCh {
		cr.log.Debug().Msg("broadcasting catalog delta")
		cr.catalogRPC.SendDelta(context.TODO(), delta)
	}
}

func (cr *catalogReplicator) sync() {

	cr.log.Info().Msg("syncing catalog")
	members := cr.cluster.LiveMembers()
	localVersion := cr.catalog.Hash()
	diff := make(map[string]string)

	cr.log.Info().Msgf("syncing catalog with local version %v", localVersion)

	for _, m := range members {

		cr.log.Debug().Msgf("processing member %v", m.Name)

		catalogVersion, ok := m.Tags["catalog"]

		if !ok {
			continue
		}

		if localVersion != catalogVersion {
			diff[catalogVersion] = m.Name
		}

	}

	if len(diff) == 0 {
		return
	}

	diffMembers := make([]string, len(diff))[:0]

	for _, m := range diff {
		diffMembers = append(diffMembers, m)
	}

	snapshots := cr.catalogRPC.GetSnapShot(context.TODO(), diffMembers)

	for _, snapshot := range snapshots {

		if snapshot.err != nil {
			cr.log.Error().Err(snapshot.err).Msgf("Error retriving snapshot from node %v", snapshot.member)
			continue
		}
		cr.catalog.MergeSnapshot(snapshot.Snapshot)
	}

	err := cr.cluster.SetTag("catalog", cr.catalog.Hash())

	if err != nil {
		cr.log.Error().Err(err).Msg("error setting catalog version tag")
	}

}
