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

func (s storageMock) OpenSegment(info *SegmentInfo) Segment {
	args := s.Called(info)
	return args.Get(0).(Segment)
}

func (s storageMock) DeleteSegment(id uuid.UUID) {
	s.Called(id)
}

type segmentMock struct {
	mock.Mock
}

func (s *segmentMock) Info() *SegmentInfo {
	args := s.Called()
	return args.Get(0).(*SegmentInfo)
}

func (s *segmentMock) Column(name string) Column {
	args := s.Called(name)
	return args.Get(0).(Column)
}

func (s *segmentMock) Close() {
	s.Called()
}

var tests = map[string]func(t *testing.T){
	"test add segment":         testAddSegment,
	"test read write registry": testReadWrite,
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

	segmentInfo := createTestSegmentInfo()

	storageMock := &storageMock{}
	segMock := &segmentMock{}
	segMock.On("Close").Return()
	storageMock.On("OpenSegment", segmentInfo).Return(segMock)

	reg := NewSegmentRegistry(os.TempDir(), storageMock)
	reg.Start()

	reg.AddSegment(segmentInfo)

	result := reg.Segments(nil, "", "")
	assert.Len(t, result, 1)
	assert.Equal(t, segMock, result[0])
	reg.Stop()
	storageMock.AssertExpectations(t)

}

func testReadWrite(t *testing.T) {

	segmentInfo := createTestSegmentInfo()

	storageMock := &storageMock{}

	regPath := os.TempDir()
	reg := NewSegmentRegistry(regPath, storageMock)
	reg.Start()

	reg.AddSegment(segmentInfo)
	reg.Stop()

	reg = NewSegmentRegistry(regPath, storageMock)
	reg.Start()

	result := reg.SegmentInfos()

	assert.Len(t, result, 1)
	assert.Equal(t, segmentInfo.Id, result[0].Id)
	reg.Stop()
	storageMock.AssertExpectations(t)

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
