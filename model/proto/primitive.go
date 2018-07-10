package proto

type Signature interface {
	GetPubKey() []byte
	GetSignature() []byte
}