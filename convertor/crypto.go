package convertor

import (
	"crypto/rand"
	"crypto/sha256"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ed25519"
)

var (
	ErrCryptoSign = errors.Errorf("Failed Sign by ed25519")
)

func CalcHash(b []byte) []byte {
	sha := sha256.New()
	sha.Write(b)
	return sha.Sum(nil)
}

func Verify(pubkey []byte, message []byte, signature []byte) error {
	if l := len(pubkey); l != ed25519.PublicKeySize {
		return errors.Errorf("ed25519: bad private key length: %d, expected %d",
			l, ed25519.PublicKeySize)
	}
	if ok := ed25519.Verify(pubkey, message, signature); !ok {
		return errors.Errorf("ed25519.Verify is invalid\n"+
			"pubkey: %x\n"+
			"message: %x\n"+
			"signature: %x",
			pubkey, message, signature)
	}
	return nil
}

func Sign(privkey []byte, message []byte) ([]byte, error) {
	if l := len(privkey); l != ed25519.PrivateKeySize {
		return nil, errors.Wrapf(ErrCryptoSign,
			"ed25519: bad private key length: %d, expected %d",
			l, ed25519.PrivateKeySize)
	}
	return ed25519.Sign(privkey, message), nil
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
