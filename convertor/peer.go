package convertor

type Peer struct {
	address string
	pubkey  []byte
}

func (p *Peer) GetAddress() string {
	return p.address
}

func (p *Peer) GetPubkey() []byte {
	return p.pubkey
}
