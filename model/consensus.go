package model

type VoteMessage interface {
	GetBlockHash() []byte
	GetSignature() Signature
}

type ConsensusSender interface {
	Propagate(tx Transaction) error
	Propose(proposal Proposal) error
	Vote(vote VoteMessage) error
	PreCommit(vote VoteMessage) error
}
