package model

type Block interface {
	GetHeader() BlockHeader
	GetTransactions() []Transaction
	GetSignature() Signature
	GetHash() ([]byte, error)
	Verify() bool
}

type BlockHeader interface {
	GetHeight() int64
	GetPreBlockHash() []byte
	GetCreatedTime() int64
	GetHash() ([]byte, error)
}

type Proposal interface {
	GetRound() int64
}
