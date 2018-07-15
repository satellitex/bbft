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

func NewTestConsensusReceiverUsecase() (dba.ProposalTxQueue, dba.PeerService, dba.Lock, dba.BlockChain, model.ConsensusSender, ConsensusReceiver) {
	testConfig := GetTestConfig()
	queue := dba.NewProposalTxQueueOnMemory(testConfig)
	ps := dba.NewPeerServiceOnMemory()
	lock := dba.NewLockOnMemory(ps, testConfig)
	pool := dba.NewReceiverPoolOnMemory(testConfig)
	bc := dba.NewBlockChainOnMemory()
	slv := convertor.NewStatelessValidator()
	sender := convertor.NewMockConsensusSender()
	return queue, ps, lock, bc, sender, NewConsensusReceiverUsecase(queue, lock, pool, bc, slv, sender)
}

func TestConsensusReceieverUsecase_Propagate(t *testing.T) {
	queue, _, _, _, sender, receiver := NewTestConsensusReceiverUsecase()
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
	_, _, _, _, sender, receiver := NewTestConsensusReceiverUsecase()
	t.Run("success case", func(t *testing.T) {
		proposal := RandomProposalWithHeightRound(t, 0, 0)
		err := receiver.Propose(proposal)
		assert.NoError(t, err)
		assert.Equal(t, proposal, sender.(*convertor.MockConsensusSender).Proposal)
	})

	t.Run("failed case input nil", func(t *testing.T) {
		err := receiver.Propose(nil)
		assert.EqualError(t, errors.Cause(err), model.ErrInvalidProposal.Error())
	})

	t.Run("failed case input propose", func(t *testing.T) {
		block := RandomInvalidBlock(t)
		proposal, err := convertor.NewModelFactory().NewProposal(block, 1)
		require.NoError(t, err)
		err = receiver.Propose(proposal)
		assert.EqualError(t, errors.Cause(err), model.ErrStatelessBlockValidate.Error())
	})

	t.Run("failed case already exist", func(t *testing.T) {
		proposal := RandomProposal(t)
		err := receiver.Propose(proposal)
		require.NoError(t, err)

		err = receiver.Propose(proposal)
		assert.EqualError(t, errors.Cause(err), ErrAlradyReceivedSameObject.Error())
	})

	t.Run("DoS safety test", func(t *testing.T) {
		waiter := &sync.WaitGroup{}
		for i := 0; i < GetTestConfig().ReceiveProposeProposalPoolLimits * 2; i++ {
			waiter.Add(1)
			go func() {
				err := receiver.Propose(RandomProposal(t))
				assert.NoError(t, err)
				waiter.Done()
			}()
		}
		waiter.Wait()
	})
}

func TestConsensusReceieverUsecase_Vote(t *testing.T) {
	_, ps, _, _, sender, receiver := NewTestConsensusReceiverUsecase()
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
		vote := RandomVoteMesssageFromPeer(t, peers[0])
		err := receiver.Vote(vote)
		assert.NoError(t, err)
		assert.Equal(t, vote, sender.(*convertor.MockConsensusSender).VoteMessage)
	})

	t.Run("failed case input nil", func(t *testing.T) {
		err := receiver.Vote(nil)
		assert.EqualError(t, errors.Cause(err), model.ErrInvalidVoteMessage.Error())
	})

}

func TestConsensusReceieverUsecase_PreCommit(t *testing.T) {

}
