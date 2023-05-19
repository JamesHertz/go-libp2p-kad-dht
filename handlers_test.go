package dht

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	proto "github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p"
	pb "github.com/libp2p/go-libp2p-kad-dht/pb"
	recpb "github.com/libp2p/go-libp2p-record/pb"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	peer "github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

func TestCleanRecordSigned(t *testing.T) {
	actual := new(recpb.Record)
	actual.TimeReceived = "time"
	actual.Value = []byte("value")
	actual.Key = []byte("key")

	cleanRecord(actual)
	actualBytes, err := proto.Marshal(actual)
	if err != nil {
		t.Fatal(err)
	}

	expected := new(recpb.Record)
	expected.Value = []byte("value")
	expected.Key = []byte("key")
	expectedBytes, err := proto.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(actualBytes, expectedBytes) {
		t.Error("failed to clean record")
	}
}

func TestCleanRecord(t *testing.T) {
	actual := new(recpb.Record)
	actual.TimeReceived = "time"
	actual.Key = []byte("key")
	actual.Value = []byte("value")

	cleanRecord(actual)
	actualBytes, err := proto.Marshal(actual)
	if err != nil {
		t.Fatal(err)
	}

	expected := new(recpb.Record)
	expected.Key = []byte("key")
	expected.Value = []byte("value")
	expectedBytes, err := proto.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(actualBytes, expectedBytes) {
		t.Error("failed to clean record")
	}
}

func TestBadMessage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dht := setupDHT(ctx, t, false)

	features := []peer.Feature{
		pb.IPFS_PUT_VALUE, pb.IPFS_GET_VALUE,
		pb.IPFS_ADD_PROVIDERS, pb.IPFS_GET_PROVIDERS,
		pb.FIND_CLOSEST_PEERS,
	}

	for _, ft := range features {
		msg := &pb.Message{
			Feature: string(ft),
		}
		_, err := dht.handlerForMsgType(ft)(ctx, dht.Host().ID(), msg)
		if err == nil {
			t.Fatalf("expected processing message to fail for type %s", pb.FIND_CLOSEST_PEERS) //Message_FIND_NODE)
		}
	}
}

func BenchmarkHandleFindPeer(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	h, err := libp2p.New()
	if err != nil {
		b.Fatal(err)
	}
	defer h.Close()

	d, err := New(ctx, h)
	if err != nil {
		b.Fatal(err)
	}

	rng := rand.New(rand.NewSource(150))
	var peers []peer.ID
	for i := 0; i < 1000; i++ {
		_, pubk, _ := crypto.GenerateEd25519Key(rng)
		id, err := peer.IDFromPublicKey(pubk)
		if err != nil {
			panic(err)
		}

		d.peerFound(ctx, id, true)

		peers = append(peers, id)
		a, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", 2000+i))
		if err != nil {
			panic(err)
		}

		d.host.Peerstore().AddAddr(id, a, time.Minute*50)
	}

	var reqs []*pb.Message
	for i := 0; i < b.N; i++ {
		reqs = append(reqs, &pb.Message{
			Key: []byte("asdasdasd"),
		})
	}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = d.handleFindPeer(ctx, peers[0], reqs[i])
		if err != nil {
			b.Error(err)
		}
	}

}
