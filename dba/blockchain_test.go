package dba_test

import (
	. "github.com/satellitex/bbft/dba"
	. "github.com/satellitex/bbft/test_utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/convertor"
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

func testBlockChain_VerifyCommit(t *testing.T, bc BlockChain) {
	t.Run("failed empty bc and add unverified Block", func (t *testing.T){
		block := ValidSignedBlock(t)

		err := bc.VerifyCommit(block)
		assert.EqualError(t, errors.Cause(err), ErrBlockChainVerifyCommit.Error())
	})

	t.Run("failed empty bc and add verified Block, but height = 0", func (t *testing.T){
		block := ValidSignedBlock(t)
		block.(*convertor.Block).Header.Height = 100

		err := bc.VerifyCommit(block)
		assert.EqualError(t, errors.Cause(err), ErrBlockChainVerifyCommit.Error())
	})

	t.Run("success empty bc and add comittable block", func (t *testing.T){
		block := RandomCommitableBlock(t, bc)

		err := bc.VerifyCommit(block)
		assert.EqualError(t, errors.Cause(err), ErrBlockChainVerifyCommit.Error())
	})

	t.Run("failed exist bc and add unverified Block", func (t *testing.T){
		block := ValidSignedBlock(t)

		err := bc.VerifyCommit(block)
		assert.EqualError(t, errors.Cause(err), ErrBlockChainVerifyCommit.Error())
	})

	t.Run("failed exist bc and add verified Block, but height = 0", func (t *testing.T){
		block := ValidSignedBlock(t)
		block.(*convertor.Block).Header.Height = 100

		err := bc.VerifyCommit(block)
		assert.EqualError(t, errors.Cause(err), ErrBlockChainVerifyCommit.Error())
	})

	t.Run("failed exist bc and add verified Block, but height = 0", func (t *testing.T){
		block := ValidSignedBlock(t)
		block.(*convertor.Block).Header.Height = 100

		err := bc.VerifyCommit(block)
		assert.EqualError(t, errors.Cause(err), ErrBlockChainVerifyCommit.Error())
	})

}

func TestBlockChainOnMemory_Commit(t *testing.T) {
	bc := NewBlockChainOnMemory()
	testBlockChain_VerifyCommit(t, bc)
}

func TestBlockChainOnMemory_Top(t *testing.T) {
	bc := NewBlockChainOnMemory()
	testBlockChain_Top(t, bc)
}
