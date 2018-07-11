package model

import "github.com/satellitex/bbft/crypto"

type Transaction interface {
	GetPaylaod() TransactionPayload
	GetSignatures() []Signature
	GetHash() (crypto.HashPtr, error)
}

type TransactionPayload interface {
	GetMessage() string
}
