package dba

import (
	"github.com/satellitex/bbft/model"
	"math/rand"
	"sort"
)

type PeerService interface {
	Size() int
	AddPeer(peer model.Peer)
	GetPeer(pubkey []byte) (model.Peer, bool)
	GetPeerFromAddress(address string) (model.Peer, bool)
	GetPeers() []model.Peer
	GetNumberOfAllowedFailedPeers() int
	GetNumberOfRequiredAcceptPeers() int
	GetPermutationPeers(height int64) []model.Peer
}

type PeerServiceOnMemory struct {
	peers       map[string]model.Peer
	fromAddress map[string]model.Peer
}

func NewPeerServiceOnMemory() PeerService {
	return &PeerServiceOnMemory{
		make(map[string]model.Peer),
		make(map[string]model.Peer),
	}
}

func (p *PeerServiceOnMemory) Size() int {
	return len(p.peers)
}

func (p *PeerServiceOnMemory) AddPeer(peer model.Peer) {
	p.peers[string(peer.GetPubkey())] = peer
	p.fromAddress[string(peer.GetAddress())] = peer
}

func (p *PeerServiceOnMemory) GetPeer(pubkey []byte) (model.Peer, bool) {
	peer, ok := p.peers[string(pubkey)]
	if !ok {
		return nil, false
	}
	return peer, true
}

func (p *PeerServiceOnMemory) GetPeerFromAddress(address string) (model.Peer, bool) {
	peer, ok := p.fromAddress[address]
	if !ok {
		return nil, false
	}
	return peer, true
}

func (p *PeerServiceOnMemory) GetPeers() []model.Peer {
	keys := make([]string, 0, p.Size())
	for key, _ := range p.fromAddress {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	ret := make([]model.Peer, 0, p.Size())
	for _, key := range keys {
		ret = append(ret, p.fromAddress[key])
	}
	return ret
}

func (p *PeerServiceOnMemory) GetNumberOfAllowedFailedPeers() int {
	return (p.Size() - 1) / 3
}

func (p *PeerServiceOnMemory) GetNumberOfRequiredAcceptPeers() int {
	return p.GetNumberOfAllowedFailedPeers()*2 + 1
}

func (p *PeerServiceOnMemory) GetPermutationPeers(height int64) []model.Peer {
	rnd := rand.New(rand.NewSource(height))
	peers := p.GetPeers()
	rnd.Shuffle(len(peers), func(i, j int) {
		peers[i], peers[j] = peers[j], peers[i]
	})
	return peers
}
