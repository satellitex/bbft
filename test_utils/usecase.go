package test_utils

import (
	"github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func RandomProposalWithRound(t *testing.T, round int32) model.Proposal {
	proposal, err := convertor.NewModelFactory().NewProposal(RandomValidBlock(t), round)
	require.NoError(t, err)
	return proposal
}
