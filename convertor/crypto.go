package convertor

import (
	"crypto/rand"
	"crypto/sha256"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ed25519"
)

func CalcHash(b []byte) []byte {
	sha := sha256.New()
	sha.Write(b)
	return sha.Sum(nil)
}

func Verify(pubkey []byte, message []byte, signature []byte) bool {
	return ed25519.Verify(pubkey, message, signature)
}

func Sign(privkey []byte, message []byte) []byte {
	return ed25519.Sign(privkey, message)
}

func NewKeyPair() ([]byte, []byte) {
	a, b, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(errors.Errorf("ed25519.GenerateKey(rand.Reader) failed: %v", err))
	}
	return a, b
}

var (
	ErrMarshalProtocolBuffer = errors.New("failed to marshal protocol buffer")
)

func CalcHashFromProto(msg proto.Message) ([]byte, error) {
	pb, err := proto.Marshal(msg)
	if err != nil {
		return nil, errors.Wrap(ErrMarshalProtocolBuffer, err.Error())
	}
	return CalcHash(pb), nil
}

func VerifyFromProto(pubkey []byte, msg proto.Message, signature []byte) bool {
	hash, err := CalcHashFromProto(msg)
	if err != nil {
		return false
	}
	return Verify(pubkey, hash, signature)
}
