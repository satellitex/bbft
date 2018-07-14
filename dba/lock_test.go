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
}

func testLock_AddVoteMessageAndGetLocked(t *testing.T, lock Lock, p PeerService) {
	// 4 peer
	peers := []model.Peer{
		RandomPeerWithPriv(),
		RandomPeerWithPriv(),
		RandomPeerWithPriv(),
		RandomPeerWithPriv(),
	}
	for _, peer := range peers {
		p.AddPeer(peer)
	}

	validGetLockedProposal := func(t *testing.T, expectedProposal model.Proposal) {
		proposal, ok := lock.GetLockedProposal()
		if expectedProposal == nil {
			assert.False(t, ok)
		} else {
			assert.True(t, ok)
			assert.Equal(t, expectedProposal, proposal)
		}
	}

	t.Run("success valid votes", func(t *testing.T) {
		vote := RandomVoteMesssageFromPeer(t, peers[0])
		in, err := lock.AddVoteMessage(vote)
		assert.NoError(t, err)
		assert.False(t, in)

		validGetLockedProposal(t, nil)
	})
	t.Run("faield nil vote", func(t *testing.T) {
		in, err := lock.AddVoteMessage(nil)
		assert.EqualError(t, errors.Cause(err), model.ErrInvalidVoteMessage.Error())
		assert.False(t, in)

		validGetLockedProposal(t, nil)
	})


	// register valid proposal
	validProposals := []model.Proposal{
		RandomProposal(t),
		RandomProposal(t),
	}
	if validProposals[0].GetRound() > validProposals[1].GetRound() {
		validProposals[0], validProposals[1] = validProposals[1], validProposals[0]
	}
	lock.RegisterProposal(validProposals[0])
	lock.RegisterProposal(validProposals[1])
	require.True(t, validProposals[0].GetRound() < validProposals[1].GetRound())

	validAddVote := func(t *testing.T, proposal model.Proposal) {
		vote := convertor.NewModelFactory().NewVoteMessage(GetHash(t, proposal.GetBlock()))
		ValidSign(t, vote)
		_, err := lock.AddVoteMessage(vote)
		require.NoError(t, err)
	}

	t.Run("success valid vote setLockedProposal", func(t *testing.T) {
		vp := validProposals[0]
		validAddVote(t, vp)
		validAddVote(t, vp)
		validGetLockedProposal(t, nil)

		vote := convertor.NewModelFactory().NewVoteMessage(GetHash(t, vp.GetBlock()))
		ValidSign(t, vote)
		in, err := lock.AddVoteMessage(vote)
		assert.NoError(t, err)
		assert.True(t, in)

		validGetLockedProposal(t, vp)

		vote = convertor.NewModelFactory().NewVoteMessage(GetHash(t, vp.GetBlock()))
		ValidSign(t, vote)
		in, err = lock.AddVoteMessage(vote)
		assert.NoError(t, err)
		assert.False(t, in)

		validGetLockedProposal(t, vp)
	})

	t.Run("success valid second vote setLockedProposal", func(t *testing.T) {
		vp := validProposals[1]
		validAddVote(t, vp)
		validAddVote(t, vp)
		validGetLockedProposal(t, validProposals[0])

		vote := convertor.NewModelFactory().NewVoteMessage(GetHash(t, vp.GetBlock()))
		ValidSign(t, vote)
		in, err := lock.AddVoteMessage(vote)
		assert.NoError(t, err)
		assert.True(t, in)

		validGetLockedProposal(t, vp)
	})

	t.Run("failed unregisterd Proposal votes, setLockedProposal", func(t *testing.T) {
		vp := RandomProposal(t)
		validAddVote(t, vp)
		validAddVote(t, vp)
		validGetLockedProposal(t, validProposals[1])

		vote := convertor.NewModelFactory().NewVoteMessage(GetHash(t, vp.GetBlock()))
		ValidSign(t, vote)
		in, err := lock.AddVoteMessage(vote)
		assert.EqualError(t, errors.Cause(err), ErrSetLockedProposal.Error())
		assert.False(t, in)

		validGetLockedProposal(t, validProposals[1])
	})
}

func TestLockOnMemory_RegisterProposal(t *testing.T) {
	lock := NewLockOnMemory(NewPeerServiceOnMemory())
	testLock_RegisterProposal(t, lock)
}

func TestLockOnMemory_AddVoteMessageAndGetLocked(t *testing.T) {
	ps := NewPeerServiceOnMemory()
	lock := NewLockOnMemory(ps)
	testLock_AddVoteMessageAndGetLocked(t, lock, ps)
}
