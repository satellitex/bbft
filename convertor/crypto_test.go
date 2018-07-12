package convertor

import (
	"testing"

	"fmt"
	"github.com/pkg/errors"
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
	assert.True(t, Verify(pubkey, hash, signature),
		"pubkey: %x \nhash: %x\nsignature %x", pubkey, hash, signature)
}

func TestFailedSign(t *testing.T) {
	hash := CalcHash([]byte("a"))
	_, err := Sign(nil, hash)
	assert.EqualError(t, errors.Cause(err), ErrCryptoSign.Error())
}

func TestFailedVerifyNilPubkey(t *testing.T) {
	_, privkey := NewKeyPair()
	hash := CalcHash([]byte("a"))
	signature, err := Sign(privkey, hash)
	require.NoError(t, err)
	assert.False(t, Verify(nil, hash, signature),
		"pubkey: %x \nhash: %x\nsignature %x", nil, hash, signature)
}

func TestFailedVerifyNilSignature(t *testing.T) {
	pubkey, privkey := NewKeyPair()
	hash := CalcHash([]byte("a"))
	_, err := Sign(privkey, hash)
	require.NoError(t, err)
	assert.False(t, Verify(pubkey, hash, nil),
		"pubkey: %x \nhash: %x\nsignature %x", pubkey, hash, nil)
}

func TestFailedVerifyNilHash(t *testing.T) {
	pubkey, privkey := NewKeyPair()
	hash := CalcHash([]byte("a"))
	signature, err := Sign(privkey, hash)
	require.NoError(t, err)
	assert.False(t, Verify(pubkey, nil, signature),
		"pubkey: %x \nhash: %x\nsignature %x", pubkey, hash, signature)
}
