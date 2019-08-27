// Code generated by protoc-gen-go. DO NOT EDIT.
// source: planets-mongo.proto

package mongov1

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	math "math"
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
	Name                 string               `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Id                   int32                `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	Climate              string               `protobuf:"bytes,3,opt,name=climate,proto3" json:"climate,omitempty"`
	Terrain              []string             `protobuf:"bytes,4,rep,name=terrain,proto3" json:"terrain,omitempty"`
	TimesOnMovies        int32                `protobuf:"varint,5,opt,name=timesOnMovies,proto3" json:"timesOnMovies,omitempty"`
	AddedAt              *timestamp.Timestamp `protobuf:"bytes,6,opt,name=added_at,json=addedAt,proto3" json:"added_at,omitempty"`
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
	proto.RegisterType((*Planet)(nil), "service.Planet")
}

func init() { proto.RegisterFile("planets-mongo.proto", fileDescriptor_ecba292957595039) }

var fileDescriptor_ecba292957595039 = []byte{
	// 211 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x8e, 0x4f, 0x4a, 0xc5, 0x30,
	0x10, 0xc6, 0x49, 0xdf, 0x7b, 0xad, 0x46, 0x74, 0x11, 0x37, 0xa1, 0x1b, 0x83, 0xb8, 0xe8, 0xc6,
	0x14, 0x14, 0x0f, 0xe0, 0x01, 0x44, 0x29, 0xee, 0x25, 0x6d, 0xc6, 0x12, 0x68, 0x92, 0x92, 0x8c,
	0xbd, 0x9d, 0x77, 0x93, 0x4e, 0xed, 0xe2, 0xed, 0xe6, 0xfb, 0xc3, 0x37, 0x3f, 0x7e, 0x3b, 0x4f,
	0x26, 0x00, 0xe6, 0x47, 0x1f, 0xc3, 0x18, 0xf5, 0x9c, 0x22, 0x46, 0x51, 0x65, 0x48, 0x8b, 0x1b,
	0xa0, 0xbe, 0x1b, 0x63, 0x1c, 0x27, 0x68, 0xc9, 0xee, 0x7f, 0xbe, 0x5b, 0x74, 0x1e, 0x32, 0x1a,
	0x3f, 0x6f, 0xcd, 0xfb, 0x5f, 0xc6, 0xcb, 0x0f, 0x5a, 0x10, 0x82, 0x1f, 0x83, 0xf1, 0x20, 0x99,
	0x62, 0xcd, 0x65, 0x47, 0xb7, 0xb8, 0xe1, 0x85, 0xb3, 0xb2, 0x50, 0xac, 0x39, 0x75, 0x85, 0xb3,
	0x42, 0xf2, 0x6a, 0x98, 0x9c, 0x37, 0x08, 0xf2, 0x40, 0xb5, 0x5d, 0xae, 0x09, 0x42, 0x4a, 0xc6,
	0x05, 0x79, 0x54, 0x87, 0x35, 0xf9, 0x97, 0xe2, 0x81, 0x5f, 0xd3, 0xd7, 0xf7, 0xf0, 0x16, 0x17,
	0x07, 0x59, 0x9e, 0x68, 0xee, 0xdc, 0x14, 0x2f, 0xfc, 0xc2, 0x58, 0x0b, 0xf6, 0xcb, 0xa0, 0x2c,
	0x15, 0x6b, 0xae, 0x9e, 0x6a, 0xbd, 0xc1, 0xeb, 0x1d, 0x5e, 0x7f, 0xee, 0xf0, 0x5d, 0x45, 0xdd,
	0x57, 0xec, 0x4b, 0x0a, 0x9f, 0xff, 0x02, 0x00, 0x00, 0xff, 0xff, 0xd9, 0x58, 0x77, 0x10, 0x07,
	0x01, 0x00, 0x00,
}
