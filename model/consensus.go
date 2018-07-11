package model

type ProposalTx interface {
	GetTransaction() Transaction
	GetSignature() Signature
}

type VoteMessage interface {
	GetBlockHash() []byte
	GetSignature() Signature
}

type ConsensusSender interface {
	Propagate(ptx ProposalTx) error
	Propose(proposal Proposal) error
	Vote(vote VoteMessage) error
	PreCommit(vote VoteMessage) error
}