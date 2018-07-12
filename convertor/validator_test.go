package convertor

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func validSignedBlock(t *testing.T) model.Block {
	validPub, validPri := NewKeyPair()
	block := randomValidBlock(t)

	err := block.Sign(validPub, validPri)
	require.NoError(t, err)
	require.NoError(t, block.Verify())
	return block
}

func invalidSingedBlock(t *testing.T) model.Block {
	validPub, validPri := NewKeyPair()
	block := randomInvalidBlock(t)

	err := block.Sign(validPub, validPri)
	require.NoError(t, err)
	require.NoError(t, block.Verify())
	return block
}

func invalidErrSignedBlock(t *testing.T) model.Block {
	inValidPub := randomByte()
	inValidPriv := randomByte()
	block := randomInvalidBlock(t)

	err := block.Sign(inValidPub, inValidPriv)
	require.Error(t, err)
	require.Error(t, block.Verify())
	return block
}

func validErrSignedBlock(t *testing.T) model.Block {
	inValidPub := randomByte()
	inValidPriv := randomByte()
	block := randomInvalidBlock(t)

	err := block.Sign(inValidPub, inValidPriv)
	require.Error(t, err)
	require.Error(t, block.Verify())
	return block
}

func TestStatelessValidator_Validate(t *testing.T) {
	slv := NewStatelessValidator()
	t.Run("success valid key and valid txs", func(t *testing.T) {
		block := validSignedBlock(t)
		assert.NoError(t, slv.Validate(block))
	})
	t.Run("failed valid key and inValid txs", func(t *testing.T) {
		block := invalidSingedBlock(t)
		assert.EqualError(t, errors.Cause(slv.Validate(block)), ErrStatelessValidate.Error())
	})
	t.Run("failed invalid key and valid block", func(t *testing.T) {
		block := validErrSignedBlock(t)
		assert.EqualError(t, errors.Cause(slv.Validate(block)), ErrStatelessValidate.Error())
	})
	t.Run("failed invalid key and invalid block", func(t *testing.T) {
		block := invalidErrSignedBlock(t)
		assert.EqualError(t, errors.Cause(slv.Validate(block)), ErrStatelessValidate.Error())
	})
}
