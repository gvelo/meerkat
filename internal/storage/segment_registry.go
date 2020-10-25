package storage

import (
	"errors"
	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	bolt "go.etcd.io/bbolt"
	"path"
	"sync"
	"time"
)

const (
	registryFileName = "segreg.db"
	bucketName       = "segments"
	cleanInterval    = 5 * time.Second
)

type SegmentRegistry interface {
	SegmentInfos() []*SegmentInfo
	Segments(partitions []uint64, dbName string, tableName string) []SegmentIF
	AddSegment(segmentInfo *SegmentInfo)
	RemoveSegment(segmentId uuid.UUID)
	MergeSegment(segmentInfo *SegmentInfo, mergedSegments []uuid.UUID)
	Release(segmentId uuid.UUID)
	Start()
	Stop()
}

type segmentEntry struct {
	id          uuid.UUID
	segmentInfo *SegmentInfo
	deleted     bool
	refCount    int
	segment     SegmentIF
	mu          sync.Mutex
}

type segmentRegistry struct {
	cache       map[uuid.UUID]*segmentEntry
	db          *bolt.DB
	mu          sync.Mutex
	segStorage  SegmentStorage
	cleanTicker *time.Ticker
	done        chan struct{}
	wg          sync.WaitGroup
	running     bool
	log         zerolog.Logger
}

func NewSegmentRegistry(dbPath string, segStorage SegmentStorage) SegmentRegistry {

	dbName := path.Join(dbPath, registryFileName)

	db, err := bolt.Open(dbName, 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		panic(err)
	}

	reg := &segmentRegistry{
		cache:      make(map[uuid.UUID]*segmentEntry),
		db:         db,
		segStorage: segStorage,
		done:       make(chan struct{}),
		log:        log.With().Str("component", "SegmentRegistry").Logger(),
	}

	reg.read()

	return reg
}

func (s *segmentRegistry) read() {

	s.log.Info().Msg("reading segments info")

	s.updateDb(func(bucket *bolt.Bucket) {

		_ = bucket.ForEach(func(k, v []byte) error {

			segmentInfo := &SegmentInfo{}

			err := proto.Unmarshal(v, segmentInfo)

			if err != nil {
				s.log.Panic().Err(err).Msg("error unmarshalling SegmentInfo")
			}

			id, err := uuid.FromBytes(segmentInfo.Id)

			if err != nil {
				s.log.Panic().Err(err).Msg("error unmarshalling SegmentId")
			}

			s.cache[id] = &segmentEntry{
				id:          id,
				segmentInfo: segmentInfo,
				deleted:     false,
				refCount:    0,
				segment:     nil,
			}

			return nil

		})

	})

}

func (s *segmentRegistry) SegmentInfos() []*SegmentInfo {

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		panic(errors.New("segment registry not running"))
	}

	segments := make([]*SegmentInfo, len(s.cache))

	for _, entry := range s.cache {
		segments = append(segments, entry.segmentInfo)
	}

	return segments

}

func (s *segmentRegistry) Segments(partitions []uint64, dbName string, tableName string) []SegmentIF {

	s.log.Debug().Msg("getSegments")

	entries := s.filterSegments(partitions, dbName, tableName)

	s.openSegments(entries)

	segments := make([]SegmentIF, len(entries))

	for i, entry := range entries {
		segments[i] = entry.segment
	}

	return segments

}

func (s *segmentRegistry) openSegments(entries []*segmentEntry) {

	defer func() {
		if r := recover(); r != nil {
			s.releaseAll(entries)
		}
	}()
	for _, entry := range entries {
		s.openSegment(entry)
	}
}

func (s *segmentRegistry) releaseAll(entries []*segmentEntry) {

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, entry := range entries {
		s.release(entry)
	}
}

func (s *segmentRegistry) release(entry *segmentEntry) {
	s.log.Debug().Str("segmentId", entry.id.String()).Msg("release")

	if entry.refCount == 0 {
		s.log.Panic().Str("segmentId", entry.id.String()).Msg("segment refCount = 0")
		return
	}

	entry.refCount--

}

func (s *segmentRegistry) openSegment(entry *segmentEntry) {

	entry.mu.Lock()
	defer entry.mu.Unlock()

	if entry.segment == nil {
		entry.segment = s.segStorage.OpenSegment(entry.segmentInfo)
	}

}

func (s *segmentRegistry) filterSegments(partitions []uint64, dbName string, tableName string) []*segmentEntry {

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		panic(errors.New("segment registry not running"))
	}

	// TODO(gvelo): implement filtering.

	var entries []*segmentEntry

	for _, entry := range s.cache {

		if entry.deleted {
			continue
		}

		entry.refCount++

		entries = append(entries, entry)

	}

	return entries
}

