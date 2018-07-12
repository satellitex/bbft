package convertor

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
	"github.com/stretchr/testify/assert"
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
		assert.NoError(t, err)

		assert.NoError(t, block.Verify())
	})
	t.Run("success valid key and inValid txs", func(t *testing.T) {
		validPub, validPri := NewKeyPair()

		block := randomInvalidBlock(t)
		err := block.Sign(validPub, validPri)
		assert.NoError(t, err)

		assert.NoError(t, block.Verify())
	})
	t.Run("failed invalid key and valid block", func(t *testing.T) {
		inValidPub := randomByte()
		inValidPriv := randomByte()

		block := randomValidBlock(t)
		err := block.Sign(inValidPub, inValidPriv)
		assert.Error(t, err)

		assert.EqualError(t, errors.Cause(block.Verify()), ErrBlockVerify.Error())
	})
	t.Run("failed invalid key and invalid block", func(t *testing.T) {
		inValidPub := randomByte()
		inValidPriv := randomByte()

		block := randomInvalidBlock(t)
		err := block.Sign(inValidPub, inValidPriv)
		assert.Error(t, err)

		assert.EqualError(t, errors.Cause(block.Verify()), ErrBlockVerify.Error())
	})
	t.Run("failed nil signature", func(t *testing.T) {
		block := randomValidBlock(t)
		block.(*Block).Signature = nil

		assert.EqualError(t, errors.Cause(block.Verify()), ErrBlockVerify.Error())
	})
}
