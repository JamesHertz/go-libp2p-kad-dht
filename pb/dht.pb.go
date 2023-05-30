// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dht.proto

package dht_pb

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	pb "github.com/libp2p/go-libp2p-record/pb"
	_ "github.com/libp2p/go-libp2p/core/crypto/pb"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type Message_ConnectionType int32

const (
	// sender does not have a connection to peer, and no extra information (default)
	Message_NOT_CONNECTED Message_ConnectionType = 0
	// sender has a live connection to peer
	Message_CONNECTED Message_ConnectionType = 1
	// sender recently connected to peer
	Message_CAN_CONNECT Message_ConnectionType = 2
	// sender recently tried to connect to peer repeatedly but failed to connect
	// ("try" here is loose, but this should signal "made strong effort, failed")
	Message_CANNOT_CONNECT Message_ConnectionType = 3
)

var Message_ConnectionType_name = map[int32]string{
	0: "NOT_CONNECTED",
	1: "CONNECTED",
	2: "CAN_CONNECT",
	3: "CANNOT_CONNECT",
}

var Message_ConnectionType_value = map[string]int32{
	"NOT_CONNECTED":  0,
	"CONNECTED":      1,
	"CAN_CONNECT":    2,
	"CANNOT_CONNECT": 3,
}

func (x Message_ConnectionType) String() string {
	return proto.EnumName(Message_ConnectionType_name, int32(x))
}

func (Message_ConnectionType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_616a434b24c97ff4, []int{0, 0}
}

type Message struct {
	// defines what type of message it is.
	Feature string `protobuf:"bytes,1,opt,name=feature,proto3" json:"feature,omitempty"`
	// Used to specify the key associated with this message.
	// PUT_VALUE, GET_VALUE, ADD_PROVIDER, GET_PROVIDERS
	Key []byte `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	// Used to return a value
	// PUT_VALUE, GET_VALUE
	Record *pb.Record `protobuf:"bytes,3,opt,name=record,proto3" json:"record,omitempty"`
	// Used to return Providers
	// GET_VALUE, ADD_PROVIDER, GET_PROVIDERS
	ProviderPeers []Message_Peer `protobuf:"bytes,9,rep,name=providerPeers,proto3" json:"providerPeers"`
	// Used to return peers closer to a key in a query
	// GET_VALUE, GET_PROVIDERS, FIND_NODE
	CloserPeers []Message_Peer `protobuf:"bytes,8,rep,name=closerPeers,proto3" json:"closerPeers"`
	// defines what coral cluster level this query/response belongs to.
	// in case we want to implement coral's cluster rings in the future.
	ClusterLevelRaw      int32                 `protobuf:"varint,10,opt,name=clusterLevelRaw,proto3" json:"clusterLevelRaw,omitempty"`
	ProviderPeersII      []Message_PeerWithKey `protobuf:"bytes,11,rep,name=providerPeersII,proto3" json:"providerPeersII"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_616a434b24c97ff4, []int{0}
}
func (m *Message) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Message.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return m.Size()
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetFeature() string {
	if m != nil {
		return m.Feature
	}
	return ""
}

func (m *Message) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *Message) GetRecord() *pb.Record {
	if m != nil {
		return m.Record
	}
	return nil
}

func (m *Message) GetProviderPeers() []Message_Peer {
	if m != nil {
		return m.ProviderPeers
	}
	return nil
}

func (m *Message) GetCloserPeers() []Message_Peer {
	if m != nil {
		return m.CloserPeers
	}
	return nil
}

func (m *Message) GetClusterLevelRaw() int32 {
	if m != nil {
		return m.ClusterLevelRaw
	}
	return 0
}

func (m *Message) GetProviderPeersII() []Message_PeerWithKey {
	if m != nil {
		return m.ProviderPeersII
	}
	return nil
}

