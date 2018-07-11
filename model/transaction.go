package model

type Transaction interface {
	GetPayload() TransactionPayload
	GetSignatures() []Signature
	GetHash() ([]byte, error)
	Verify() bool
}

type TransactionPayload interface {
	GetMessage() string
}
