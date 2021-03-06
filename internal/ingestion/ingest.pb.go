// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ingest.proto

package ingestion

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	io "io"
	math "math"
	storage "meerkat/internal/storage"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type IngestionRequest struct {
	Table *Table `protobuf:"bytes,1,opt,name=table,proto3" json:"table,omitempty"`
}

func (m *IngestionRequest) Reset()         { *m = IngestionRequest{} }
func (m *IngestionRequest) String() string { return proto.CompactTextString(m) }
func (*IngestionRequest) ProtoMessage()    {}
func (*IngestionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ff993cce43359ffa, []int{0}
}
func (m *IngestionRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *IngestionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_IngestionRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *IngestionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IngestionRequest.Merge(m, src)
}
func (m *IngestionRequest) XXX_Size() int {
	return m.Size()
}
func (m *IngestionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_IngestionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_IngestionRequest proto.InternalMessageInfo

func (m *IngestionRequest) GetTable() *Table {
	if m != nil {
		return m.Table
	}
	return nil
}

type IngestResponse struct {
}

func (m *IngestResponse) Reset()         { *m = IngestResponse{} }
func (m *IngestResponse) String() string { return proto.CompactTextString(m) }
func (*IngestResponse) ProtoMessage()    {}
func (*IngestResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ff993cce43359ffa, []int{1}
}
func (m *IngestResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *IngestResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_IngestResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *IngestResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IngestResponse.Merge(m, src)
}
func (m *IngestResponse) XXX_Size() int {
	return m.Size()
}
func (m *IngestResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_IngestResponse.DiscardUnknown(m)
}

var xxx_messageInfo_IngestResponse proto.InternalMessageInfo

type Table struct {
	Name       string       `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Partitions []*Partition `protobuf:"bytes,2,rep,name=partitions,proto3" json:"partitions,omitempty"`
}

func (m *Table) Reset()         { *m = Table{} }
func (m *Table) String() string { return proto.CompactTextString(m) }
func (*Table) ProtoMessage()    {}
func (*Table) Descriptor() ([]byte, []int) {
	return fileDescriptor_ff993cce43359ffa, []int{2}
}
func (m *Table) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Table) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Table.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Table) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Table.Merge(m, src)
}
func (m *Table) XXX_Size() int {
	return m.Size()
}
func (m *Table) XXX_DiscardUnknown() {
	xxx_messageInfo_Table.DiscardUnknown(m)
}

var xxx_messageInfo_Table proto.InternalMessageInfo

func (m *Table) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Table) GetPartitions() []*Partition {
	if m != nil {
		return m.Partitions
	}
	return nil
}

type Partition struct {
	Id      uint64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Columns []*Column `protobuf:"bytes,2,rep,name=columns,proto3" json:"columns,omitempty"`
	Data    []byte    `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *Partition) Reset()         { *m = Partition{} }
func (m *Partition) String() string { return proto.CompactTextString(m) }
func (*Partition) ProtoMessage()    {}
func (*Partition) Descriptor() ([]byte, []int) {
	return fileDescriptor_ff993cce43359ffa, []int{3}
}
func (m *Partition) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Partition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Partition.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Partition) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Partition.Merge(m, src)
}
func (m *Partition) XXX_Size() int {
	return m.Size()
}
func (m *Partition) XXX_DiscardUnknown() {
	xxx_messageInfo_Partition.DiscardUnknown(m)
}

var xxx_messageInfo_Partition proto.InternalMessageInfo

func (m *Partition) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Partition) GetColumns() []*Column {
	if m != nil {
		return m.Columns
	}
	return nil
}

