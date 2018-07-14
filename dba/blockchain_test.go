package dba_test

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/convertor"
	. "github.com/satellitex/bbft/dba"
	"github.com/satellitex/bbft/model"
	. "github.com/satellitex/bbft/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func testBlockChain_Top(t *testing.T, bc BlockChain) {
	// Empty case
	_, ok := bc.Top()
	assert.False(t, ok)

	for i := 0; i < 10; i++ {
		commitableBlock := RandomCommitableBlock(t, bc)
		bc.Commit(commitableBlock)
		top, ok := bc.Top()
		assert.True(t, ok)
		assert.Equal(t, top, commitableBlock)
	}
}

type HackHashBlock struct {
	*convertor.Block
	hash []byte
}

func (e *HackHashBlock) GetHash() ([]byte, error) {
	return e.hash, nil
}

func testBlockChain_VerifyCommit(t *testing.T, bc BlockChain) {

	t.Run("failed empty bc and add unverified Block", func(t *testing.T) {
		block := ValidErrSignedBlock(t)

		err := bc.VerifyCommit(block)
		assert.EqualError(t, errors.Cause(err), model.ErrBlockVerify.Error())
	})

	t.Run("failed empty bc and add verified Block, but height = 0", func(t *testing.T) {
		block := RandomCommitableBlock(t, bc)
		block.(*convertor.Block).Header.Height = 100
		ValidSign(t,block)

		err := bc.VerifyCommit(block)
		assert.EqualError(t, errors.Cause(err), ErrBlockChainVerifyCommitInvalidHeight.Error())
	})

	t.Run("success empty bc and add comittable block", func(t *testing.T) {
		block := RandomCommitableBlock(t, bc)

		err := bc.VerifyCommit(block)
		assert.NoError(t, err)
	})

	// Commit 1 Block
	require.NoError(t, bc.Commit(RandomCommitableBlock(t, bc)))

	t.Run("failed exist bc and add unverified Block", func(t *testing.T) {
		block := ValidErrSignedBlock(t)

		err := bc.VerifyCommit(block)
		assert.EqualError(t, errors.Cause(err), model.ErrBlockVerify.Error())
	})

	t.Run("failed exist bc and add verified Block, but height = 0", func(t *testing.T) {
		block := ValidSignedBlock(t)
		block.(*convertor.Block).Header.Height = 100
		ValidSign(t,block)

		err := bc.VerifyCommit(block)
		assert.EqualError(t, errors.Cause(err), ErrBlockChainVerifyCommitInvalidHeight.Error())
	})

	t.Run("failed exist bc and add verified Block, but preblock is not top.block.hash ", func(t *testing.T) {
		block := RandomCommitableBlock(t, bc)
		block.(*convertor.Block).Header.PreBlockHash = RandomByte()
		ValidSign(t,block)

		err := bc.VerifyCommit(block)
		assert.EqualError(t, errors.Cause(err), ErrBlockChainVerifyCommitInvalidPreBlockHash.Error())
	})

	t.Run("failed exist bc and add verified Block, but createdTime so faster ", func(t *testing.T) {
		block := RandomCommitableBlock(t, bc)
		block.(*convertor.Block).Header.CreatedTime = 0
		ValidSign(t,block)

		err := bc.VerifyCommit(block)
		assert.EqualError(t, errors.Cause(err), ErrBlockChainVerifyCommitInvalidCreatedTime.Error())
	})

	t.Run("failed exist bc and add verified Block, but Already Exist", func(t *testing.T) {
		top, ok := bc.Top()
		require.True(t, ok)
		block := &HackHashBlock{RandomCommitableBlock(t, bc).(*convertor.Block), GetHash(t, top)}
		ValidSign(t,block)

		err := bc.VerifyCommit(block)
		assert.EqualError(t, errors.Cause(err), ErrBlockChainVerifyCommitAlreadyExist.Error())
	})

	t.Run("success exist bc and add comittable block", func(t *testing.T) {
		block := RandomCommitableBlock(t, bc)

		err := bc.VerifyCommit(block)
		assert.NoError(t, err)
	})

}

func testBlockChain_Commit(t *testing.T, bc BlockChain) {
	t.Run("failed commit invalid block", func(t *testing.T) {
		err := bc.Commit(RandomBlock(t))
		require.EqualError(t, errors.Cause(err), ErrBlockChainVerifyCommit.Error())
	})

	t.Run("failed commit valid block", func(t *testing.T) {
		err := bc.Commit(RandomCommitableBlock(t, bc))
		require.NoError(t, err)
	})
}

func TestBlockChainOnMemory_Top(t *testing.T) {
	bc := NewBlockChainOnMemory()
	testBlockChain_Top(t, bc)
}

func TestBlockChainOnMemory_VerifyCommit(t *testing.T) {
	bc := NewBlockChainOnMemory()
	testBlockChain_VerifyCommit(t, bc)
}

func TestBlockChainOnMemory_Commit(t *testing.T) {
	bc := NewBlockChainOnMemory()
	testBlockChain_Commit(t, bc)
}
