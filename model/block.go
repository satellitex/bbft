package model

import "github.com/satellitex/bbft/crypto"

type Block interface {
	GetHeader() BlockHeader
	GetTransactions() []Transaction
	GetSignature() Signature
	GetHash() (crypto.HashPtr, error)
}

type BlockHeader interface {
	GetHeight() int64
	GetPreBlockHash() []byte
	GetCreatedTime() int64
	GetHash() (crypto.HashPtr, error)
}

type Proposal interface {
	GetRound() int64
}