type Message_Peer struct {
	// ID of a given peer.
	Id byteString `protobuf:"bytes,1,opt,name=id,proto3,customtype=byteString" json:"id"`
	// multiaddrs for a given peer
	Addrs [][]byte `protobuf:"bytes,2,rep,name=addrs,proto3" json:"addrs,omitempty"`
	// used to signal the sender's connection capabilities to the peer
	Connection           Message_ConnectionType `protobuf:"varint,3,opt,name=connection,proto3,enum=dht.pb.Message_ConnectionType" json:"connection,omitempty"`
	Features             []string               `protobuf:"bytes,4,rep,name=features,proto3" json:"features,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *Message_Peer) Reset()         { *m = Message_Peer{} }
func (m *Message_Peer) String() string { return proto.CompactTextString(m) }
func (*Message_Peer) ProtoMessage()    {}
func (*Message_Peer) Descriptor() ([]byte, []int) {
	return fileDescriptor_616a434b24c97ff4, []int{0, 0}
}
func (m *Message_Peer) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Message_Peer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Message_Peer.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Message_Peer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message_Peer.Merge(m, src)
}
func (m *Message_Peer) XXX_Size() int {
	return m.Size()
}
func (m *Message_Peer) XXX_DiscardUnknown() {
	xxx_messageInfo_Message_Peer.DiscardUnknown(m)
}

var xxx_messageInfo_Message_Peer proto.InternalMessageInfo

func (m *Message_Peer) GetAddrs() [][]byte {
	if m != nil {
		return m.Addrs
	}
	return nil
}

func (m *Message_Peer) GetConnection() Message_ConnectionType {
	if m != nil {
		return m.Connection
	}
	return Message_NOT_CONNECTED
}

func (m *Message_Peer) GetFeatures() []string {
	if m != nil {
		return m.Features
	}
	return nil
}

// TODO: should I have features over here???
type Message_PeerWithKey struct {
	// ID of a given peer.
	Id byteString `protobuf:"bytes,1,opt,name=id,proto3,customtype=byteString" json:"id"`
	// multiaddrs for a given peer
	Addrs [][]byte `protobuf:"bytes,2,rep,name=addrs,proto3" json:"addrs,omitempty"`
	// used to signal the sender's connection capabilities to the peer
	Connection Message_ConnectionType `protobuf:"varint,3,opt,name=connection,proto3,enum=dht.pb.Message_ConnectionType" json:"connection,omitempty"`
	// key provided by peer; only set for GET_PROVIDER prefix lookups
	Key []byte `protobuf:"bytes,4,opt,name=key,proto3" json:"key,omitempty"`
	// signature of the provided key + encrypted peer ID for ADD_PROVIDER messages
	Signature []byte `protobuf:"bytes,5,opt,name=signature,proto3" json:"signature,omitempty"`
	// public key of the peer for ADD_PROVIDER messages
	PublicKey            []byte   `protobuf:"bytes,6,opt,name=publicKey,proto3" json:"publicKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message_PeerWithKey) Reset()         { *m = Message_PeerWithKey{} }
func (m *Message_PeerWithKey) String() string { return proto.CompactTextString(m) }
func (*Message_PeerWithKey) ProtoMessage()    {}
func (*Message_PeerWithKey) Descriptor() ([]byte, []int) {
	return fileDescriptor_616a434b24c97ff4, []int{0, 1}
}
func (m *Message_PeerWithKey) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Message_PeerWithKey) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Message_PeerWithKey.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Message_PeerWithKey) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message_PeerWithKey.Merge(m, src)
}
func (m *Message_PeerWithKey) XXX_Size() int {
	return m.Size()
}
func (m *Message_PeerWithKey) XXX_DiscardUnknown() {
	xxx_messageInfo_Message_PeerWithKey.DiscardUnknown(m)
}

var xxx_messageInfo_Message_PeerWithKey proto.InternalMessageInfo

func (m *Message_PeerWithKey) GetAddrs() [][]byte {
	if m != nil {
		return m.Addrs
	}
	return nil
}

func (m *Message_PeerWithKey) GetConnection() Message_ConnectionType {
	if m != nil {
		return m.Connection
	}
	return Message_NOT_CONNECTED
}

func (m *Message_PeerWithKey) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *Message_PeerWithKey) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *Message_PeerWithKey) GetPublicKey() []byte {
	if m != nil {
		return m.PublicKey
	}
	return nil
}

