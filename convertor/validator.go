package convertor

import (
	"github.com/satellitex/bbft/dba"
	"github.com/satellitex/bbft/model"
)

type StatefulValidator struct {
	bc dba.BlockChain
}

func (v *StatefulValidator) Validate(block model.Block) bool {
	bc, ok := v.bc.(*dba.BlockChainOnMemory)
	if !ok {
		return false
	}
	_, ok = bc.GetIndex(block)
	return !ok
}

type StatelessValidator struct {
}

func NewStatelessValidator() model.StatelessValidator {
	return &StatelessValidator{}
}

func (v *StatelessValidator) Validate(block model.Block) bool {
	for _, tx := range block.GetTransactions() {
		if !tx.Verify() {
			return false
		}
	}
	return block.Verify()
}
