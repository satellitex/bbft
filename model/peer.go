package model

type Peer interface {
	GetAddress() string
	GetPubkey() []byte
}
