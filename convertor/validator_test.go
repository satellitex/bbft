package convertor_test

import (
	"github.com/pkg/errors"
	. "github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/dba"
	"github.com/satellitex/bbft/model"
	. "github.com/satellitex/bbft/test_utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStatefulValidator_Validate(t *testing.T) {
	bc := dba.NewBlockChainOnMemory()
	sfv := NewStatefulValidator(bc)

	t.Run("success valid commitable Block", func(t *testing.T) {
		block := RandomCommitableBlock(t, bc)
		err := sfv.Validate(block)
		assert.NoError(t, err)

	})

	t.Run("failed nil Block", func(t *testing.T) {
		err := sfv.Validate(nil)
		assert.EqualError(t, errors.Cause(err), model.ErrInvalidBlock.Error())
	})

	t.Run("failed uncommitable Block", func(t *testing.T) {
		block := RandomBlock(t)
		err := sfv.Validate(block)
		assert.EqualError(t, errors.Cause(err), dba.ErrBlockChainVerifyCommit.Error())
	})

	t.Run("fialed invalid getHash Tx including Block", func(t *testing.T) {
		block := RandomCommitableBlock(t, bc)
		block.(*Block).Transactions[0].Payload = nil
		err := sfv.Validate(block)
		assert.EqualError(t, errors.Cause(err), model.ErrTransactionGetHash.Error())
	})

	// commit block
	commitBlock := RandomCommitableBlock(t, bc)
	bc.Commit(commitBlock)

	t.Run("succes valid commitable Block, exist bc", func(t *testing.T) {
		block := RandomCommitableBlock(t, bc)
		err := sfv.Validate(block)
		assert.NoError(t, err)
	})

	t.Run("failed alrady exist Tx including Block", func(t *testing.T) {
		block := RandomCommitableBlock(t, bc)
		block.(*Block).Transactions = commitBlock.(*Block).Transactions
		err := sfv.Validate(block)
		MultiErrorInCheck(t, err, ErrStatefulValidateAlreadyExistTx)
	})
}

func TestStatelessValidator_Validate(t *testing.T) {
	slv := NewStatelessValidator()
	t.Run("success valid key and valid txs", func(t *testing.T) {
		block := ValidSignedBlock(t)
		assert.NoError(t, slv.Validate(block))
	})
	t.Run("failed valid key and inValid txs", func(t *testing.T) {
		block := InvalidSingedBlock(t)
		MultiErrorInCheck(t, slv.Validate(block), model.ErrTransactionVerify)
	})
	t.Run("failed invalid key and valid block", func(t *testing.T) {
		block := ValidErrSignedBlock(t)
		MultiErrorInCheck(t, slv.Validate(block), model.ErrBlockVerify)
	})
	t.Run("failed invalid key and invalid block", func(t *testing.T) {
		block := InvalidErrSignedBlock(t)
		MultiErrorInCheck(t, errors.Cause(slv.Validate(block)), model.ErrTransactionVerify)
	})
	t.Run("failed nil block", func(t *testing.T) {
		assert.EqualError(t, errors.Cause(slv.Validate(nil)), model.ErrInvalidBlock.Error())
	})
}
