// Code generated by protoc-gen-go. DO NOT EDIT.
// source: planets-mongo.proto

package v1

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Planet struct {
	Name                 string               `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty" bson:"name"`
	Id                   int32                `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty" bson:"id"`
	Climate              string               `protobuf:"bytes,3,opt,name=climate,proto3" json:"climate,omitempty" bson:"climate"`
	Terrain              []string             `protobuf:"bytes,4,rep,name=terrain,proto3" json:"terrain,omitempty" bson:"terrain"`
	TimesOnMovies        int32                `protobuf:"varint,5,opt,name=timesOnMovies,proto3" json:"timesOnMovies,omitempty" bson:"times_on_movies"`
	AddedAt              *timestamp.Timestamp `protobuf:"bytes,6,opt,name=added_at,json=addedAt,proto3" json:"added_at,omitempty" bson:"added_at"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Planet) Reset()         { *m = Planet{} }
func (m *Planet) String() string { return proto.CompactTextString(m) }
func (*Planet) ProtoMessage()    {}
func (*Planet) Descriptor() ([]byte, []int) {
	return fileDescriptor_ecba292957595039, []int{0}
}

func (m *Planet) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Planet.Unmarshal(m, b)
}
func (m *Planet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Planet.Marshal(b, m, deterministic)
}
func (m *Planet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Planet.Merge(m, src)
}
func (m *Planet) XXX_Size() int {
	return xxx_messageInfo_Planet.Size(m)
}
func (m *Planet) XXX_DiscardUnknown() {
	xxx_messageInfo_Planet.DiscardUnknown(m)
}

var xxx_messageInfo_Planet proto.InternalMessageInfo

func (m *Planet) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Planet) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Planet) GetClimate() string {
	if m != nil {
		return m.Climate
	}
	return ""
}

func (m *Planet) GetTerrain() []string {
	if m != nil {
		return m.Terrain
	}
	return nil
}

func (m *Planet) GetTimesOnMovies() int32 {
	if m != nil {
		return m.TimesOnMovies
	}
	return 0
}

func (m *Planet) GetAddedAt() *timestamp.Timestamp {
	if m != nil {
		return m.AddedAt
	}
	return nil
}

func init() {
	proto.RegisterType((*Planet)(nil), "v1.Planet")
}

func init() { proto.RegisterFile("planets-mongo.proto", fileDescriptor_ecba292957595039) }

var fileDescriptor_ecba292957595039 = []byte{
	// 207 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x8e, 0x4b, 0x4a, 0xc6, 0x30,
	0x10, 0xc7, 0x49, 0xfa, 0xd2, 0x88, 0x2e, 0xe2, 0x26, 0x74, 0x63, 0x10, 0x17, 0xdd, 0x98, 0xa2,
	0xe2, 0x01, 0x3c, 0x80, 0x28, 0xc5, 0xbd, 0xa4, 0x66, 0x2c, 0x81, 0x26, 0x29, 0xed, 0xd8, 0xdb,
	0x79, 0x37, 0xe9, 0xd4, 0x2e, 0xbe, 0xdd, 0xfc, 0x1f, 0x33, 0xf3, 0x13, 0xd7, 0xd3, 0x68, 0x23,
	0xe0, 0x72, 0x1f, 0x52, 0x1c, 0x92, 0x99, 0xe6, 0x84, 0x49, 0xf2, 0xf5, 0xa1, 0xbe, 0x19, 0x52,
	0x1a, 0x46, 0x68, 0xc9, 0xe9, 0x7f, 0xbe, 0x5b, 0xf4, 0x01, 0x16, 0xb4, 0x61, 0xda, 0x4b, 0xb7,
	0xbf, 0x4c, 0x94, 0xef, 0xb4, 0x2c, 0xa5, 0xc8, 0xa3, 0x0d, 0xa0, 0x98, 0x66, 0xcd, 0x79, 0x47,
	0xb3, 0xbc, 0x12, 0xdc, 0x3b, 0xc5, 0x35, 0x6b, 0x8a, 0x8e, 0x7b, 0x27, 0x95, 0xa8, 0xbe, 0x46,
	0x1f, 0x2c, 0x82, 0xca, 0xa8, 0x76, 0xc8, 0x2d, 0x41, 0x98, 0x67, 0xeb, 0xa3, 0xca, 0x75, 0xb6,
	0x25, 0xff, 0x52, 0xde, 0x89, 0x4b, 0xfa, 0xfa, 0x16, 0x5f, 0xd3, 0xea, 0x61, 0x51, 0x05, 0x9d,
	0x3b, 0x35, 0xe5, 0xb3, 0x38, 0xb3, 0xce, 0x81, 0xfb, 0xb4, 0xa8, 0x4a, 0xcd, 0x9a, 0x8b, 0xc7,
	0xda, 0xec, 0xf0, 0xe6, 0x80, 0x37, 0x1f, 0x07, 0x7c, 0x57, 0x51, 0xf7, 0x05, 0xfb, 0x92, 0xc2,
	0xa7, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x48, 0x49, 0xab, 0xd6, 0x02, 0x01, 0x00, 0x00,
}
