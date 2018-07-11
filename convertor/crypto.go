package convertor

import (
	"crypto/rand"
	"crypto/sha256"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
	"golang.org/x/crypto/ed25519"
)

type Cryptor struct{}

func (_ *Cryptor) CalcHash(b []byte) []byte {
	sha := sha256.New()
	sha.Write(b)
	return sha.Sum(nil)
}

func (_ *Cryptor) Verify(pubkey []byte, message []byte, signature []byte) bool {
	return ed25519.Verify(pubkey, message, signature)
}

func (_ *Cryptor) NewKeyPair() ([]byte, []byte) {
	a, b, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(errors.Errorf("ed25519.GenerateKey(rand.Reader) failed: %v", err))
	}
	return a, b
}

var (
	ErrMarshalProtocolBuffer = errors.New("failed to marshal protocol buffer")
)

func CalcHashFromProto(msg proto.Message, c model.Cryptor) ([]byte, error) {
	pb, err := proto.Marshal(msg)
	if err != nil {
		return nil, errors.Wrap(ErrMarshalProtocolBuffer, err.Error())
	}
	return c.CalcHash(pb), nil
}

func VerifyFromProto(pubkey []byte, msg proto.Message, signature []byte, c model.Cryptor) bool {
	hashPtr, err := CalcHashFromProto(msg, c)
	if err != nil {
		return false
	}
	return c.Verify(pubkey, hashPtr, signature)
}
