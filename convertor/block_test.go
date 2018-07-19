package convertor_test

import (
	"github.com/pkg/errors"
	. "github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/model"
	. "github.com/satellitex/bbft/test_utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBlock_GetHash(t *testing.T) {
	blocks := make([]model.Block, 20)
	for id, _ := range blocks {
		blocks[id] = RandomBlock(t)
	}
	for id, a := range blocks {
		for jd, b := range blocks {
			if id != jd {
				assert.NotEqual(t, GetHash(t, a), GetHash(t, b))
			} else {
				assert.Equal(t, GetHash(t, a), GetHash(t, b))
			}
		}
	}
}

type EvilTx struct {
	*Transaction
}

func TestBlock_FailedGetHash(t *testing.T) {
	t.Run("failed nil header", func(t *testing.T) {
		block := ValidSignedBlock(t)
		block.(*Block).Header = nil

		_, err := block.GetHash()
		assert.EqualError(t, errors.Cause(err), model.ErrBlockHeaderGetHash.Error())
	})
	t.Run("failed nil tx in transactions", func(t *testing.T) {
		block := ValidSignedBlock(t)
		block.(*Block).Transactions[0] = nil

		_, err := block.GetHash()
		assert.EqualError(t, errors.Cause(err), model.ErrTransactionGetHash.Error())
	})
	t.Run("failed nil bbft Block", func(t *testing.T) {
		block := ValidSignedBlock(t)
		block.(*Block).Block = nil

		_, err := block.GetHash()
		assert.Error(t, err)
	})
}

func TestBlock_SignAndVerify(t *testing.T) {
	t.Run("success valid key and valid txs", func(t *testing.T) {
		validPub, validPri := NewKeyPair()

		block := RandomValidBlock(t)
		err := block.Sign(validPub, validPri)
		assert.NoError(t, err)

		assert.NoError(t, block.Verify())
	})
	t.Run("success valid key and inValid txs", func(t *testing.T) {
		validPub, validPri := NewKeyPair()

		block := RandomInvalidBlock(t)
		err := block.Sign(validPub, validPri)
		assert.NoError(t, err)

		assert.NoError(t, block.Verify())
	})
	t.Run("failed invalid key and valid block", func(t *testing.T) {
		inValidPub := RandomByte()
		inValidPriv := RandomByte()

		block := RandomValidBlock(t)
		err := block.Sign(inValidPub, inValidPriv)
		assert.Error(t, err)

		assert.EqualError(t, errors.Cause(block.Verify()), ErrCryptoVerify.Error())
	})
	t.Run("failed invalid key and invalid block", func(t *testing.T) {
		inValidPub := RandomByte()
		inValidPriv := RandomByte()

		block := RandomInvalidBlock(t)
		err := block.Sign(inValidPub, inValidPriv)
		assert.Error(t, err)

		assert.EqualError(t, errors.Cause(block.Verify()), ErrCryptoVerify.Error())
	})
	t.Run("failed nil signature", func(t *testing.T) {
		block := ValidSignedBlock(t)
		block.(*Block).Signature = nil

		assert.EqualError(t, errors.Cause(block.Verify()), model.ErrInvalidSignature.Error())
	})
	t.Run("failed nil header", func(t *testing.T) {
		block := ValidSignedBlock(t)
		block.(*Block).Header = nil

		assert.EqualError(t, errors.Cause(block.Verify()), model.ErrBlockGetHash.Error())
	})
	t.Run("failed nil tx in transactions", func(t *testing.T) {
		block := ValidSignedBlock(t)
		block.(*Block).Transactions[0] = nil

		assert.EqualError(t, errors.Cause(block.Verify()), model.ErrBlockGetHash.Error())
	})
}
