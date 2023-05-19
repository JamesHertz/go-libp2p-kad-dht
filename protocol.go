package dht

import (
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/core/peer"
	pb "github.com/libp2p/go-libp2p-kad-dht/pb"
)

var (
	// ProtocolDHT is the default DHT protocol.
	ProtocolDHT protocol.ID = "/ipfs/kad/1.0.0"
	// DefaultProtocols spoken by the DHT.
	DefaultProtocols = []protocol.ID{ProtocolDHT}
	Features = peer.FeatureList{
		pb.IPFS_GET_PROVIDERS, pb.IPFS_ADD_PROVIDERS, pb.IPFS_GET_VALUE, 
	    pb.IPFS_PUT_VALUE, pb.IPFS_PING, pb.FIND_CLOSEST_PEERS,
	}
)
