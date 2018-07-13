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

func TestBlock_FailedGetHash(t *testing.T) {
	t.Run("failed nil header", func(t *testing.T) {
		block := randomValidBlock(t)
		block.(*Block).Header = nil

		_, err := block.GetHash()
		assert.EqualError(t, errors.Cause(err), ErrBlockGetHash.Error())
	})
	t.Run("failed nil tx in transactions", func(t *testing.T) {
		block := randomValidBlock(t)
		for id, _ := range block.GetTransactions() {
			block.(*Block).Transactions[id] = nil
		}

		hash, err := block.GetHash()
		assert.EqualError(t, errors.Cause(err), ErrBlockGetHash.Error(), "%x", hash)
	})
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
	t.Run("failed nil header", func(t *testing.T) {
		block := randomValidBlock(t)
		block.(*Block).Header = nil

		assert.EqualError(t, errors.Cause(block.Verify()), ErrBlockVerify.Error())
	})
	t.Run("failed nil tx in transactions", func(t *testing.T) {
		block := randomValidBlock(t)
		block.GetTransactions()[0].(*Transaction).Transaction = nil

		assert.EqualError(t, errors.Cause(block.Verify()), ErrBlockVerify.Error())
	})
}
