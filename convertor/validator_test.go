package convertor

import (
	"github.com/satellitex/bbft/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func validSignedBlock(t *testing.T) model.Block {
	validPub, validPri := NewKeyPair()
	block := randomValidBlock(t)

	err := block.Sign(validPub, validPri)
	require.NoError(t, err)
	require.True(t, block.Verify())
	return block
}

func invalidSingedBlock(t *testing.T) model.Block {
	validPub, validPri := NewKeyPair()
	block := randomInvalidBlock(t)

	err := block.Sign(validPub, validPri)
	require.NoError(t, err)
	require.True(t, block.Verify())
	return block
}

func invalidErrSignedBlock(t *testing.T) model.Block {
	inValidPub := randomByte()
	inValidPriv := randomByte()
	block := randomInvalidBlock(t)

	err := block.Sign(inValidPub, inValidPriv)
	require.Error(t, err)
	require.False(t, block.Verify())
	return block
}

func validErrSignedBlock(t *testing.T) model.Block {
	inValidPub := randomByte()
	inValidPriv := randomByte()
	block := randomInvalidBlock(t)

	err := block.Sign(inValidPub, inValidPriv)
	require.Error(t, err)
	require.False(t, block.Verify())
	return block
}

func TestStatelessValidator_Validate(t *testing.T) {
	slv := NewStatelessValidator()
	t.Run("success valid key and valid txs", func(t *testing.T) {
		block := validSignedBlock(t)
		require.True(t, slv.Validate(block))
	})
	t.Run("failed valid key and inValid txs", func(t *testing.T) {
		block := invalidSingedBlock(t)
		require.False(t, slv.Validate(block))
	})
	t.Run("failed invalid key and valid block", func(t *testing.T) {
		block := validErrSignedBlock(t)
		require.False(t, slv.Validate(block))
	})
	t.Run("failed invalid key and invalid block", func(t *testing.T) {
		block := invalidErrSignedBlock(t)
		require.False(t, slv.Validate(block))
	})
}
