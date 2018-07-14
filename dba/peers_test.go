package dba_test

import (
	. "github.com/satellitex/bbft/dba"
	"github.com/satellitex/bbft/model"
	. "github.com/satellitex/bbft/test_utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testPeerService(t *testing.T, p PeerService) {
	peers := []model.Peer{
		RandomPeer(),
		RandomPeer(),
		RandomPeer(),
		RandomPeer(),
	}
	for cnt, peer := range peers {
		assert.Equal(t, cnt, p.Size())

		_, ok := p.GetPeer(peer.GetPubkey())
		assert.False(t, ok)

		_, ok = p.GetPeerFromAddress(peer.GetAddress())
		assert.False(t, ok)

		p.AddPeer(peer)

		assert.Equal(t, cnt+1, p.Size())

		actPeer, ok := p.GetPeer(peer.GetPubkey())
		assert.True(t, ok)
		assert.Equal(t, peer, actPeer)

		actPeer, ok = p.GetPeerFromAddress(peer.GetAddress())
		assert.True(t, ok)
		assert.Equal(t, peer, actPeer)
	}
}

func TestPeerServiceOnMemory(t *testing.T) {
	peerService := NewPeerServiceOnMemory()
	testPeerService(t, peerService)
}
