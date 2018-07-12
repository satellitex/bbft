package convertor

import (
	"github.com/satellitex/bbft/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBlock_GetHash(t *testing.T) {
	blocks := make([]model.Block, 20)
	for id, _ := range blocks {
		blocks[id] = randomBlock(t)
	}
	for id, a := range blocks {
		for jd, b := range blocks {
			if id != jd {
				assert.NotEqual(t, getHash(t, a), getHash(t, b))
			} else {
				assert.Equal(t, getHash(t, a), getHash(t, b))
			}
		}
	}
}

func TestBlock_SignAndVerify(t *testing.T) {
	t.Run("success valid key and valid txs", func(t *testing.T) {
		validPub, validPri := NewKeyPair()

		block := randomValidBlock(t)
		err := block.Sign(validPub, validPri)
		require.NoError(t, err)

		assert.True(t, block.Verify())
	})
	t.Run("success valid key and inValid txs", func(t *testing.T) {
		validPub, validPri := NewKeyPair()

		block := randomInvalidBlock(t)
		err := block.Sign(validPub, validPri)
		require.NoError(t, err)

		assert.False(t, block.Verify())
	})
	t.Run("failed invalid key and valid block", func(t *testing.T) {
		inValidPub := randomByte()
		inValidPriv := randomByte()

		block := randomValidBlock(t)
		err := block.Sign(inValidPub, inValidPriv)
		require.Error(t, err)
	})
	t.Run("failed invalid key and invalid block", func(t *testing.T) {
		inValidPub := randomByte()
		inValidPriv := randomByte()

		block := randomInvalidBlock(t)
		err := block.Sign(inValidPub, inValidPriv)
		require.Error(t, err)
	})
}
