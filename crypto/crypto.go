package crypto

import (
	"crypto/rand"
	"fmt"

	"crypto/sha256"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ed25519"
)

var (
	ErrMarshalProtocolBuffer = errors.New("failed to marshal protocol buffer")
)

const (
	HashSize = sha256.Size
)

type Hash [HashSize]byte
type HashPtr *Hash

// CalcSha256 calculates a SHA3-256 hash.
func CalcSha256(b []byte) HashPtr {
	var hash Hash = sha256.Sum256(b)
	return &hash
}

func CalcHashFromProto(msg proto.Message) (HashPtr, error) {
	pb, err := proto.Marshal(msg)
	if err != nil {
		return nil, errors.Wrap(ErrMarshalProtocolBuffer, err.Error())
	}
	return CalcSha256(pb), nil
}

func VerifyFromProto(pubkey []byte, msg proto.Message, signature []byte) bool {
	hash, err := CalcHashFromProto(msg)
	if err != nil {
		return false
	}
	return ed25519.Verify(pubkey, (*hash)[:], signature)
}

func NewKeyPair() (ed25519.PublicKey, ed25519.PrivateKey) {
	a, b, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(fmt.Errorf("ed25519.GenerateKey(rand.Reader) failed: %v", err))
	}
	return a, b
}
