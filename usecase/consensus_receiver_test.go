package usecase_test

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/dba"
	"github.com/satellitex/bbft/model"
	. "github.com/satellitex/bbft/test_utils"
	. "github.com/satellitex/bbft/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func NewTestConsensusReceiverUsecase() (dba.ProposalTxQueue, dba.PeerService, dba.Lock, dba.BlockChain, model.ConsensusSender, *ReceiveChannel, ConsensusReceiver) {
	testConfig := GetTestConfig()
	queue := dba.NewProposalTxQueueOnMemory(testConfig)
	ps := dba.NewPeerServiceOnMemory()
	lock := dba.NewLockOnMemory(ps, testConfig)
	pool := dba.NewReceiverPoolOnMemory(testConfig)
	bc := dba.NewBlockChainOnMemory()
	slv := convertor.NewStatelessValidator()
	sender := convertor.NewMockConsensusSender()
	receivChan := NewReceiveChannel(testConfig)
	return queue, ps, lock, bc, sender, receivChan, NewConsensusReceiverUsecase(queue, ps, lock, pool, bc, slv, sender, receivChan)
}

func TestConsensusReceieverUsecase_Propagate(t *testing.T) {
	queue, _, _, _, sender, _, receiver := NewTestConsensusReceiverUsecase()
	t.Run("success case", func(t *testing.T) {
		tx := RandomValidTx(t)
		err := receiver.Propagate(tx)
		assert.NoError(t, err)
		assert.Equal(t, tx, sender.(*convertor.MockConsensusSender).Tx)
	})

	t.Run("failed case input nil", func(t *testing.T) {
		err := receiver.Propagate(nil)
		assert.EqualError(t, errors.Cause(err), model.ErrStatelessTxValidate.Error())
	})

	t.Run("failed case unverify tx", func(t *testing.T) {
		err := receiver.Propagate(RandomInvalidTx(t))
		assert.EqualError(t, errors.Cause(err), model.ErrStatelessTxValidate.Error())
	})

	t.Run("failed case already exist tx", func(t *testing.T) {
		tx := RandomValidTx(t)
		err := receiver.Propagate(tx)
		assert.NoError(t, err)
		err = receiver.Propagate(tx)
		assert.EqualError(t, errors.Cause(err), ErrAlradyReceivedSameObject.Error())
	})

	// To Empty queue
	for {
		_, ok := queue.Pop()
		if !ok {
			break
		}
	}

	t.Run("DoS safety test", func(t *testing.T) {
		waiter := &sync.WaitGroup{}
		for i := 0; i < GetTestConfig().QueueLimits; i++ {
			waiter.Add(1)
			go func() {
				err := receiver.Propagate(RandomValidTx(t))
				assert.NoError(t, err)
				waiter.Done()
			}()
		}
		waiter.Wait()
		for i := 0; i < 100; i++ {
			waiter.Add(1)
			go func() {
				err := receiver.Propagate(RandomValidTx(t))
				assert.EqualError(t, errors.Cause(err), dba.ErrProposalTxQueuePush.Error())
				waiter.Done()
			}()
		}
		waiter.Wait()
	})
}

func TestConsensusReceieverUsecase_Propose(t *testing.T) {
	_, ps, _, _, sender, channel, receiver := NewTestConsensusReceiverUsecase()

	peer := RandomPeerWithPriv()
	ps.AddPeer(peer)

	t.Run("success case", func(t *testing.T) {
		proposal := RandomProposalWithPeer(t, 0, 0, peer)
		err := receiver.Propose(proposal)
		require.NoError(t, err)
		assert.Equal(t, proposal, sender.(*convertor.MockConsensusSender).Proposal)
		assert.Equal(t, proposal, <-channel.Propose)
	})

	t.Run("failed case input nil", func(t *testing.T) {
		err := receiver.Propose(nil)
		assert.EqualError(t, errors.Cause(err), model.ErrInvalidProposal.Error())
	})

	t.Run("failed case input unverified propose", func(t *testing.T) {
		block := RandomInvalidBlock(t)
		proposal, err := convertor.NewModelFactory().NewProposal(block, 1)
		require.NoError(t, err)
		err = receiver.Propose(proposal)
		assert.EqualError(t, errors.Cause(err), model.ErrStatelessBlockValidate.Error())
	})

	t.Run("failed case not leader signed", func(t *testing.T) {
		proposal := RandomProposalWithHeightRound(t, 0, 0)
		err := receiver.Propose(proposal)
		assert.EqualError(t, errors.Cause(err), ErrVerifyOnlyLeader.Error())
	})

	t.Run("failed case already exist", func(t *testing.T) {
		proposal := RandomProposalWithPeer(t, 1, 0, peer)
		err := receiver.Propose(proposal)
		require.NoError(t, err)
		require.Equal(t, proposal, <-channel.Propose)

		err = receiver.Propose(proposal)
		assert.EqualError(t, errors.Cause(err), ErrAlradyReceivedSameObject.Error())
	})

	t.Run("DoS safety test", func(t *testing.T) {
		waiter := &sync.WaitGroup{}
		for i := 0; i < GetTestConfig().ReceiveProposeProposalPoolLimits*2; i++ {
			waiter.Add(1)
			go func(i int64) {
				err := receiver.Propose(RandomProposalWithPeer(t, i, 0, peer))
				assert.NoError(t, err)
				waiter.Done()
			}(int64(i+2))
			go func() {
				<-channel.Propose
			}()
		}
		waiter.Wait()
	})
}

