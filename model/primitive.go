package model

type Signature interface {
	GetPubkey() []byte
	GetSignature() []byte
}