func init() {
	proto.RegisterEnum("dht.pb.Message_ConnectionType", Message_ConnectionType_name, Message_ConnectionType_value)
	proto.RegisterType((*Message)(nil), "dht.pb.Message")
	proto.RegisterType((*Message_Peer)(nil), "dht.pb.Message.Peer")
	proto.RegisterType((*Message_PeerWithKey)(nil), "dht.pb.Message.PeerWithKey")
}

func init() { proto.RegisterFile("dht.proto", fileDescriptor_616a434b24c97ff4) }

var fileDescriptor_616a434b24c97ff4 = []byte{
	// 502 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x53, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0xed, 0xc6, 0x69, 0x5a, 0x8f, 0xf3, 0xd5, 0x55, 0x0f, 0xab, 0x80, 0x52, 0x2b, 0x27, 0x73,
	0x88, 0x23, 0x85, 0x03, 0x17, 0x84, 0x68, 0x42, 0x0f, 0x55, 0x20, 0x45, 0x4b, 0xa4, 0x1e, 0x51,
	0x6c, 0x6f, 0x9d, 0x15, 0x21, 0x6b, 0xad, 0xd7, 0x45, 0xfe, 0x2f, 0xdc, 0xf9, 0x2b, 0x3d, 0x72,
	0xe1, 0xc2, 0xa1, 0x42, 0xf9, 0x25, 0xc8, 0x6b, 0x27, 0x71, 0x22, 0x04, 0x47, 0x6e, 0xef, 0xcd,
	0xbc, 0xe7, 0x37, 0x1e, 0x8f, 0xc1, 0x0c, 0x16, 0xca, 0x8d, 0xa4, 0x50, 0x02, 0xd7, 0x34, 0xf4,
	0x3a, 0x2f, 0x42, 0xae, 0x16, 0x89, 0xe7, 0xfa, 0xe2, 0xf3, 0x60, 0xc9, 0xbd, 0x68, 0x18, 0x0d,
	0x42, 0xd1, 0x2f, 0x90, 0x2f, 0x24, 0x1b, 0xf8, 0x32, 0x8d, 0x94, 0x18, 0x44, 0x5e, 0x81, 0xf2,
	0x07, 0x74, 0x86, 0x7f, 0x31, 0xf6, 0x25, 0xf3, 0x85, 0x0c, 0x32, 0x4f, 0x8e, 0x0a, 0x4f, 0xbf,
	0xe4, 0x09, 0x45, 0x28, 0x06, 0xba, 0xec, 0x25, 0x77, 0x9a, 0x69, 0xa2, 0x51, 0x2e, 0xef, 0x7d,
	0xab, 0xc1, 0xc9, 0x3b, 0x16, 0xc7, 0xf3, 0x90, 0x61, 0x02, 0x27, 0x77, 0x6c, 0xae, 0x12, 0xc9,
	0x08, 0xb2, 0x91, 0x63, 0xd2, 0x0d, 0xc5, 0x6d, 0x30, 0x3e, 0xb1, 0x94, 0x54, 0x6c, 0xe4, 0xd4,
	0x69, 0x06, 0xf1, 0x33, 0xa8, 0xe5, 0xb1, 0xc4, 0xb0, 0x91, 0x63, 0x0d, 0xcf, 0xdc, 0xcd, 0x14,
	0x9e, 0x4b, 0x35, 0xa2, 0x85, 0x00, 0xbf, 0x86, 0x46, 0x24, 0xc5, 0x3d, 0x0f, 0x98, 0x7c, 0xcf,
	0x98, 0x8c, 0x89, 0x69, 0x1b, 0x8e, 0x35, 0x3c, 0x77, 0xf3, 0xf5, 0xb8, 0x45, 0xbc, 0x9b, 0x35,
	0x47, 0xd5, 0x87, 0xc7, 0x8b, 0x23, 0xba, 0x6f, 0xc0, 0x2f, 0xc1, 0xf2, 0x97, 0x22, 0xde, 0xf8,
	0x4f, 0xff, 0xe9, 0x2f, 0xcb, 0xb1, 0x03, 0x2d, 0x7f, 0x99, 0xc4, 0x8a, 0xc9, 0xb7, 0xec, 0x9e,
	0x2d, 0xe9, 0xfc, 0x0b, 0x01, 0x1b, 0x39, 0xc7, 0xf4, 0xb0, 0x8c, 0x27, 0xd0, 0xda, 0x0b, 0xbe,
	0xbe, 0x26, 0x96, 0xce, 0x7a, 0xf2, 0xa7, 0xac, 0x5b, 0xae, 0x16, 0x13, 0x96, 0x16, 0x91, 0x87,
	0xce, 0xce, 0x57, 0x04, 0xd5, 0x0c, 0xe3, 0x1e, 0x54, 0x78, 0xa0, 0x37, 0x5a, 0x1f, 0xe1, 0x4c,
	0xfb, 0xf3, 0xf1, 0x02, 0xbc, 0x54, 0xb1, 0x0f, 0x4a, 0xf2, 0x55, 0x48, 0x2b, 0x3c, 0xc0, 0xe7,
	0x70, 0x3c, 0x0f, 0x02, 0x19, 0x93, 0x8a, 0x6d, 0x38, 0x75, 0x9a, 0x13, 0xfc, 0x0a, 0xc0, 0x17,
	0xab, 0x15, 0xf3, 0x15, 0x17, 0x2b, 0xbd, 0xe8, 0xe6, 0xb0, 0x7b, 0x38, 0xca, 0x78, 0xab, 0x98,
	0xa5, 0x11, 0xa3, 0x25, 0x07, 0xee, 0xc0, 0x69, 0xf1, 0x05, 0x63, 0x52, 0xb5, 0x0d, 0xc7, 0xa4,
	0x5b, 0xde, 0xf9, 0x81, 0xc0, 0x2a, 0xbd, 0xc5, 0x7f, 0x9c, 0xb2, 0x38, 0xae, 0xea, 0xee, 0xb8,
	0x9e, 0x82, 0x19, 0xf3, 0x70, 0x95, 0x9f, 0xe2, 0xb1, 0xae, 0xef, 0x0a, 0x59, 0x37, 0x4a, 0xbc,
	0x25, 0xf7, 0x27, 0x2c, 0x25, 0xb5, 0xbc, 0xbb, 0x2d, 0xf4, 0x6e, 0xa1, 0xb9, 0x9f, 0x85, 0xcf,
	0xa0, 0x31, 0xbd, 0x99, 0x7d, 0x1c, 0xdf, 0x4c, 0xa7, 0x57, 0xe3, 0xd9, 0xd5, 0x9b, 0xf6, 0x11,
	0x6e, 0x80, 0xb9, 0xa3, 0x08, 0xb7, 0xc0, 0x1a, 0x5f, 0x4e, 0x37, 0x8a, 0x76, 0x05, 0x63, 0x68,
	0x8e, 0x2f, 0xa7, 0x25, 0x57, 0xdb, 0x18, 0xd5, 0x1f, 0xd6, 0x5d, 0xf4, 0x7d, 0xdd, 0x45, 0xbf,
	0xd6, 0x5d, 0xe4, 0xd5, 0xf4, 0xef, 0xf3, 0xfc, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0x95, 0x36,
	0xf4, 0xb5, 0xef, 0x03, 0x00, 0x00,
}

