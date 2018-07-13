package test_utils

import (
	"github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/dba"
	"github.com/satellitex/bbft/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func RandomCommitableBlock(t *testing.T, bc dba.BlockChain) model.Block {
	if pre, ok := bc.Top(); ok {
		block, err := convertor.NewModelFactory().NewBlock(
			pre.GetHeader().GetHeight()+1,
			GetHash(t, pre),
			pre.GetHeader().GetCreatedTime()+10,
			RandomValidTxs(t),
		)
		require.NoError(t, err)
		validPub, validPri := convertor.NewKeyPair()
		block.Sign(validPub, validPri)
		return block
	}
	return RandomValidBlock(t)
}

func RandomCommitBlock(t *testing.T, bc dba.BlockChain, height int) dba.BlockChain {
	for i := 0; i < height; i++ {
		block := RandomCommitableBlock(t, bc)
		err := bc.Commit(block)
		require.NoError(t,err)
	}
	return bc
}