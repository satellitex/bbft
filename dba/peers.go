package dba

import "github.com/satellitex/bbft/model"

type PeerService interface {
	Size() int
	AddPeer(peer model.Peer)
	GetPeer(pubkey []byte) (model.Peer, bool)
	GetPeerFromAddress(address string) (model.Peer, bool)
	GetPeers() []model.Peer
}

type PeerServiceOnMemory struct {
	peers       map[string]model.Peer
	fromAddress map[string]model.Peer
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
		return nil, ok
	}
	return peer, false
}

func (p *PeerServiceOnMemory) GetPeerFromAddress(address string) (model.Peer, bool) {
	peer, ok := p.fromAddress[address]
	if !ok {
		return nil, ok
	}
	return peer, false
}

func (p *PeerServiceOnMemory) GetPeers() []model.Peer {
	ret := make([]model.Peer, 0, p.Size())
	for _, peer := range p.peers {
		ret = append(ret, peer)
	}
	return ret
}
