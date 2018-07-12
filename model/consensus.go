package model

type VoteMessage interface {
	GetBlockHash() []byte
	GetSignature() Signature
	Sign(pubKey []byte, privKey []byte) error
	Verify() bool
}

type ConsensusSender interface {
	Propagate(tx Transaction) error
	Propose(proposal Proposal) error
	Vote(vote VoteMessage) error
	PreCommit(vote VoteMessage) error
}
