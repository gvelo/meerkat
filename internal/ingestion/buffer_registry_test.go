package ingestion

import (
	"github.com/rs/zerolog"
	"meerkat/internal/storage"
	"meerkat/internal/util/testutil"
	"testing"
	"time"
)

func TestRegistry(t *testing.T) {

	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	bufReg := NewBufferRegistry(10,
		5*time.Second,
		-1,
		1*time.Second,
		5,
	)

	bufReg.Start()

	table := GetTable()

	for i := 0; i < 20; i++ {

		bufReg.AddToBuffer(table)
		time.Sleep(1 * time.Second)

	}

	time.Sleep(time.Hour)

}

func GetTable() *Table {

	rsw := NewRowSetWriter(0)

	tsCol := &Column{
		Idx:     0,
		Name:    storage.TSColumnName,
		ColSize: 0,
		Len:     0,
		Type:    storage.ColumnType_TIMESTAMP,
	}

	testCol := &Column{
		Idx:     1,
		Name:    "testCol",
		ColSize: 0,
		Len:     0,
		Type:    storage.ColumnType_STRING,
	}

	for i := 0; i < 10; i++ {

		rsw.WriteFixedInt64(0, time.Now().UnixNano())
		tsCol.Len++

		str := testutil.RandomString(3)
		rsw.WriteString(1, str)
		testCol.Len++
		testCol.ColSize += uint64(len(str))

	}

	partition := &Partition{
		Id:      0,
		Columns: []*Column{tsCol, testCol},
		Data:    rsw.Buf.Data(),
	}

	t := &Table{
		Name:       "testtable",
		Partitions: []*Partition{partition},
	}

	return t

}
