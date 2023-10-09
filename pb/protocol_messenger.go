package dht_pb

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	logging "github.com/ipfs/go-log"
	recpb "github.com/libp2p/go-libp2p-record/pb"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multihash"

	"github.com/libp2p/go-libp2p/core/crypto"

	"github.com/libp2p/go-libp2p-kad-dht/internal"
)

var logger = logging.Logger("dht")


// features (at least by now c:)
const (
	// new features :)
	IPFS_DH_GET_PROVIDERS peer.Feature = "/ipfs/dh/getproviders"
	IPFS_DH_ADD_PROVIDERS peer.Feature = "/ipfs/dh/putproviders"

	// old features
	IPFS_GET_PROVIDERS peer.Feature = "/ipfs/getproviders"
	IPFS_ADD_PROVIDERS peer.Feature = "/ipfs/putproviders"

	// standard features
	IPFS_GET_VALUE     peer.Feature = "/ipfs/get"
	IPFS_PUT_VALUE     peer.Feature = "/ipfs/put"
	IPFS_PING          peer.Feature = "/ipfs/ping"
	FIND_CLOSEST_PEERS peer.Feature = "/libp2p/findclosestpeer"
)

// ProtocolMessenger can be used for sending DHT messages to peers and processing their responses.
// This decouples the wire protocol format from both the DHT protocol implementation and from the implementation of the
// routing.Routing interface.
//
// Note: the ProtocolMessenger's MessageSender still needs to deal with some wire protocol details such as using
// varint-delineated protobufs
type ProtocolMessenger struct {
	m MessageSender
}

type ProtocolMessengerOption func(*ProtocolMessenger) error

// NewProtocolMessenger creates a new ProtocolMessenger that is used for sending DHT messages to peers and processing
// their responses.
func NewProtocolMessenger(msgSender MessageSender, opts ...ProtocolMessengerOption) (*ProtocolMessenger, error) {
	pm := &ProtocolMessenger{
		m: msgSender,
	}

	for _, o := range opts {
		if err := o(pm); err != nil {
			return nil, err
		}
	}

	return pm, nil
}

// MessageSender handles sending wire protocol messages to a given peer
type MessageSender interface {
	// SendRequest sends a peer a message and waits for its response
	SendRequest(ctx context.Context, p peer.ID, pmes *Message) (*Message, error)
	// SendMessage sends a peer a message without waiting on a response
	SendMessage(ctx context.Context, p peer.ID, pmes *Message) error
}

// PutValue asks a peer to store the given key/value pair.
func (pm *ProtocolMessenger) PutValue(ctx context.Context, p peer.ID, rec *recpb.Record) error {
	pmes := NewMessage(IPFS_PUT_VALUE, rec.Key, 0)
	pmes.Record = rec
	rpmes, err := pm.m.SendRequest(ctx, p, pmes)
	if err != nil {
		logger.Debugw("failed to put value to peer", "to", p, "key", internal.LoggableRecordKeyBytes(rec.Key), "error", err)
		return err
	}

	if !bytes.Equal(rpmes.GetRecord().Value, pmes.GetRecord().Value) {
		const errStr = "value not put correctly"
		logger.Infow(errStr, "put-message", pmes, "get-message", rpmes)
		return errors.New(errStr)
	}

	return nil
}

// GetValue asks a peer for the value corresponding to the given key. Also returns the K closest peers to the key
// as described in GetClosestPeers.
func (pm *ProtocolMessenger) GetValue(ctx context.Context, p peer.ID, key string) (*recpb.Record, []*peer.AddrInfo, error) {
	pmes := NewMessage(IPFS_GET_VALUE, []byte(key), 0)
	respMsg, err := pm.m.SendRequest(ctx, p, pmes)
	if err != nil {
		return nil, nil, err
	}

	// Perhaps we were given closer peers
	peers := PBPeersToPeerInfos(respMsg.GetCloserPeers())

	if rec := respMsg.GetRecord(); rec != nil {
		// Success! We were given the value
		logger.Debug("got value")

		// Check that record matches the one we are looking for (validation of the record does not happen here)
		if !bytes.Equal([]byte(key), rec.GetKey()) {
			logger.Debug("received incorrect record")
			return nil, nil, internal.ErrIncorrectRecord
		}

		return rec, peers, err
	}

	return nil, peers, nil
}

// GetClosestPeers asks a peer to return the K (a DHT-wide parameter) DHT server peers closest in XOR space to the id
// Note: If the peer happens to know another peer whose peerID exactly matches the given id it will return that peer
// even if that peer is not a DHT server node.
func (pm *ProtocolMessenger) GetClosestPeers(ctx context.Context, p peer.ID, id peer.ID) ([]*peer.AddrInfo, error) {
	pmes := NewMessage(FIND_CLOSEST_PEERS, []byte(id), 0)
	respMsg, err := pm.m.SendRequest(ctx, p, pmes)
	if err != nil {
		return nil, err
	}
	peers := PBPeersToPeerInfos(respMsg.GetCloserPeers())
	return peers, nil
}

