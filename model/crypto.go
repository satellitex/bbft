package model

type Cryptor interface {
	CalcHash(b []byte) []byte
	Verify(pubkey []byte, message []byte, signature []byte) bool
	NewKeyPair() ([]byte, []byte)
}
