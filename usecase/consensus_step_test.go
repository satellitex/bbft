package usecase_test

import (
	"github.com/satellitex/bbft/dba"
	"testing"

	"bytes"
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/config"
	"github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/model"
	. "github.com/satellitex/bbft/test_utils"
	. "github.com/satellitex/bbft/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"time"
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

	t.Run("success2", func(t *testing.T) {
		proposal := RandomProposalWithHeightRound(t, 100, 3)
		err := finder.Set(proposal)
		assert.NoError(t, err)
		actual, ok := finder.Find(100, 3)
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

func NewTestConsensusStepUsecase() (*config.BBFTConfig, dba.BlockChain, dba.PeerService, dba.Lock,
	dba.ProposalTxQueue, model.ConsensusSender, *ReceiveChannel, ConsensusStep) {

	conf := GetTestConfig()
	bc := dba.NewBlockChainOnMemory()
	ps := dba.NewPeerServiceOnMemory()
	lock := dba.NewLockOnMemory(ps, conf)
	queue := dba.NewProposalTxQueueOnMemory(conf)
	sender := convertor.NewMockConsensusSender()
	slv := convertor.NewStatelessValidator()
	sfv := convertor.NewStatefulValidator(bc)
	factory := convertor.NewModelFactory()
	channel := NewReceiveChannel(conf)

	consensusStep := NewConsensusStepUsecase(conf, bc, ps, lock, queue, sender, slv, sfv, factory, channel)
	return conf, bc, ps, lock, queue, sender, channel, consensusStep
}

func TestConsensusStepUsecase_Propose(t *testing.T) {
	conf, bc, ps, lock, queue, sender, channel, c := NewTestConsensusStepUsecase()
	factory := convertor.NewModelFactory()

	//First Commit
	bc.Commit(RandomCommitableBlock(t, bc))
	top, ok := bc.Top()
	require.True(t, ok)

	ps.AddPeer(&PeerWithPriv{
		&convertor.Peer{"myself", conf.PublicKey},
		conf.SecretKey,
	})
	ps.AddPeer(RandomPeerWithPriv())
	ps.AddPeer(RandomPeerWithPriv())
	ps.AddPeer(RandomPeerWithPriv())

	// myselfId is Round when myself is Leader
	myselfId := func() int32 {
		for id, p := range ps.GetPermutationPeers(0) {
			if bytes.Equal(p.GetPubkey(), conf.PublicKey) {
				return int32(id)
			}
		}
		return 0
	}()

	t.Run("leader case", func(t *testing.T) {
		validTx := RandomValidTx(t)
		invalidTx := RandomInvalidTx(t)
		queue.Push(validTx)
		queue.Push(invalidTx)

		// set CommitTime
		c.(*ConsensusStepUsecase).RoundCommitTime = time.Duration(Now())
		// not leader
		err := c.Propose(0, myselfId)
		assert.NoError(t, err)

		tmp, err := factory.NewBlock(0, GetHash(t, top),
			int64(c.(*ConsensusStepUsecase).RoundCommitTime), []model.Transaction{validTx})
		tmp.Sign(conf.PublicKey, conf.SecretKey)
		require.NoError(t, err)
		expectedProposal, err := factory.NewProposal(tmp, myselfId)
		require.NoError(t, err)

		assert.Equal(t, expectedProposal, sender.(*convertor.MockConsensusSender).Proposal)
		assert.Equal(t, expectedProposal, c.(*ConsensusStepUsecase).ThisRoundProposal)
	})

	t.Run("timeOut case not leader ", func(t *testing.T) {
		startTime := Now()
		c.(*ConsensusStepUsecase).ProposeTimeOut = time.Duration(Now()) + TimeParseDuration(t, "200ms")
		// not leader
		err := c.Propose(0, int32((myselfId+1)%4))
		endTime := Now()

		assert.NoError(t, err)
		assert.True(t, TimeParseDuration(t, "200ms") < time.Duration(endTime-startTime))
		assert.True(t, TimeParseDuration(t, "210ms") > time.Duration(endTime-startTime))
	})

	t.Run("not leader get proposal case", func(t *testing.T) {
		c.(*ConsensusStepUsecase).ProposeTimeOut = time.Duration(Now()) + conf.ProposeMaxCalcTime + conf.AllowedConnectDelayTime

		expectedProposal := RandomProposalWithHeightRound(t, 0, int32((myselfId+2)%3))
		go func() {
			err := c.Propose(0, int32((myselfId+2)%4))
			assert.NoError(t, err)
			assert.Equal(t, expectedProposal, c.(*ConsensusStepUsecase).ThisRoundProposal)
		}()
		channel.Propose <- expectedProposal
	})

	t.Run("not leader get vote locked case", func(t *testing.T) {
		c.(*ConsensusStepUsecase).ProposeTimeOut = time.Duration(Now()) + conf.ProposeMaxCalcTime + conf.AllowedConnectDelayTime

		expectedProposal := RandomProposalWithHeightRound(t, 0, int32((myselfId+2)%3))
		lock.RegisterProposal(expectedProposal)

		go func() {
			err := c.Propose(0, int32((myselfId+3)%4))
			assert.NoError(t, err)
			actual, ok := lock.GetLockedProposal(0)
			require.True(t, ok)
			assert.Equal(t, expectedProposal, actual)
		}()
		for _, p := range ps.GetPeers()[1:] {
			vote := RandomVoteMessageFromPeerWithBlock(t, p, expectedProposal.GetBlock())
			require.NoError(t, lock.AddVoteMessage(vote))
			channel.Vote <- vote
		}
		actualProposal, ok := lock.GetLockedProposal(0)
		require.True(t, ok)
		require.Equal(t, expectedProposal, actualProposal)
	})

}
