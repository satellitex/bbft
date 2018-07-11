package model

type Transaction interface {
	GetPaylaod() TransactionPayload
	GetSignatures() []Signature
	GetHash() ([]byte, error)
	Verify() error
}

type TransactionPayload interface {
	GetMessage() string
}