func (m *Message) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Message) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Message) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.ProviderPeersII) > 0 {
		for iNdEx := len(m.ProviderPeersII) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ProviderPeersII[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintDht(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x5a
		}
	}
	if m.ClusterLevelRaw != 0 {
		i = encodeVarintDht(dAtA, i, uint64(m.ClusterLevelRaw))
		i--
		dAtA[i] = 0x50
	}
	if len(m.ProviderPeers) > 0 {
		for iNdEx := len(m.ProviderPeers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ProviderPeers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintDht(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x4a
		}
	}
	if len(m.CloserPeers) > 0 {
		for iNdEx := len(m.CloserPeers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.CloserPeers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintDht(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x42
		}
	}
	if m.Record != nil {
		{
			size, err := m.Record.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintDht(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Key) > 0 {
		i -= len(m.Key)
		copy(dAtA[i:], m.Key)
		i = encodeVarintDht(dAtA, i, uint64(len(m.Key)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Feature) > 0 {
		i -= len(m.Feature)
		copy(dAtA[i:], m.Feature)
		i = encodeVarintDht(dAtA, i, uint64(len(m.Feature)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Message_Peer) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Message_Peer) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Message_Peer) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Features) > 0 {
		for iNdEx := len(m.Features) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Features[iNdEx])
			copy(dAtA[i:], m.Features[iNdEx])
			i = encodeVarintDht(dAtA, i, uint64(len(m.Features[iNdEx])))
			i--
			dAtA[i] = 0x22
		}
	}
	if m.Connection != 0 {
		i = encodeVarintDht(dAtA, i, uint64(m.Connection))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Addrs) > 0 {
		for iNdEx := len(m.Addrs) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Addrs[iNdEx])
			copy(dAtA[i:], m.Addrs[iNdEx])
			i = encodeVarintDht(dAtA, i, uint64(len(m.Addrs[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size := m.Id.Size()
		i -= size
		if _, err := m.Id.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintDht(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *Message_PeerWithKey) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Message_PeerWithKey) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Message_PeerWithKey) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.PublicKey) > 0 {
		i -= len(m.PublicKey)
		copy(dAtA[i:], m.PublicKey)
		i = encodeVarintDht(dAtA, i, uint64(len(m.PublicKey)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.Signature) > 0 {
		i -= len(m.Signature)
		copy(dAtA[i:], m.Signature)
		i = encodeVarintDht(dAtA, i, uint64(len(m.Signature)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Key) > 0 {
		i -= len(m.Key)
		copy(dAtA[i:], m.Key)
		i = encodeVarintDht(dAtA, i, uint64(len(m.Key)))
		i--
		dAtA[i] = 0x22
	}
	if m.Connection != 0 {
		i = encodeVarintDht(dAtA, i, uint64(m.Connection))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Addrs) > 0 {
		for iNdEx := len(m.Addrs) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Addrs[iNdEx])
			copy(dAtA[i:], m.Addrs[iNdEx])
			i = encodeVarintDht(dAtA, i, uint64(len(m.Addrs[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size := m.Id.Size()
		i -= size
		if _, err := m.Id.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintDht(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintDht(dAtA []byte, offset int, v uint64) int {
	offset -= sovDht(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Message) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Feature)
	if l > 0 {
		n += 1 + l + sovDht(uint64(l))
	}
	l = len(m.Key)
	if l > 0 {
		n += 1 + l + sovDht(uint64(l))
	}
	if m.Record != nil {
		l = m.Record.Size()
		n += 1 + l + sovDht(uint64(l))
	}
	if len(m.CloserPeers) > 0 {
		for _, e := range m.CloserPeers {
			l = e.Size()
			n += 1 + l + sovDht(uint64(l))
		}
	}
	if len(m.ProviderPeers) > 0 {
		for _, e := range m.ProviderPeers {
			l = e.Size()
			n += 1 + l + sovDht(uint64(l))
		}
	}
	if m.ClusterLevelRaw != 0 {
		n += 1 + sovDht(uint64(m.ClusterLevelRaw))
	}
	if len(m.ProviderPeersII) > 0 {
		for _, e := range m.ProviderPeersII {
			l = e.Size()
			n += 1 + l + sovDht(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Message_Peer) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Id.Size()
	n += 1 + l + sovDht(uint64(l))
	if len(m.Addrs) > 0 {
		for _, b := range m.Addrs {
			l = len(b)
			n += 1 + l + sovDht(uint64(l))
		}
	}
	if m.Connection != 0 {
		n += 1 + sovDht(uint64(m.Connection))
	}
	if len(m.Features) > 0 {
		for _, s := range m.Features {
			l = len(s)
			n += 1 + l + sovDht(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Message_PeerWithKey) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Id.Size()
	n += 1 + l + sovDht(uint64(l))
	if len(m.Addrs) > 0 {
		for _, b := range m.Addrs {
			l = len(b)
			n += 1 + l + sovDht(uint64(l))
		}
	}
	if m.Connection != 0 {
		n += 1 + sovDht(uint64(m.Connection))
	}
	l = len(m.Key)
	if l > 0 {
		n += 1 + l + sovDht(uint64(l))
	}
	l = len(m.Signature)
	if l > 0 {
		n += 1 + l + sovDht(uint64(l))
	}
	l = len(m.PublicKey)
	if l > 0 {
		n += 1 + l + sovDht(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovDht(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozDht(x uint64) (n int) {
	return sovDht(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Message) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDht
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
			return fmt.Errorf("proto: Message: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Message: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Feature", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDht
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
				return ErrInvalidLengthDht
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDht
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Feature = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDht
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
				return ErrInvalidLengthDht
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthDht
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Key = append(m.Key[:0], dAtA[iNdEx:postIndex]...)
			if m.Key == nil {
				m.Key = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Record", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDht
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
				return ErrInvalidLengthDht
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthDht
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Record == nil {
				m.Record = &pb.Record{}
			}
			if err := m.Record.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CloserPeers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDht
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
				return ErrInvalidLengthDht
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthDht
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CloserPeers = append(m.CloserPeers, Message_Peer{})
			if err := m.CloserPeers[len(m.CloserPeers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProviderPeers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDht
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
				return ErrInvalidLengthDht
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthDht
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ProviderPeers = append(m.ProviderPeers, Message_Peer{})
			if err := m.ProviderPeers[len(m.ProviderPeers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClusterLevelRaw", wireType)
			}
			m.ClusterLevelRaw = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDht
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ClusterLevelRaw |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProviderPeersII", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDht
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
				return ErrInvalidLengthDht
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthDht
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ProviderPeersII = append(m.ProviderPeersII, Message_PeerWithKey{})
			if err := m.ProviderPeersII[len(m.ProviderPeersII)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDht(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDht
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Message_Peer) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDht
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
			return fmt.Errorf("proto: Peer: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Peer: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDht
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
				return ErrInvalidLengthDht
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthDht
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Id.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Addrs", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDht
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
				return ErrInvalidLengthDht
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthDht
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Addrs = append(m.Addrs, make([]byte, postIndex-iNdEx))
			copy(m.Addrs[len(m.Addrs)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Connection", wireType)
			}
			m.Connection = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDht
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Connection |= Message_ConnectionType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Features", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDht
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
				return ErrInvalidLengthDht
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDht
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Features = append(m.Features, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDht(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDht
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Message_PeerWithKey) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDht
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
			return fmt.Errorf("proto: PeerWithKey: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PeerWithKey: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDht
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
				return ErrInvalidLengthDht
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthDht
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Id.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Addrs", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDht
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
				return ErrInvalidLengthDht
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthDht
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Addrs = append(m.Addrs, make([]byte, postIndex-iNdEx))
			copy(m.Addrs[len(m.Addrs)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Connection", wireType)
			}
			m.Connection = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDht
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Connection |= Message_ConnectionType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDht
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
				return ErrInvalidLengthDht
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthDht
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Key = append(m.Key[:0], dAtA[iNdEx:postIndex]...)
			if m.Key == nil {
				m.Key = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signature", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDht
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
				return ErrInvalidLengthDht
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthDht
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signature = append(m.Signature[:0], dAtA[iNdEx:postIndex]...)
			if m.Signature == nil {
				m.Signature = []byte{}
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PublicKey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDht
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
				return ErrInvalidLengthDht
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthDht
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PublicKey = append(m.PublicKey[:0], dAtA[iNdEx:postIndex]...)
			if m.PublicKey == nil {
				m.PublicKey = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDht(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDht
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipDht(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDht
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
					return 0, ErrIntOverflowDht
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowDht
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
				return 0, ErrInvalidLengthDht
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupDht
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthDht
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthDht        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDht          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupDht = fmt.Errorf("proto: unexpected end of group")
)
