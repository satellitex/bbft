package convertor

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/dba"
	"github.com/satellitex/bbft/model"
	"go.uber.org/multierr"
)

var (
	ErrStatefulValidate  = errors.Errorf("Failed StatefulValidate")
	ErrStatelessValidate = errors.Errorf("Failed StatelessValidate")

	ErrStatefulValidateAlreadyExistTx = errors.Errorf("Failed Already Exist Transaction")
)

type StatefulValidator struct {
	bc dba.BlockChain
}

func (v *StatefulValidator) Validate(block model.Block) error {
	// bloc is already stateless validate
	var result error
	if block == nil {
		return errors.Wrapf(model.ErrInvalidBlock, "Block is nil")
	}
	if err := v.bc.VerifyCommit(block); err != nil {
		result = multierr.Append(result, errors.Wrapf(dba.ErrBlockChainVerifyCommit, err.Error()))
	}
	for _, tx := range block.GetTransactions() {
		hash, err := tx.GetHash()
		if err != nil {
			result = multierr.Append(result, errors.Wrapf(model.ErrTransactionGetHash, err.Error()))
			continue
		}
		if _, ok := v.bc.FindTx(hash); ok {
			result = multierr.Append(result, errors.Wrapf(ErrStatefulValidateAlreadyExistTx, "Alrady exist transaction hash : %x", hash))
		}
	}
	return result
}

func NewStatefulValidator(bc dba.BlockChain) model.StatefulValidator {
	return &StatefulValidator{bc}
}

type StatelessValidator struct {
}

func NewStatelessValidator() model.StatelessValidator {
	return &StatelessValidator{}
}

func (v *StatelessValidator) BlockValidate(block model.Block) error {
	var result error
	if block == nil {
		return errors.Wrapf(model.ErrInvalidBlock, "Block is nil")
	}
	for _, tx := range block.GetTransactions() {
		if err := v.TxValidate(tx); err != nil {
			result = multierr.Append(result, errors.Wrapf(model.ErrStatelessTxValidate, err.Error()))
		}
	}
	if err := block.Verify(); err != nil {
		result = multierr.Append(result, errors.Wrapf(model.ErrBlockVerify, err.Error()))
	}
	return result
}

func (v *StatelessValidator) TxValidate(tx model.Transaction) error {
	if tx == nil {
		return errors.Wrapf(model.ErrInvalidTransaction, "tx is nil")
	}
	if err := tx.Verify(); err != nil {
		return errors.Wrapf(model.ErrTransactionVerify, err.Error())
	}
	return nil
}
