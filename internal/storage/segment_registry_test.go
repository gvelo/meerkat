package storage

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"path"
	"testing"
	"time"
)

type storageMock struct {
	mock.Mock
}

func (s storageMock) WriteSegment(src SegmentSource) *SegmentInfo {
	args := s.Called(src)
	return args.Get(0).(*SegmentInfo)
}

func (s storageMock) OpenSegment(info *SegmentInfo) SegmentIF {
	args := s.Called(info)
	return args.Get(0).(*Segment)
}

func (s storageMock) DeleteSegment(id uuid.UUID) {
	s.Called(id)
}

var tests = map[string]func(t *testing.T){
	"test add segment": testAddSegment,
}

func TestSegmentRegistry(t *testing.T) {
	for name, f := range tests {
		deleteRegistry()
		t.Run(name, f)
	}
}

func deleteRegistry() {
	_ = os.Remove(path.Join(os.TempDir(), registryFileName))
}

func testAddSegment(t *testing.T) {

	storageMock := &storageMock{}

	reg := NewSegmentRegistry(os.TempDir(), storageMock)
	reg.Start()

	segmentInfo := createTestSegmentInfo()

	storageMock.On("OpenSegment", segmentInfo).Return(&Segment{})

	reg.AddSegment(segmentInfo)

	reg.Stop()

	result := reg.Segments(nil, "", "")

	assert.Len(t, result, 1)
	assert.Equal(t, &Segment{}, result[0])

}

func createTestSegmentInfo() *SegmentInfo {

	id := uuid.New()

	return &SegmentInfo{
		Id:           id[:],
		DatabaseName: "test-db",
		TableName:    "test-table",
		PartitionId:  0,
		Len:          0,
		Interval: &Interval{
			From: time.Now(),
			To:   time.Now(),
		},
		Columns: []*ColumnInfo{
			{
				Name:           "test-columnt",
				ColumnType:     ColumnType_STRING,
				IndexType:      IndexType_FULLTEXT,
				Encoding:       Encoding_SNAPPY,
				Nullable:       false,
				Len:            12344323,
				Cardinality:    23,
				SizeOnDisk:     413341,
				CompressedSize: 13433,
				NullCount:      0,
			},
		},
	}

}
