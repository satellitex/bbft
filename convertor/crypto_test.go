package convertor

import (
	"testing"

	"fmt"
	"github.com/stretchr/testify/assert"
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
	signature := Sign(privkey, hash)
	assert.True(t, Verify(pubkey, hash, signature),
		"pubkey: %x \nhash: %x\nsignature %x", pubkey, hash, signature)
}
