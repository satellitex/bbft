package test_utils

import (
	"github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/dba"
	"github.com/satellitex/bbft/model"
	"github.com/stretchr/testify/require"
	"math/rand"
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
	block := RandomValidBlock(t)
	block.(*convertor.Block).Header.Height = 0
	ValidSign(t, block)
	return block
}

func RandomProposal(t *testing.T) model.Proposal {
	proposal, err := convertor.NewModelFactory().NewProposal(RandomValidBlock(t), rand.Int63())
	require.NoError(t, err)
	return proposal
}

func RandomVoteMessage(t *testing.T) model.VoteMessage {
	vote := convertor.NewModelFactory().NewVoteMessage(RandomByte())
	ValidSign(t, vote)
	return vote
}

func ValidSignToBlock(b model.Block) model.Block {
	pub, pri := convertor.NewKeyPair()
	b.Sign(pub, pri)
	return b
}

type Signer interface {
	Sign(pub []byte, pri []byte) error
}

func ValidSign(t *testing.T, s Signer) {
	pub, pri := convertor.NewKeyPair()
	require.NoError(t, s.Sign(pub, pri))
}
