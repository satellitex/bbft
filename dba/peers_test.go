package dba_test

import (
	. "github.com/satellitex/bbft/dba"
	"github.com/satellitex/bbft/model"
	. "github.com/satellitex/bbft/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func testPeerService(t *testing.T, p PeerService) {
	peers := []model.Peer{
		RandomPeer(),
		RandomPeer(),
		RandomPeer(),
		RandomPeer(),
	}

	t.Run("empty peerService, test", func(t *testing.T) {
		peers := p.GetPeers()
		assert.Empty(t, peers)

		peers = p.GetPermutationPeers(0)
		assert.Empty(t, peers)
	})

	t.Run("test Add And Get", func(t *testing.T) {
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
	})

	t.Run("Get Numer", func(t *testing.T) {
		// 4 peers (Minimum peers is 4)
		assert.Equal(t, 3, p.GetNumberOfRequiredAcceptPeers())
		assert.Equal(t, 1, p.GetNumberOfAllowedFailedPeers())

		// 5 peers
		p.AddPeer(RandomPeer())
		assert.Equal(t, 3, p.GetNumberOfRequiredAcceptPeers())
		assert.Equal(t, 1, p.GetNumberOfAllowedFailedPeers())

		// 6 peers
		p.AddPeer(RandomPeer())
		assert.Equal(t, 3, p.GetNumberOfRequiredAcceptPeers())
		assert.Equal(t, 1, p.GetNumberOfAllowedFailedPeers())

		// 7 peers
		p.AddPeer(RandomPeer())
		assert.Equal(t, 5, p.GetNumberOfRequiredAcceptPeers())
		assert.Equal(t, 2, p.GetNumberOfAllowedFailedPeers())

		// 8 peers
		p.AddPeer(RandomPeer())
		assert.Equal(t, 5, p.GetNumberOfRequiredAcceptPeers())
		assert.Equal(t, 2, p.GetNumberOfAllowedFailedPeers())

		// 9 peers
		p.AddPeer(RandomPeer())
		assert.Equal(t, 5, p.GetNumberOfRequiredAcceptPeers())
		assert.Equal(t, 2, p.GetNumberOfAllowedFailedPeers())

		// 10 peers
		p.AddPeer(RandomPeer())
		assert.Equal(t, 7, p.GetNumberOfRequiredAcceptPeers())
		assert.Equal(t, 3, p.GetNumberOfAllowedFailedPeers())

		// 11 peers
		p.AddPeer(RandomPeer())
		assert.Equal(t, 7, p.GetNumberOfRequiredAcceptPeers())
		assert.Equal(t, 3, p.GetNumberOfAllowedFailedPeers())
	})

	t.Run("Get Permutation", func(t *testing.T) {
		peers := p.GetPeers()
		peers2 := p.GetPeers()
		assert.Equal(t, peers, peers2)
		for i := 0; i < 20; i++ {
			x := rand.Int63()
			peers_x := p.GetPermutationPeers(x)
			peers_x2 := p.GetPermutationPeers(x)
			assert.Equal(t, peers_x, peers_x2)

			y := rand.Int63()
			require.NotEqual(t, x, y)
			peers_y := p.GetPermutationPeers(y)
			assert.NotEqual(t, peers_x, peers_y)
		}
	})

}

func TestPeerServiceOnMemory(t *testing.T) {
	peerService := NewPeerServiceOnMemory()
	testPeerService(t, peerService)
}