// PutProvider asks a peer to store that we are a provider for the given key.
// func (pm *ProtocolMessenger) PutProvider(ctx context.Context, p peer.ID, key multihash.Multihash, host host.Host) error { -removed
func (pm *ProtocolMessenger) PutProvider(ctx context.Context, p peer.ID, key multihash.Multihash, host host.Host, encID []byte) error {
// +added
	pi := peer.AddrInfo{
   		ID:    peer.ID(encID),
		Addrs: host.Addrs(),
	}	
// +added

/* -removed
	pi := peer.AddrInfo{
		ID:    host.ID(),
		Addrs: host.Addrs(),
	}
*/

	// TODO: We may want to limit the type of addresses in our provider records
	// For example, in a WAN-only DHT prohibit sharing non-WAN addresses (e.g. 192.168.0.100)
	if len(pi.Addrs) < 1 {
		return fmt.Errorf("no known addresses for self, cannot put provider")
	}

/* -removed
	pmes := NewMessage(IPFS_ADD_PROVIDERS, key, 0)
	pmes.ProviderPeers = RawPeerInfosToPBPeers([]peer.AddrInfo{pi})
*/

// +added
	// sign ( key || encID )
	privKey := host.Peerstore().PrivKey(host.ID())
	sig, err := privKey.Sign(append(key, encID...))
	if err != nil {
		return err
	}

	pubKey := host.Peerstore().PubKey(host.ID())
	// TODO: an method for pbPubKey.Marshall()
	pbPubKey, err := crypto.MarshalPublicKey(pubKey)
	// pbPubKey, err := crypto.PublicKeyToProto(pubKey)
	// bytesPubKey, err := crypto.MarshalPublicKey(pubKey)
	if err != nil {
		return err
	}

	pmes := NewMessage(IPFS_DH_ADD_PROVIDERS, key, 0)
	pmes.ProviderPeersII = PeersToPeersWithKey(RawPeerInfosToPBPeers([]peer.AddrInfo{pi}))
	pmes.ProviderPeersII[0].Signature = sig
	pmes.ProviderPeersII[0].PublicKey = pbPubKey
// +added

	return pm.m.SendMessage(ctx, p, pmes)
}

// +added
// GetProviders asks a peer for the providers it knows of for a given key. Also returns the K closest peers to the key
// as described in GetClosestPeers.
func (pm *ProtocolMessenger) GetProviders(
		ctx context.Context,
		p peer.ID,
		key []byte,
	) ([]*peer.AddrInfo, []*peer.AddrInfo, error) {
	pmes := NewMessage(IPFS_DH_GET_PROVIDERS, key, 0)
	resp, err := pm.m.SendRequest(ctx, p, pmes)
	if err != nil {
		return nil, nil, err
	}
	provs := PBPeersToAddrInfos(resp.GetProviderPeersII())
	closerPeers := PBPeersToPeerInfos(resp.GetCloserPeers())
	return provs, closerPeers, nil
}

func (pm *ProtocolMessenger) GetProvidersDefault(ctx context.Context, p peer.ID, key multihash.Multihash) ([]*peer.AddrInfo, []*peer.AddrInfo, error) { 

	pmes := NewMessage(IPFS_GET_PROVIDERS, key, 0)
	respMsg, err := pm.m.SendRequest(ctx, p, pmes)
	if err != nil {
		return nil, nil, err
	}
	provs := PBPeersToPeerInfos(respMsg.GetProviderPeers())
	closerPeers := PBPeersToPeerInfos(respMsg.GetCloserPeers())
	return provs, closerPeers, nil
}

// Ping sends a ping message to the passed peer and waits for a response.
func (pm *ProtocolMessenger) Ping(ctx context.Context, p peer.ID) error {
	req := NewMessage(IPFS_PING, nil, 0)
	resp, err := pm.m.SendRequest(ctx, p, req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	if resp.Feature!= string(IPFS_PING){
		return fmt.Errorf("got unexpected response type: %v", resp.Feature)
	}
	return nil
}


// NEW METHODS :)

// GetProvidersByPrefix asks a peer for the providers it knows of for a given key prefix.
// The returned providers list is a list of providers that have the given fullKey.
// also returns K closest peers.
func (pm *ProtocolMessenger) GetProvidersByPrefix(
	ctx context.Context,
	p peer.ID,
	lookupKey []byte,
	fullKey []byte,
) ([]*peer.AddrInfo, []*peer.AddrInfo, error) {
	pmes := NewMessage(IPFS_DH_GET_PROVIDERS, lookupKey, 0) //NewGetProvidersMessage(IPFS_DH_GET_PROVIDERS, lookupKey, 0)
	resp, err := pm.m.SendRequest(ctx, p, pmes)
	if err != nil {
		return nil, nil, err
	}

	logger.Infof("GetProvidersByPrefix received %d provs", len(resp.GetProviderPeers()))

	// if this is a prefix lookup, the providers might not actually have
	// the content we're looking for. discard all that don't
	provs := PBPeersWithKeyToAddrInfos(resp.GetProviderPeersII(), fullKey)
	closerPeers := PBPeersToPeerInfos(resp.GetCloserPeers())
	return provs, closerPeers, nil
}