func TestConsensusReceieverUsecase_Vote(t *testing.T) {
	_, ps, _, _, sender, channel, receiver := NewTestConsensusReceiverUsecase()
	peers := []model.Peer{
		RandomPeerWithPriv(),
		RandomPeerWithPriv(),
		RandomPeerWithPriv(),
		RandomPeerWithPriv(),
	}
	for _, p := range peers {
		ps.AddPeer(p)
	}

	t.Run("success case", func(t *testing.T) {
		vote := RandomVoteMessageFromPeer(t, peers[0])
		err := receiver.Vote(vote)
		assert.NoError(t, err)
		assert.Equal(t, vote, sender.(*convertor.MockConsensusSender).VoteMessage)
		assert.Equal(t, vote, <-channel.Vote)
	})

	t.Run("failed case input nil", func(t *testing.T) {
		err := receiver.Vote(nil)
		assert.EqualError(t, errors.Cause(err), model.ErrInvalidVoteMessage.Error())
	})

	t.Run("failed case input unverified vote", func(t *testing.T) {
		vote := RandomVoteMessage(t)
		vote.(*convertor.VoteMessage).Signature = nil
		err := receiver.Vote(vote)
		assert.EqualError(t, errors.Cause(err), model.ErrVoteMessageVerify.Error())
	})

	t.Run("failed case input not peers vote", func(t *testing.T) {
		vote := RandomVoteMessage(t)
		err := receiver.Vote(vote)
		assert.EqualError(t, errors.Cause(err), ErrVoteNotInPeerService.Error())
	})

	t.Run("fialed case already exist vote", func(t *testing.T) {
		vote := RandomVoteMessageFromPeer(t, peers[0])
		err := receiver.Vote(vote)
		require.NoError(t, err)
		require.Equal(t, vote, <-channel.Vote)

		err = receiver.Vote(vote)
		assert.EqualError(t, errors.Cause(err), ErrAlradyReceivedSameObject.Error())
	})

	t.Run("DoS safety test", func(t *testing.T) {
		waiter := &sync.WaitGroup{}
		for i := 0; i < GetTestConfig().ReceiveVoteVoteMessagePoolLimits*2; i++ {
			waiter.Add(1)
			go func() {
				err := receiver.Vote(RandomVoteMessageFromPeer(t, peers[1]))
				assert.NoError(t, err)
				waiter.Done()
			}()
			go func() {
				<-channel.Vote
			}()
		}
		waiter.Wait()
	})
}

func TestConsensusReceieverUsecase_PreCommit(t *testing.T) {
	_, ps, _, _, sender, channel, receiver := NewTestConsensusReceiverUsecase()
	peers := []model.Peer{
		RandomPeerWithPriv(),
		RandomPeerWithPriv(),
		RandomPeerWithPriv(),
		RandomPeerWithPriv(),
	}
	for _, p := range peers {
		ps.AddPeer(p)
	}

	t.Run("success case", func(t *testing.T) {
		preCommit := RandomVoteMessageFromPeer(t, peers[0])
		err := receiver.PreCommit(preCommit)
		assert.NoError(t, err)
		assert.Equal(t, preCommit, sender.(*convertor.MockConsensusSender).PreCommitMessage)
		assert.Equal(t, preCommit, <-channel.PreCommit)
	})

	t.Run("failed case input nil", func(t *testing.T) {
		err := receiver.PreCommit(nil)
		assert.EqualError(t, errors.Cause(err), model.ErrInvalidVoteMessage.Error())
	})

	t.Run("failed case input unverified preCommit", func(t *testing.T) {
		preCommit := RandomVoteMessage(t)
		preCommit.(*convertor.VoteMessage).Signature = nil
		err := receiver.PreCommit(preCommit)
		assert.EqualError(t, errors.Cause(err), model.ErrVoteMessageVerify.Error())
	})

	t.Run("failed case input not peers preCommit", func(t *testing.T) {
		preCommit := RandomVoteMessage(t)
		err := receiver.PreCommit(preCommit)
		assert.EqualError(t, errors.Cause(err), ErrPreCommitNotInPeerService.Error())
	})

	t.Run("fialed case already exist preCommit", func(t *testing.T) {
		preCommit := RandomVoteMessageFromPeer(t, peers[0])
		err := receiver.PreCommit(preCommit)
		require.NoError(t, err)
		require.Equal(t, preCommit, <-channel.PreCommit)

		err = receiver.PreCommit(preCommit)
		assert.EqualError(t, errors.Cause(err), ErrAlradyReceivedSameObject.Error())
	})

	t.Run("DoS safety test", func(t *testing.T) {
		waiter := &sync.WaitGroup{}
		for i := 0; i < GetTestConfig().ReceivePreCommitVoteMessagePoolLimits*2; i++ {
			waiter.Add(1)
			go func() {
				err := receiver.PreCommit(RandomVoteMessageFromPeer(t, peers[1]))
				assert.NoError(t, err)
				waiter.Done()
			}()
			go func() {
				<-channel.PreCommit
			}()
		}
		waiter.Wait()
	})
}
