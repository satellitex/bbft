package model

type Block interface {
	GetHeader() BlockHeader
	GetTransactions() []Transaction
	GetSignature() Signature
	GetHash() ([]byte, error)
	Verify() error
	Sign(pubKey []byte, privKey []byte) error
}

type BlockHeader interface {
	GetHeight() int64
	GetPreBlockHash() []byte
	GetCreatedTime() int64
	GetHash() ([]byte, error)
}

type Proposal interface {
	GetBlock() Block
	GetRound() int64
}
