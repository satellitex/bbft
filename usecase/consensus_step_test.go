package usecase_test

import (
	"github.com/satellitex/bbft/dba"
	"testing"

	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
	. "github.com/satellitex/bbft/test_utils"
	. "github.com/satellitex/bbft/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProposalFinder(t *testing.T) {
	finder := NewProposalFinder()

	t.Run("success", func(t *testing.T) {
		proposal := RandomProposalWithHeightRound(t, 0, 1)
		err := finder.Set(proposal)
		assert.NoError(t, err)
		actual, ok := finder.Find(0, 1)
		assert.True(t, ok)
		assert.Equal(t, proposal, actual)
	})

	t.Run("failed nil vote", func(t *testing.T) {
		err := finder.Set(nil)
		assert.EqualError(t, errors.Cause(err), model.ErrInvalidProposal.Error())

		_, ok := finder.Find(0, 0)
		assert.False(t, ok)
	})

	t.Run("Clear", func(t *testing.T) {
		expectedProposals := []model.Proposal{
			RandomProposalWithHeightRound(t, 1, 0),
			RandomProposalWithHeightRound(t, 1, 1),
			RandomProposalWithHeightRound(t, 1, 2),
			RandomProposalWithHeightRound(t, 1, 3),
		}
		for i := 0; i < 4; i++ {
			require.NoError(t, finder.Set(expectedProposals[i]))
			actual, ok := finder.Find(1, int32(i))
			require.True(t, ok)
			require.Equal(t, expectedProposals[i], actual)
		}

		finder.Clear(2)

		for i := 0; i < 4; i++ {
			_, ok := finder.Find(1, int32(i))
			assert.False(t, ok)
		}
	})

}

func TestPreCommitFinder(t *testing.T) {
	ps := dba.NewPeerServiceOnMemory()

	conf := GetTestConfig()
	for i := 0; i < 4; i++ {
		ps.AddPeer(RandomPeerWithPriv())
	}

	finder := NewPreCommitFinder(ps, conf)
	t.Run("success", func(t *testing.T) {
		err := finder.Set(RandomVoteMessage(t))
		assert.NoError(t, err)
		_, ok := finder.Get()
		assert.False(t, ok)
	})

	t.Run("failed nil vote", func(t *testing.T) {
		err := finder.Set(nil)
		assert.EqualError(t, errors.Cause(err), model.ErrInvalidVoteMessage.Error())
		_, ok := finder.Get()
		assert.False(t, ok)
	})

	t.Run("many set", func(t *testing.T) {
		for i := 0; i < conf.PreCommitFinderLimits*2; i++ {
			err := finder.Set(RandomVoteMessage(t))
			assert.NoError(t, err)
		}
	})

	t.Run("collec Get", func(t *testing.T) {
		vote := RandomVoteMessage(t)
		for i := 0; i < 2; i++ {
			assert.NoError(t, finder.Set(vote))
			_, ok := finder.Get()
			assert.False(t, ok)
		}
		assert.NoError(t, finder.Set(vote))
		hash, ok := finder.Get()
		assert.True(t, ok)

		assert.Equal(t, vote.GetBlockHash(), hash)

		_, ok = finder.Get()
		assert.False(t, ok)
	})

}
