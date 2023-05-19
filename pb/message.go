package dht_pb

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	logging "github.com/ipfs/go-log"
	ma "github.com/multiformats/go-multiaddr"
)

var log = logging.Logger("dht.pb")

type PeerRoutingInfo struct {
	peer.AddrInfo
	network.Connectedness
}

// TODO: CHANGE FEATURES LIST THING :)

// type marshableMsg interface {
// 	XXX_Unmarshal([]byte) error
// 	XXX_Marshal([]byte, bool) ([]byte, error)
// 	XXX_Size() int
// }

// TODO: think about deduplicating this thing
func (msg *Message) GetIpfsMsg() (*IpfsMessage, error) {
	res := &IpfsMessage{}
	if err := res.Unmarshal(msg.Data); err != nil {
		return nil, fmt.Errorf("Invalid IpfsMsg: %v")
	}
	return res, nil
}

func (msg *Message) GetBareMsg() (*Message_BareMsg, error) {
	res := &Message_BareMsg{}
	res.Marshal()
	if err := res.Unmarshal(msg.Data); err != nil {
		return nil, fmt.Errorf("Invalid BareMsg: %v")
	}
	return res, nil
}

func (msg *Message) GetMsgFeature() peer.Feature {
	return peer.Feature(msg.Feature)
}

func ToDhtMessage[
	T interface {
		Marshal() ([]byte, error)
		Size() int
	}](msg T, feature peer.Feature) *Message {
	data, err := msg.Marshal()
	if err != nil {
		panic("Error marshalling something... You have to fix this...")
	}
	return &Message{
		Feature: string(feature),
		Data:    data,
	}
}

// NewMessage constructs a new dht message with given type, key, and level
func NewIpfsMsg(key []byte, level int) *IpfsMessage {
	m := &IpfsMessage{
		Key: key,
	}
	m.SetClusterLevel(level)
	return m
}

func peerRoutingInfoToPBPeer(p PeerRoutingInfo) *Peer {
	var pbp Peer

	pbp.Addrs = make([][]byte, len(p.Addrs))
	for i, maddr := range p.Addrs {
		pbp.Addrs[i] = maddr.Bytes() // Bytes, not String. Compressed.
	}
	pbp.Id = byteString(p.ID)
	pbp.Connection = ConnectionType(p.Connectedness)
	return &pbp
}

func peerInfoToPBPeer(p peer.AddrInfo) Peer {
	var pbp Peer

	pbp.Addrs = make([][]byte, len(p.Addrs))
	for i, maddr := range p.Addrs {
		pbp.Addrs[i] = maddr.Bytes() // Bytes, not String. Compressed.
	}
	pbp.Id = byteString(p.ID)
	pbp.Features = nil // TODO: change this :)
	return pbp
}

// PBPeerToPeer turns a *Message_Peer into its peer.AddrInfo counterpart
func PBPeerToPeerInfo(pbp Peer) peer.AddrInfo {
	return peer.AddrInfo{
		ID:    peer.ID(pbp.Id),
		Addrs: pbp.Addresses(),
	}
}

// RawPeerInfosToPBPeers converts a slice of Peers into a slice of *Message_Peers,
// ready to go out on the wire.
func RawPeerInfosToPBPeers(peers []peer.AddrInfo) []Peer {
	pbpeers := make([]Peer, len(peers))
	for i, p := range peers {
		pbpeers[i] = peerInfoToPBPeer(p)
	}
	return pbpeers
}

// PeersToPBPeers converts given []peer.Peer into a set of []*Message_Peer,
// which can be written to a message and sent out. the key thing this function
// does (in addition to PeersToPBPeers) is set the ConnectionType with
// information from the given network.Network.
func PeerInfosToPBPeers(n network.Network, peers []peer.AddrInfo) []Peer {
	pbps := RawPeerInfosToPBPeers(peers)
	for i, pbp := range pbps {
		c := ConnectionType(n.Connectedness(peers[i].ID))
		pbp.Connection = c
	}
	return pbps
}

func PeerRoutingInfosToPBPeers(peers []PeerRoutingInfo) []*Peer {
	pbpeers := make([]*Peer, len(peers))
	for i, p := range peers {
		pbpeers[i] = peerRoutingInfoToPBPeer(p)
	}
	return pbpeers
}

// PBPeersToPeerInfos converts given []*Message_Peer into []peer.AddrInfo
// Invalid addresses will be silently omitted.
func PBPeersToPeerInfos(pbps []Peer) []*peer.AddrInfo {
	peers := make([]*peer.AddrInfo, 0, len(pbps))
	for _, pbp := range pbps {
		ai := PBPeerToPeerInfo(pbp)
		peers = append(peers, &ai)
	}
	return peers
}

// Addresses returns a multiaddr associated with the Message_Peer entry
func (m *Peer) Addresses() []ma.Multiaddr {
	if m == nil {
		return nil
	}

	maddrs := make([]ma.Multiaddr, 0, len(m.Addrs))
	for _, addr := range m.Addrs {
		maddr, err := ma.NewMultiaddrBytes(addr)
		if err != nil {
			log.Debugw("error decoding multiaddr for peer", "peer", peer.ID(m.Id), "error", err)
			continue
		}

		maddrs = append(maddrs, maddr)
	}
	return maddrs
}

// GetClusterLevel gets and adjusts the cluster level on the message.
// a +/- 1 adjustment is needed to distinguish a valid first level (1) and
// default "no value" protobuf behavior (0)
func (m *IpfsMessage) GetClusterLevel() int {
	level := m.GetClusterLevelRaw() - 1
	if level < 0 {
		return 0
	}
	return int(level)
}

// SetClusterLevel adjusts and sets the cluster level on the message.
// a +/- 1 adjustment is needed to distinguish a valid first level (1) and
// default "no value" protobuf behavior (0)
func (m *IpfsMessage) SetClusterLevel(level int) {
	lvl := int32(level)
	m.ClusterLevelRaw = lvl + 1
}

// ConnectionType returns a Message_ConnectionType associated with the
// network.Connectedness.
func ConnectionType(c network.Connectedness) Peer_ConnectionType {
	switch c {
	default:
		return Peer_NOT_CONNECTED //Message_NOT_CONNECTED
	case network.NotConnected:
		return Peer_NOT_CONNECTED // Message_NOT_CONNECTED
	case network.Connected:
		return Peer_CONNECTED //Message_CONNECTED
	case network.CanConnect:
		return Peer_CAN_CONNECT //Message_CAN_CONNECT
	case network.CannotConnect:
		return Peer_CANNOT_CONNECT //Message_CANNOT_CONNECT
	}
}

// Connectedness returns an network.Connectedness associated with the
// Message_ConnectionType.
func Connectedness(c Peer_ConnectionType) network.Connectedness {
	switch c {
	default:
		return network.NotConnected
	case Peer_NOT_CONNECTED: //Message_NOT_CONNECTED:
		return network.NotConnected
	case Peer_CONNECTED: //Message_CONNECTED:
		return network.Connected
	case Peer_CAN_CONNECT: //Message_CAN_CONNECT:
		return network.CanConnect
	case Peer_CANNOT_CONNECT: //Message_CANNOT_CONNECT:
		return network.CannotConnect
	}
}
