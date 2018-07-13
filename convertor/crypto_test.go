package convertor_test

import (
	"testing"

	"fmt"
	. "github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalcHash(t *testing.T) {
	message := "a"
	expectedHash := "ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb"
	actualHash := CalcHash([]byte(message))

	assert.Equal(t, expectedHash, fmt.Sprintf("%x", actualHash))
}

func TestSignAndVerify(t *testing.T) {
	pubkey, privkey := NewKeyPair()
	hash := CalcHash([]byte("a"))
	signature, err := Sign(privkey, hash)
	require.NoError(t, err)
	assert.NoError(t, Verify(pubkey, hash, signature),
		"pubkey: %x \nhash: %x\nsignature %x", pubkey, hash, signature)
}

func TestFailedSign(t *testing.T) {
	hash := CalcHash([]byte("a"))
	_, err := Sign(nil, hash)
	assert.Error(t, err)
}

func TestFailedVerifyNilPubkey(t *testing.T) {
	_, privkey := NewKeyPair()
	hash := CalcHash([]byte("a"))
	signature, err := Sign(privkey, hash)
	require.NoError(t, err)
	assert.Error(t, Verify(nil, hash, signature),
		"pubkey: %x \nhash: %x\nsignature %x", nil, hash, signature)
}

func TestFailedVerifyNilSignature(t *testing.T) {
	pubkey, privkey := NewKeyPair()
	hash := CalcHash([]byte("a"))
	_, err := Sign(privkey, hash)
	require.NoError(t, err)
	assert.Error(t, Verify(pubkey, hash, nil),
		"pubkey: %x \nhash: %x\nsignature %x", pubkey, hash, nil)
}

func TestFailedVerifyNilHash(t *testing.T) {
	pubkey, privkey := NewKeyPair()
	hash := CalcHash([]byte("a"))
	signature, err := Sign(privkey, hash)
	require.NoError(t, err)
	assert.Error(t, Verify(pubkey, nil, signature),
		"pubkey: %x \nhash: %x\nsignature %x", pubkey, nil, signature)
}

func TestSuccessdVerifyNilMarshal(t *testing.T) {
	tx := &Transaction{&bbft.Transaction{}}
	_, err := CalcHashFromProto(tx)

	assert.NoError(t, err)
}

func TestFailedVerifyNilMarshal(t *testing.T) {
	tx := &Transaction{}
	_, err := CalcHashFromProto(tx)

	assert.Error(t, err)
}
