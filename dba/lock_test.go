package dba_test

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/convertor"
	. "github.com/satellitex/bbft/dba"
	"github.com/satellitex/bbft/model"
	. "github.com/satellitex/bbft/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func testLock_RegisterProposal(t *testing.T, lock Lock) {
	t.Run("success register valid random proposal", func(t *testing.T) {
		err := lock.RegisterProposal(RandomProposal(t))
		assert.NoError(t, err)
	})
	t.Run("failed register nil", func(t *testing.T) {
		err := lock.RegisterProposal(nil)
		assert.EqualError(t, errors.Cause(err), model.ErrInvalidProposal.Error())
	})
	t.Run("failed invalid block GetHash", func(t *testing.T) {
		proposal := RandomProposal(t)
		proposal.(*convertor.Proposal).Block.Header = nil

		err := lock.RegisterProposal(proposal)
		assert.EqualError(t, errors.Cause(err), model.ErrBlockGetHash.Error())
	})
	t.Run("failed already register proposal", func(t *testing.T) {
		proposal := RandomProposal(t)
		err := lock.RegisterProposal(proposal)
		assert.NoError(t, err)

		err = lock.RegisterProposal(proposal)
		assert.EqualError(t, errors.Cause(err), ErrAlreadyRegisterProposal.Error())
	})

	t.Run("No Memory over flow, when RegisterProposal Attack", func(t *testing.T) {
		t.Skipf("This test is too long")
		for i := 0; i < 1000000; i++ {
			go func() {
				err := lock.RegisterProposal(RandomProposal(t))
				assert.NoError(t, err)
			}()
		}
	})
}

func testLock_AddVoteMessageAndGetLocked(t *testing.T, lock Lock, p PeerService) {
	// 4 peer
	peers := []model.Peer{
		RandomPeer(),
		RandomPeer(),
		RandomPeer(),
		RandomPeer(),
	}
	for _, peer := range peers {
		p.AddPeer(peer)
	}

	validGetLockedProposal := func(t *testing.T, expectedProposal model.Proposal) {
		proposal, ok := lock.GetLockedProposal(0)
		if expectedProposal == nil {
			assert.False(t, ok)
		} else {
			assert.True(t, ok)
			assert.Equal(t, expectedProposal, proposal)
		}
	}

	t.Run("success valid votes", func(t *testing.T) {
		vote := RandomVoteMessage(t)
		err := lock.AddVoteMessage(vote)
		assert.NoError(t, err)

		validGetLockedProposal(t, nil)
	})
	t.Run("faield nil vote", func(t *testing.T) {
		err := lock.AddVoteMessage(nil)
		assert.EqualError(t, errors.Cause(err), model.ErrInvalidVoteMessage.Error())

		validGetLockedProposal(t, nil)
	})

	// register valid proposal
	validProposals := []model.Proposal{
		RandomProposalWithHeightRound(t, 0, 0),
		RandomProposalWithHeightRound(t, 0, 1),
		RandomProposalWithHeightRound(t, 0, 2),
	}
	lock.RegisterProposal(validProposals[0])
	lock.RegisterProposal(validProposals[1])
	require.True(t, validProposals[0].GetRound() < validProposals[1].GetRound())

	validAddVote := func(t *testing.T, proposal model.Proposal) {
		vote := convertor.NewModelFactory().NewVoteMessage(GetHash(t, proposal.GetBlock()))
		ValidSign(t, vote)
		err := lock.AddVoteMessage(vote)
		require.NoError(t, err)
	}

	t.Run("success valid vote setLockedProposal", func(t *testing.T) {
		vp := validProposals[0]
		validAddVote(t, vp)
		validAddVote(t, vp)
		validGetLockedProposal(t, nil)

		vote := convertor.NewModelFactory().NewVoteMessage(GetHash(t, vp.GetBlock()))
		ValidSign(t, vote)
		err := lock.AddVoteMessage(vote)
		assert.NoError(t, err)

		validGetLockedProposal(t, vp)

		vote = convertor.NewModelFactory().NewVoteMessage(GetHash(t, vp.GetBlock()))
		ValidSign(t, vote)
		err = lock.AddVoteMessage(vote)
		assert.NoError(t, err)

		validGetLockedProposal(t, vp)
	})

	t.Run("success valid second vote setLockedProposal", func(t *testing.T) {
		vp := validProposals[1]
		validAddVote(t, vp)
		validAddVote(t, vp)
		validGetLockedProposal(t, validProposals[0])

		vote := convertor.NewModelFactory().NewVoteMessage(GetHash(t, vp.GetBlock()))
		ValidSign(t, vote)
		err := lock.AddVoteMessage(vote)
		assert.NoError(t, err)

		validGetLockedProposal(t, vp)
	})

	t.Run("success collect Proposal votes, before register proposal", func(t *testing.T) {
		vp := validProposals[2]
		validAddVote(t, vp)
		validAddVote(t, vp)
		validGetLockedProposal(t, validProposals[1])

		vote := convertor.NewModelFactory().NewVoteMessage(GetHash(t, vp.GetBlock()))
		ValidSign(t, vote)
		err := lock.AddVoteMessage(vote)
		assert.NoError(t, err)

		validGetLockedProposal(t, validProposals[1])

		err = lock.RegisterProposal(vp)
		assert.NoError(t, err)
		validGetLockedProposal(t, vp)
	})

	t.Run("failed alrady exist voteMessage", func(t *testing.T) {
		vote := convertor.NewModelFactory().NewVoteMessage(GetHash(t, RandomBlock(t)))
		ValidSign(t, vote)

		err := lock.AddVoteMessage(vote)
		assert.NoError(t, err)

		err = lock.AddVoteMessage(vote)
		assert.EqualError(t, errors.Cause(err), ErrAlreadyAddVoteMessage.Error())

	})

	t.Run("Execute clean lock", func(t *testing.T) {
		validGetLockedProposal(t, validProposals[2])

		lock.Clean(1)

		validGetLockedProposal(t, nil)
	})

	t.Run("No Memory over flow, when AddVote Attack", func(t *testing.T) {
		t.Skip("This test is too long")
		for i := 0; i < 1000000; i++ {
			go func() {
				vote := convertor.NewModelFactory().NewVoteMessage(GetHash(t, RandomBlock(t)))
				ValidSign(t, vote)

				err := lock.AddVoteMessage(vote)
				require.NoError(t, err)
			}()
		}
	})

}

func TestLockOnMemory_RegisterProposal(t *testing.T) {
	lock := NewLockOnMemory(NewPeerServiceOnMemory(), GetTestConfig())
	testLock_RegisterProposal(t, lock)
}

func TestLockOnMemory_AddVoteMessageAndGetLocked(t *testing.T) {
	ps := NewPeerServiceOnMemory()
	lock := NewLockOnMemory(ps, GetTestConfig())
	testLock_AddVoteMessageAndGetLocked(t, lock, ps)
}