func (s *segmentRegistry) AddSegment(segmentInfo *SegmentInfo) {

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		panic(errors.New("segment registry not running"))
	}

	s.log.Debug().Msg("AddSegment")

	bytes, err := proto.Marshal(segmentInfo)

	if err != nil {
		s.log.Panic().Err(err).Msg("cannot marshal SegmentInfo")
	}

	s.updateDb(func(bucket *bolt.Bucket) {

		err := bucket.Put(segmentInfo.Id, bytes)

		if err != nil {
			s.log.Panic().Err(err).Msg("cannot store SegmentInfo")
		}

	})

	id, err := uuid.FromBytes(segmentInfo.Id)

	if err != nil {
		s.log.Panic().Err(err).Msg("cannot unmarshal uuid")
	}

	entry := &segmentEntry{
		id:          id,
		segmentInfo: segmentInfo,
	}

	s.cache[id] = entry

}

func (s *segmentRegistry) RemoveSegment(segmentId uuid.UUID) {

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		panic(errors.New("segment registry not running"))
	}

	s.updateDb(func(bucket *bolt.Bucket) {

		err := bucket.Delete(segmentId[:])

		if err != nil {
			s.log.Panic().Err(err).Msg("cannot remove SegmentInfo")
		}

	})

	entry := s.cache[segmentId]

	if entry == nil {
		s.log.Panic().Str("segmentId", segmentId.String()).Msg("segment not found")
		return
	}

	entry.deleted = true

}

func (s *segmentRegistry) MergeSegment(segmentInfo *SegmentInfo, mergedSegments []uuid.UUID) {

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		panic(errors.New("segment registry not running"))
	}

	s.updateDb(func(bucket *bolt.Bucket) {

		for _, mergedSegmentId := range mergedSegments {

			err := bucket.Delete(mergedSegmentId[:])

			if err != nil {
				s.log.Panic().Err(err).Msg("cannot remove SegmentInfo")
			}

		}

		bytes, err := proto.Marshal(segmentInfo)

		if err != nil {
			s.log.Panic().Err(err).Msg("cannot marshal SegmentInfo")
		}

		err = bucket.Put(segmentInfo.Id, bytes)

		if err != nil {
			s.log.Panic().Err(err).Msg("cannot store SegmentInfo")
		}

	})

	id, err := uuid.FromBytes(segmentInfo.Id)

	if err != nil {
		s.log.Panic().Err(err).Msg("cannot unmarshal uuid")
	}

	entry := &segmentEntry{
		id:          id,
		segmentInfo: segmentInfo,
	}

	for _, mergedSegmentId := range mergedSegments {

		entry := s.cache[mergedSegmentId]

		if entry == nil {
			s.log.Panic().Str("segmentId", mergedSegmentId.String()).Msg("segment not found")
			return
		}

		entry.deleted = true
	}

	s.cache[id] = entry

}

func (s *segmentRegistry) Release(segmentId uuid.UUID) {

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		panic(errors.New("segment registry not running"))
	}

	entry := s.cache[segmentId]

	if entry == nil {
		s.log.Panic().Str("segmentId", segmentId.String()).Msg("segment not found")
		return
	}

	s.release(entry)
}

func (s *segmentRegistry) clean() {

	s.mu.Lock()

	var deletedSegments []*segmentEntry

	for id, entry := range s.cache {
		if entry.deleted && entry.refCount == 0 {
			delete(s.cache, id)
			deletedSegments = append(deletedSegments, entry)
		}
	}

	s.mu.Unlock()

	for _, entry := range deletedSegments {
		entry.mu.Lock()
		if entry.segment != nil {
			entry.segment.Close()
			s.log.Debug().Str("segmentId", entry.id.String()).Msg("deleting segment")
			s.segStorage.DeleteSegment(entry.id)
		}
		entry.mu.Unlock()
	}
}

func (s *segmentRegistry) updateDb(f func(bucket *bolt.Bucket)) {

	tx, err := s.db.Begin(true)

	if err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			err := tx.Rollback()
			s.log.Error().Err(err).Msg("cannot rollback tx")
			panic(r)
		}
		err := tx.Commit()
		if err != nil {
			panic(err)
		}
	}()

	bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))

	if err != nil {
		panic(err)
	}

	f(bucket)

}

func (s *segmentRegistry) run() {

	defer s.wg.Done()

	for {
		select {
		case <-s.cleanTicker.C:
			s.clean()
		case <-s.done:
			s.clean()
			return
		}
	}

}

func (s *segmentRegistry) Start() {

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return
	}

	s.running = true

	s.cleanTicker = time.NewTicker(cleanInterval)
	s.wg.Add(1)
	go s.run()

}

func (s *segmentRegistry) Stop() {

	s.mu.Lock()

	if !s.running {
		return
	}

	s.running = false

	s.mu.Unlock()

	err := s.db.Close()

	if err != nil {
		s.log.Err(err).Msg("closing db")
	}

	s.cleanTicker.Stop()
	close(s.done)
	s.wg.Wait()

	for _, entry := range s.cache {
		entry.segment.Close()
	}

}
