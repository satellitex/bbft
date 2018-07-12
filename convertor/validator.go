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
)

type StatefulValidator struct {
	bc dba.BlockChain
}

func (v *StatefulValidator) Validate(block model.Block) error {
	bc, ok := v.bc.(*dba.BlockChainOnMemory)
	if !ok {
		return errors.Wrapf(ErrStatefulValidate,
			"Can not cast dba.BlockChainOnMemory %#v", v.bc)
	}
	if id, ok := bc.GetIndex(block); ok {
		return errors.Wrapf(ErrStatefulValidate,
			"Block %v is Already exist %d-th Blcok", block, id)
	}
	return nil
}

type StatelessValidator struct {
}

func NewStatelessValidator() model.StatelessValidator {
	return &StatelessValidator{}
}

func (v *StatelessValidator) Validate(block model.Block) error {
	var result error
	for _, tx := range block.GetTransactions() {
		if err := tx.Verify(); err != nil {
			result = multierr.Append(result, err)
		}
	}
	if err := block.Verify(); err != nil {
		result = multierr.Append(result, err)
	}
	if result != nil {
		return errors.Wrapf(ErrStatelessValidate, result.Error())
	}
	return nil
}
