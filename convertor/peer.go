package convertor

type Peer struct {
	Address string
	Pubkey  []byte
}

func (p *Peer) GetAddress() string {
	return p.Address
}

func (p *Peer) GetPubkey() []byte {
	return p.Pubkey
}
