package test_utils

import (
	"github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/model"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func RandomProposalWithHeightRound(t *testing.T, height int64, round int32) model.Proposal {
	block, err := convertor.NewModelFactory().NewBlock(height, RandomByte(), rand.Int63(), RandomValidTxs(t))
	require.NoError(t, err)
	ValidSign(t, block)
	proposal, err := convertor.NewModelFactory().NewProposal(block, round)
	require.NoError(t, err)
	return proposal
}

func TimeParseDuration(t *testing.T, s string) time.Duration {
	d, err := time.ParseDuration(s)
	require.NoError(t, err)
	return d
}
