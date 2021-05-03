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
// 	"github.com/stretchr/testify/assert"
// 	"os"
// 	"path"
// 	"strconv"
// 	"testing"
// 	"time"
// )
//
// var dbpath = path.Join(os.TempDir(), "test-catalog.db")
//
// func createTestCatalog() (*catalog, error) {
// 	deleteTestCatalog()
// 	return createCatalog(dbpath)
// }
//
// func deleteTestCatalog() {
// 	_ = os.Remove(dbpath)
// }
//
// func Test_GetSet(t *testing.T) {
//
// 	assert := assert.New(t)
//
// 	c, err := createTestCatalog()
// 	assert.NoError(err)
// 	defer deleteTestCatalog()
//
// 	present := Entry{
// 		Time:    time.Now(),
// 		Value:   []byte{},
// 		Deleted: false,
// 		Key:     "test",
// 		MapName: "test",
// 	}
//
// 	eventChan := make(chan []Entry, 10)
//
// 	c.AddEventHandler("test-handler", eventChan)
//
// 	eventCount := 0
// 	replicaCount := 0
//
// 	// add new entry
//
// 	assert.Empty(c.Hash(), "hash should be empty on empty catalog")
// 	r := c.Set(present)
// 	hashPresent := c.Hash()
// 	eventCount++
// 	replicaCount++
// 	assert.NotEmpty(hashPresent, "hash should not be empty")
// 	assert.True(r, "value should be added")
// 	assert.Len(eventChan, eventCount, "event should be emmited")
// 	assert.Len(c.replicaChan, replicaCount, "event should be replicated")
//
// 	e1, found := c.Get(present.MapName, present.Key)
//
// 	assert.True(found, "entry should be on the catalog")
// 	assert.True(cmp(present, e1))
//
// 	// Test LWW Remove Bias
//
// 	r = c.Set(present)
// 	h1 := c.Hash()
// 	eventCount++
// 	replicaCount++
// 	assert.True(r, "value should be added")
// 	assert.Len(eventChan, eventCount, "event should be emmited")
// 	assert.Len(c.replicaChan, replicaCount, "event should be replicated")
// 	assert.Equal(hashPresent, h1, "hash should not be updated")
//
// 	// Test LWW (timestamp on the past).
//
// 	past := present
// 	past.Time = present.Time.Add(-time.Second)
//
// 	r = c.Set(past)
// 	h3 := c.Hash()
// 	assert.False(r, "entrie should not be added")
// 	assert.Len(eventChan, eventCount, "event should be emmited")
// 	assert.Len(c.replicaChan, replicaCount, "event should be replicated")
// 	assert.Equal(hashPresent, h3, "hash should not be updated")
//
// 	e2, found := c.Get(present.MapName, present.Key)
//
// 	assert.True(found, "entry should be on the catalog")
// 	assert.True(cmp(present, e2))
//
// 	// Test LWW (timestamp on the future).
//
// 	future := present
// 	future.Time = present.Time.Add(time.Second)
//
// 	r = c.Set(future)
// 	h4 := c.Hash()
// 	eventCount++
// 	replicaCount++
// 	assert.True(r, "entrie should be added")
// 	assert.Len(eventChan, eventCount, "event should be emmited")
// 	assert.Len(c.replicaChan, replicaCount, "event should be replicated")
// 	assert.NotEqual(h3, h4, "hash should be updated")
//
// 	e3, found := c.Get(future.MapName, future.Key)
//
// 	assert.True(found, "entry should be on the catalog")
// 	assert.True(cmp(future, e3))
//
// 	// Test LWW remove (timestamp on the future).
//
// 	future1 := future
// 	future1.Time = future.Time.Add(time.Second)
// 	future1.Deleted = true
//
// 	r = c.Set(future1)
// 	h5 := c.Hash()
// 	eventCount++
// 	replicaCount++
// 	assert.True(r, "entrie should be added")
// 	assert.Len(eventChan, eventCount, "event should be emmited")
// 	assert.Len(c.replicaChan, replicaCount, "event should be replicated")
// 	assert.NotEqual(h4, h5, "hash should be updated")
//
// 	_, found = c.Get(future1.MapName, future1.Key)
//
// 	assert.False(found, "entry should not be on the catalog")
//
// 	// Test LWW add before remove (timestamp on the future).
//
// 	future2 := future1
// 	future2.Time = future1.Time.Add(time.Second)
// 	future2.Deleted = false
//
// 	r = c.Set(future2)
// 	h6 := c.Hash()
// 	assert.False(r, "entrie should be added")
// 	assert.Len(eventChan, eventCount, "event should be emmited")
// 	assert.Len(c.replicaChan, replicaCount, "event should be replicated")
// 	assert.Equal(h5, h6, "hash should not be updated")
//
// 	_, found = c.Get(future2.MapName, future2.Key)
//
// 	assert.False(found, "entry should not be on the catalog")
//
// }
//
// func Test_GetSetAll(t *testing.T) {
//
// 	assert := assert.New(t)
//
// 	c, err := createTestCatalog()
// 	assert.NoError(err)
// 	defer deleteTestCatalog()
//
// 	eventChan := make(chan []Entry, 10)
// 	c.AddEventHandler("test-handler", eventChan)
//
// 	entries1 := createEntries("test1", 10)
// 	entries2 := createEntries("test1", 10)
//
// 	c.SetAll(entries1)
// 	c.SetAll(entries2)
//
// 	r1 := c.GetAll("test1")
// 	r2 := c.GetAll("test1")
// 	s := c.SnapShot()
//
// 	assert.NotEmpty(c.Hash(), "catalog content should be hashed.")
// 	assert.Len(r1, 10, "map test1 should contains 10 entries")
// 	assert.Len(r2, 10, "map test1 should contains 10 entries")
// 	assert.Len(s, 10, "snapshot should contains 10 entries")
// 	assert.Len(eventChan, 2, "event should be emmited")
// 	assert.Len(c.replicaChan, 2, "event should be replicated")
//
// }
//
// func Test_MergeSnapshot(t *testing.T) {
//
// 	assert := assert.New(t)
//
// 	c, err := createTestCatalog()
// 	assert.NoError(err)
// 	defer deleteTestCatalog()
//
// 	eventChan := make(chan []Entry, 10)
// 	c.AddEventHandler("test-handler", eventChan)
//
// 	snapshot := createEntries("test1", 10)
//
// 	c.MergeSnapshot(snapshot)
//
// 	s := c.SnapShot()
//
// 	assert.NotEmpty(c.Hash(), "catalog content should be hashed.")
// 	assert.Len(s, 10, "snapshot should contains 10 entries")
// 	assert.Len(eventChan, 1, "event should be emmited")
// 	assert.Len(c.replicaChan, 0, "event should be replicated")
//
// }
//
// func cmp(e1 Entry, e2 Entry) bool {
// 	return e1.Time.Equal(e2.Time) &&
// 		e1.Deleted == e2.Deleted &&
// 		e1.MapName == e2.MapName &&
// 		e1.Key == e2.Key
// }
//
// func createEntries(mapName string, count int) []Entry {
//
// 	var entries []Entry
// 	t := time.Now()
//
// 	for i := 0; i < count; i++ {
// 		e := Entry{
// 			Key:     strconv.Itoa(i),
// 			MapName: mapName,
// 			Deleted: false,
// 			Time:    t,
// 			Value:   []byte{},
// 		}
// 		entries = append(entries, e)
// 	}
// 	return entries
// }
