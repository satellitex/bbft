package proto

type Transaction interface {
	GetPaylaod() TransactionPayload
	GetSignatures() []Signature
}

type TransactionPayload interface {
	GetMessage() string
}
