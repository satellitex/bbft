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

// CalcSha256 calculates a SHA3-256 hash.
func CalcSha256(b []byte) []byte {
	hash := make([]byte, 0, 32)
	sha := sha256.New()
	// Note that Write method of hash.Hash interface never return an error.
	// See the documentation of hash.Hash.
	sha.Write(b)
	return sha.Sum(hash)
}

func CalcHashFromProto(msg proto.Message) ([]byte, error) {
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
	return ed25519.Verify(pubkey, hash, signature)
}

func NewKeyPair() (ed25519.PublicKey, ed25519.PrivateKey) {
	a, b, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(fmt.Errorf("ed25519.GenerateKey(rand.Reader) failed: %v", err))
	}
	return a, b
}
