package model

type Transaction interface {
	GetPayload() TransactionPayload
	GetSignatures() []Signature
	GetHash() ([]byte, error)
	Verify() error
}

type TransactionPayload interface {
	GetMessage() string
}
