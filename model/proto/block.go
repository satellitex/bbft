package proto

type Block interface {
	GetHeader() BlockHeader
	GetTransactions() []Transaction
	GetSignature() Signature
}

type BlockHeader interface {
	GetHeight() int64
	GetPreBlockHash() []byte
	GetRound() int64
	GetCreatedTime() int64
}