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

func (v *StatelessValidator) Validate(block model.Block) bool {
	return block.Verify()
}