package proto

type Block interface {
	GetHeader() BlockHeader
	GetTransactions() []Transaction
	GetSignature() Signature
}

type BlockHeader interface {
	GetHeight() int64
	GetPreBlockHash() []byte
	GetCreatedTime() int64
}

type Proposal interface {
	GetRound() int64
}
