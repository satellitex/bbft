package model

type Signature interface {
	GetPubKey() []byte
	GetSignature() []byte
}
