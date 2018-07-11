package proto

type ProposalTx interface {
	GetTransaction() Transaction
	GetSignature() Signature
}

type VoteMessage interface {
	GetBlockHash() []byte
	GetSignature() Signature
}