func (m *Partition) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type Column struct {
	Idx     uint64             `protobuf:"varint,1,opt,name=idx,proto3" json:"idx,omitempty"`
	Name    string             `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	ColSize uint64             `protobuf:"varint,3,opt,name=colSize,proto3" json:"colSize,omitempty"`
	Len     uint64             `protobuf:"varint,4,opt,name=len,proto3" json:"len,omitempty"`
	Type    storage.ColumnType `protobuf:"varint,5,opt,name=type,proto3,enum=meerkat.storage.ColumnType" json:"type,omitempty"`
}

func (m *Column) Reset()         { *m = Column{} }
func (m *Column) String() string { return proto.CompactTextString(m) }
func (*Column) ProtoMessage()    {}
func (*Column) Descriptor() ([]byte, []int) {
	return fileDescriptor_ff993cce43359ffa, []int{4}
}
func (m *Column) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Column) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Column.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Column) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Column.Merge(m, src)
}
func (m *Column) XXX_Size() int {
	return m.Size()
}
func (m *Column) XXX_DiscardUnknown() {
	xxx_messageInfo_Column.DiscardUnknown(m)
}

var xxx_messageInfo_Column proto.InternalMessageInfo

func (m *Column) GetIdx() uint64 {
	if m != nil {
		return m.Idx
	}
	return 0
}

func (m *Column) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Column) GetColSize() uint64 {
	if m != nil {
		return m.ColSize
	}
	return 0
}

func (m *Column) GetLen() uint64 {
	if m != nil {
		return m.Len
	}
	return 0
}

func (m *Column) GetType() storage.ColumnType {
	if m != nil {
		return m.Type
	}
	return storage.ColumnType_TIMESTAMP
}

func init() {
	proto.RegisterType((*IngestionRequest)(nil), "meerkat.ingestion.IngestionRequest")
	proto.RegisterType((*IngestResponse)(nil), "meerkat.ingestion.IngestResponse")
	proto.RegisterType((*Table)(nil), "meerkat.ingestion.Table")
	proto.RegisterType((*Partition)(nil), "meerkat.ingestion.Partition")
	proto.RegisterType((*Column)(nil), "meerkat.ingestion.Column")
}

func init() { proto.RegisterFile("ingest.proto", fileDescriptor_ff993cce43359ffa) }

var fileDescriptor_ff993cce43359ffa = []byte{
	// 383 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0xb1, 0x8e, 0xd3, 0x40,
	0x10, 0x86, 0xbd, 0x8e, 0x13, 0xc8, 0x24, 0x44, 0x61, 0x45, 0x61, 0x02, 0x32, 0xc6, 0x69, 0x5c,
	0xd9, 0x92, 0xd3, 0x52, 0x85, 0x8a, 0x0e, 0x2d, 0x69, 0xa0, 0x40, 0xda, 0xe0, 0xc1, 0x5a, 0x61,
	0x7b, 0x8d, 0xbd, 0x91, 0x08, 0xaf, 0x40, 0xc3, 0x63, 0x5d, 0x99, 0xf2, 0xca, 0x53, 0xf2, 0x22,
	0x27, 0xef, 0xc6, 0x56, 0x74, 0x97, 0xeb, 0xfe, 0xc9, 0x7c, 0xf3, 0xff, 0x3b, 0x19, 0xc3, 0x54,
	0x94, 0x19, 0x36, 0x2a, 0xaa, 0x6a, 0xa9, 0x24, 0x7d, 0x59, 0x20, 0xd6, 0xbf, 0xb8, 0x8a, 0xcc,
	0xaf, 0x42, 0x96, 0x8b, 0x57, 0x99, 0xcc, 0xa4, 0xee, 0xc6, 0xad, 0x32, 0xe0, 0xe2, 0x5d, 0x26,
	0x65, 0x96, 0x63, 0xac, 0xab, 0xed, 0xee, 0x67, 0xac, 0x44, 0x81, 0x8d, 0xe2, 0x45, 0x75, 0x06,
	0x5e, 0x34, 0x4a, 0xd6, 0x3c, 0x43, 0x53, 0x06, 0x6b, 0x98, 0x7f, 0xea, 0x2c, 0x19, 0xfe, 0xde,
	0x61, 0xa3, 0x68, 0x04, 0x43, 0xc5, 0xb7, 0x39, 0xba, 0xc4, 0x27, 0xe1, 0x24, 0x71, 0xa3, 0x47,
	0xe1, 0xd1, 0xa6, 0xed, 0x33, 0x83, 0x05, 0x73, 0x98, 0x19, 0x0f, 0x86, 0x4d, 0x25, 0xcb, 0x06,
	0x83, 0xaf, 0x30, 0xd4, 0x04, 0xa5, 0xe0, 0x94, 0xbc, 0x30, 0x4e, 0x63, 0xa6, 0x35, 0xfd, 0x00,
	0x50, 0xf1, 0x5a, 0x89, 0xd6, 0xa8, 0x71, 0x6d, 0x7f, 0x10, 0x4e, 0x92, 0xb7, 0x57, 0x32, 0x3e,
	0x77, 0x10, 0xbb, 0xe0, 0x83, 0x14, 0xc6, 0x7d, 0x83, 0xce, 0xc0, 0x16, 0xa9, 0x36, 0x77, 0x98,
	0x2d, 0x52, 0xba, 0x82, 0x67, 0x3f, 0x64, 0xbe, 0x2b, 0x7a, 0xdf, 0xd7, 0x57, 0x7c, 0x3f, 0x6a,
	0x82, 0x75, 0x64, 0xfb, 0xc6, 0x94, 0x2b, 0xee, 0x0e, 0x7c, 0x12, 0x4e, 0x99, 0xd6, 0xc1, 0x3f,
	0x02, 0x23, 0xc3, 0xd1, 0x39, 0x0c, 0x44, 0xfa, 0xe7, 0x1c, 0xd2, 0xca, 0x7e, 0x29, 0xfb, 0x62,
	0x29, 0x57, 0x27, 0x7f, 0x11, 0x7f, 0x51, 0xfb, 0x38, 0xac, 0x2b, 0xdb, 0xf9, 0x1c, 0x4b, 0xd7,
	0x31, 0xf3, 0x39, 0x96, 0x34, 0x06, 0x47, 0xed, 0x2b, 0x74, 0x87, 0x3e, 0x09, 0x67, 0xc9, 0x9b,
	0xfe, 0x89, 0xdd, 0x65, 0x4c, 0xf0, 0x66, 0x5f, 0x21, 0xd3, 0x60, 0xf2, 0x1d, 0x9e, 0x9b, 0x3f,
	0x18, 0x6b, 0xca, 0x60, 0x64, 0x34, 0x5d, 0x5e, 0xd9, 0xed, 0xe1, 0x2d, 0x17, 0xef, 0x9f, 0x84,
	0xfa, 0x63, 0x59, 0xeb, 0xe5, 0xcd, 0xd1, 0x23, 0x87, 0xa3, 0x47, 0xee, 0x8e, 0x1e, 0xf9, 0x7f,
	0xf2, 0xac, 0xc3, 0xc9, 0xb3, 0x6e, 0x4f, 0x9e, 0xf5, 0x6d, 0xdc, 0x4f, 0x6d, 0x47, 0xfa, 0x83,
	0x59, 0xdd, 0x07, 0x00, 0x00, 0xff, 0xff, 0x4c, 0xe0, 0x5f, 0x00, 0x99, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// IngesterClient is the client API for Ingester service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type IngesterClient interface {
	Ingest(ctx context.Context, in *IngestionRequest, opts ...grpc.CallOption) (*IngestResponse, error)
}

type ingesterClient struct {
	cc *grpc.ClientConn
}

func NewIngesterClient(cc *grpc.ClientConn) IngesterClient {
	return &ingesterClient{cc}
}

func (c *ingesterClient) Ingest(ctx context.Context, in *IngestionRequest, opts ...grpc.CallOption) (*IngestResponse, error) {
	out := new(IngestResponse)
	err := c.cc.Invoke(ctx, "/meerkat.ingestion.Ingester/Ingest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IngesterServer is the server API for Ingester service.
type IngesterServer interface {
	Ingest(context.Context, *IngestionRequest) (*IngestResponse, error)
}

func RegisterIngesterServer(s *grpc.Server, srv IngesterServer) {
	s.RegisterService(&_Ingester_serviceDesc, srv)
}

func _Ingester_Ingest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IngestionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IngesterServer).Ingest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meerkat.ingestion.Ingester/Ingest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IngesterServer).Ingest(ctx, req.(*IngestionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Ingester_serviceDesc = grpc.ServiceDesc{
	ServiceName: "meerkat.ingestion.Ingester",
	HandlerType: (*IngesterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ingest",
			Handler:    _Ingester_Ingest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ingest.proto",
}

func (m *IngestionRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *IngestionRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Table != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintIngest(dAtA, i, uint64(m.Table.Size()))
		n1, err := m.Table.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	return i, nil
}

func (m *IngestResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *IngestResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	return i, nil
}

func (m *Table) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Table) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Name) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintIngest(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	if len(m.Partitions) > 0 {
		for _, msg := range m.Partitions {
			dAtA[i] = 0x12
			i++
			i = encodeVarintIngest(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *Partition) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Partition) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Id != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintIngest(dAtA, i, uint64(m.Id))
	}
	if len(m.Columns) > 0 {
		for _, msg := range m.Columns {
			dAtA[i] = 0x12
			i++
			i = encodeVarintIngest(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Data) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintIngest(dAtA, i, uint64(len(m.Data)))
		i += copy(dAtA[i:], m.Data)
	}
	return i, nil
}

func (m *Column) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Column) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Idx != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintIngest(dAtA, i, uint64(m.Idx))
	}
	if len(m.Name) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintIngest(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	if m.ColSize != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintIngest(dAtA, i, uint64(m.ColSize))
	}
	if m.Len != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintIngest(dAtA, i, uint64(m.Len))
	}
	if m.Type != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintIngest(dAtA, i, uint64(m.Type))
	}
	return i, nil
}

func encodeVarintIngest(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *IngestionRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Table != nil {
		l = m.Table.Size()
		n += 1 + l + sovIngest(uint64(l))
	}
	return n
}

func (m *IngestResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *Table) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovIngest(uint64(l))
	}
	if len(m.Partitions) > 0 {
		for _, e := range m.Partitions {
			l = e.Size()
			n += 1 + l + sovIngest(uint64(l))
		}
	}
	return n
}

func (m *Partition) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovIngest(uint64(m.Id))
	}
	if len(m.Columns) > 0 {
		for _, e := range m.Columns {
			l = e.Size()
			n += 1 + l + sovIngest(uint64(l))
		}
	}
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovIngest(uint64(l))
	}
	return n
}

func (m *Column) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Idx != 0 {
		n += 1 + sovIngest(uint64(m.Idx))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovIngest(uint64(l))
	}
	if m.ColSize != 0 {
		n += 1 + sovIngest(uint64(m.ColSize))
	}
	if m.Len != 0 {
		n += 1 + sovIngest(uint64(m.Len))
	}
	if m.Type != 0 {
		n += 1 + sovIngest(uint64(m.Type))
	}
	return n
}

func sovIngest(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozIngest(x uint64) (n int) {
	return sovIngest(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *IngestionRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIngest
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: IngestionRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IngestionRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Table", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIngest
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthIngest
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthIngest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Table == nil {
				m.Table = &Table{}
			}
			if err := m.Table.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipIngest(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthIngest
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthIngest
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *IngestResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIngest
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: IngestResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IngestResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipIngest(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthIngest
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthIngest
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Table) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIngest
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Table: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Table: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIngest
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIngest
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthIngest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Partitions", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIngest
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthIngest
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthIngest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Partitions = append(m.Partitions, &Partition{})
			if err := m.Partitions[len(m.Partitions)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipIngest(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthIngest
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthIngest
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Partition) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIngest
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Partition: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Partition: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIngest
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Columns", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIngest
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthIngest
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthIngest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Columns = append(m.Columns, &Column{})
			if err := m.Columns[len(m.Columns)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIngest
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthIngest
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthIngest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = append(m.Data[:0], dAtA[iNdEx:postIndex]...)
			if m.Data == nil {
				m.Data = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipIngest(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthIngest
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthIngest
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Column) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIngest
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Column: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Column: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Idx", wireType)
			}
			m.Idx = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIngest
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Idx |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIngest
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIngest
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthIngest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ColSize", wireType)
			}
			m.ColSize = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIngest
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ColSize |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Len", wireType)
			}
			m.Len = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIngest
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Len |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIngest
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= storage.ColumnType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipIngest(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthIngest
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthIngest
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipIngest(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowIngest
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowIngest
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowIngest
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthIngest
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthIngest
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowIngest
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipIngest(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthIngest
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthIngest = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowIngest   = fmt.Errorf("proto: integer overflow")
)
