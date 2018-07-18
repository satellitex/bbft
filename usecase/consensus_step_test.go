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
	"sync"
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

func NewTestConsensusStepUsecase(t *testing.T) (*config.BBFTConfig, dba.BlockChain, dba.PeerService, dba.Lock,
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

	//First Commit
	bc.Commit(RandomCommitableBlock(t, bc))

	ps.AddPeer(&PeerWithPriv{
		&convertor.Peer{"myself", conf.PublicKey},
		conf.SecretKey,
	})
	ps.AddPeer(RandomPeerWithPriv())
	ps.AddPeer(RandomPeerWithPriv())
	ps.AddPeer(RandomPeerWithPriv())

	consensusStep := NewConsensusStepUsecase(conf, bc, ps, lock, queue, sender, slv, sfv, factory, channel)
	return conf, bc, ps, lock, queue, sender, channel, consensusStep
}

func mySelfId(conf *config.BBFTConfig, ps dba.PeerService, height int64) int32 {
	for id, p := range ps.GetPermutationPeers(height) {
		if bytes.Equal(p.GetPubkey(), conf.PublicKey) {
			return int32(id)
		}
	}
	return 0
}

func TestConsensusStepUsecase_Propose(t *testing.T) {
	conf, bc, ps, lock, queue, sender, channel, c := NewTestConsensusStepUsecase(t)
	factory := convertor.NewModelFactory()

	top, ok := bc.Top()
	require.True(t, ok)

	var height int64 = 1

	// myselfId is Round when myself is Leader
	myselfId := mySelfId(conf, ps, height)

	t.Run("leader case", func(t *testing.T) {
		validTx := RandomValidTx(t)
		invalidTx := RandomInvalidTx(t)
		queue.Push(validTx)
		queue.Push(invalidTx)

		// set CommitTime
		c.(*ConsensusStepUsecase).RoundCommitTime = time.Duration(Now())
		// not leader
		err := c.Propose(height, myselfId)
		assert.NoError(t, err)

		tmp, err := factory.NewBlock(height, GetHash(t, top),
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
		err := c.Propose(height, int32((myselfId+1)%4))
		endTime := Now()

		assert.NoError(t, err)
		assert.True(t, TimeParseDuration(t, "200ms") < time.Duration(endTime-startTime))
		assert.True(t, TimeParseDuration(t, "210ms") > time.Duration(endTime-startTime))
	})

	t.Run("not leader get proposal case", func(t *testing.T) {
		c.(*ConsensusStepUsecase).ProposeTimeOut = time.Duration(Now()) + conf.ProposeMaxCalcTime + conf.AllowedConnectDelayTime

		expectedProposal := RandomProposalWithHeightRound(t, height, int32((myselfId+2)%4))
		go func() {
			err := c.Propose(height, int32((myselfId+2)%4))
			assert.NoError(t, err)
			assert.Equal(t, expectedProposal, c.(*ConsensusStepUsecase).ThisRoundProposal)
		}()
		channel.Propose <- expectedProposal
	})

	t.Run("not leader get vote locked case", func(t *testing.T) {
		c.(*ConsensusStepUsecase).ProposeTimeOut = time.Duration(Now()) + conf.ProposeMaxCalcTime + conf.AllowedConnectDelayTime

		expectedProposal := RandomProposalWithHeightRound(t, height, int32((myselfId+2)%4))
		lock.RegisterProposal(expectedProposal)

		go func() {
			err := c.Propose(height, int32((myselfId+3)%4))
			assert.NoError(t, err)
			actual, ok := lock.GetLockedProposal(height)
			require.True(t, ok)
			assert.Equal(t, expectedProposal, actual)
		}()
		for _, p := range ps.GetPeers()[1:] {
			vote := RandomVoteMessageFromPeerWithBlock(t, p, expectedProposal.GetBlock())
			require.NoError(t, lock.AddVoteMessage(vote))
			channel.Vote <- vote
		}
		actualProposal, ok := lock.GetLockedProposal(height)
		require.True(t, ok)
		require.Equal(t, expectedProposal, actualProposal)
	})
}

func TestConsensusStepUsecase_Vote(t *testing.T) {
	conf, bc, ps, lock, _, sender, channel, c := NewTestConsensusStepUsecase(t)
	factory := convertor.NewModelFactory()

	_, ok := bc.Top()
	require.True(t, ok)

	var height int64 = 1

	// myselfId is Round when myself is Leader
	myselfId := mySelfId(conf, ps, height)

	t.Run("normal case, voteTimeOut", func(t *testing.T) {
		c.(*ConsensusStepUsecase).VoteTimeOut = time.Duration(Now()) + TimeParseDuration(t, "200ms")

		startTime := Now()
		err := c.Vote(height, 0)
		endTime := Now()
		assert.NoError(t, err)
		assert.True(t, TimeParseDuration(t, "200ms") < time.Duration(endTime-startTime), "%v", time.Duration(endTime-startTime))
		assert.True(t, TimeParseDuration(t, "210ms") > time.Duration(endTime-startTime), "%v", time.Duration(endTime-startTime))
	})

	t.Run("normal case, validate proposal and sendVote", func(t *testing.T) {
		waiter := &sync.WaitGroup{}
		waiter.Add(1)
		validProposal, err := factory.NewProposal(RandomCommitableBlock(t, bc), 0)
		require.NoError(t, err)

		c.(*ConsensusStepUsecase).ThisRoundProposal = validProposal
		c.(*ConsensusStepUsecase).VoteTimeOut = time.Duration(Now()) + conf.VoteMaxCalcTime + conf.AllowedConnectDelayTime
		lock.RegisterProposal(validProposal)

		go func() {
			err := c.Vote(height, 0)
			assert.NoError(t, err)

			actual, ok := lock.GetLockedProposal(height)
			require.True(t, ok)
			assert.Equal(t, validProposal, actual)

			actualVote := sender.(*convertor.MockConsensusSender).VoteMessage
			expectedVote := RandomVoteMessageFromPeerWithBlock(t, ps.GetPermutationPeers(height)[myselfId], validProposal.GetBlock())
			assert.Equal(t, expectedVote, actualVote)
			waiter.Done()
		}()
		for _, p := range ps.GetPeers()[1:] {
			vote := RandomVoteMessageFromPeerWithBlock(t, p, validProposal.GetBlock())
			require.NoError(t, lock.AddVoteMessage(vote))
			channel.Vote <- vote
		}

		actualProposal, ok := lock.GetLockedProposal(height)
		require.True(t, ok)
		require.Equal(t, validProposal, actualProposal)
		waiter.Wait()
	})

	t.Run("normal case, if has lock, no wait voteTimeOut", func(t *testing.T) {
		c.(*ConsensusStepUsecase).VoteTimeOut = time.Duration(Now()) + TimeParseDuration(t, "200ms")

		startTime := Now()
		err := c.Vote(height, 0)
		endTime := Now()
		assert.NoError(t, err)
		assert.True(t, TimeParseDuration(t, "190ms") > time.Duration(endTime-startTime), "%v", time.Duration(endTime-startTime))
	})

}

func TestConsensusStepUsecase_PreCommit(t *testing.T) {
	conf, bc, ps, lock, _, sender, channel, c := NewTestConsensusStepUsecase(t)
	factory := convertor.NewModelFactory()

	_, ok := bc.Top()
	require.True(t, ok)

	var height int64 = 1
	// myselfId is Round when myself is Leader
	myselfId := mySelfId(conf, ps, height)

	t.Run("invalid case, preCommitTimeOut", func(t *testing.T) {
		c.(*ConsensusStepUsecase).PreCommitTimeOut = time.Duration(Now()) + TimeParseDuration(t, "200ms")

		startTime := Now()
		err := c.PreCommit(height, 0)
		endTime := Now()
		assert.EqualError(t, errors.Cause(err), ErrConsensusPreCommit.Error())
		assert.True(t, TimeParseDuration(t, "200ms") < time.Duration(endTime-startTime), "%v", time.Duration(endTime-startTime))
		assert.True(t, TimeParseDuration(t, "210ms") > time.Duration(endTime-startTime), "%v", time.Duration(endTime-startTime))
	})

	t.Run("normal case, sendPreCommit and collected preCommit", func(t *testing.T) {
		proposal, err := factory.NewProposal(RandomCommitableBlock(t, bc), 0)
		require.NoError(t, err)
		lock.RegisterProposal(proposal)

		for _, p := range ps.GetPeers()[1:] {
			vote := RandomVoteMessageFromPeerWithBlock(t, p, proposal.GetBlock())
			require.NoError(t, lock.AddVoteMessage(vote))
		}
		locked, ok := lock.GetLockedProposal(height)
		require.True(t, ok)
		require.Equal(t, proposal, locked)

		c.(*ConsensusStepUsecase).PreCommitTimeOut = time.Duration(Now()) + conf.PreCommitMaxCalcTime + conf.AllowedConnectDelayTime

		go func() {
			err := c.PreCommit(height, 0)
			assert.NoError(t, err)

			expectedPreCommit := RandomVoteMessageFromPeerWithBlock(t, ps.GetPermutationPeers(height)[myselfId], proposal.GetBlock())
			actualPreCommit := sender.(*convertor.MockConsensusSender).PreCommitMessage
			assert.Equal(t, expectedPreCommit, actualPreCommit)
		}()
		for _, p := range ps.GetPeers()[1:] {
			vote := RandomVoteMessageFromPeerWithBlock(t, p, proposal.GetBlock())
			channel.PreCommit <- vote
		}
	})
}

func TestConsensusStepUsecase_Commit(t *testing.T) {
	_, bc, ps, lock, _, _, _, c := NewTestConsensusStepUsecase(t)
	factory := convertor.NewModelFactory()

	_, ok := bc.Top()
	require.True(t, ok)

	var height int64 = 1

	t.Run("success, commit!", func(t *testing.T) {
		proposal, err := factory.NewProposal(RandomCommitableBlock(t, bc), 0)
		require.NoError(t, err)

		lock.RegisterProposal(proposal)
		for _, p := range ps.GetPeers()[1:] {
			vote := RandomVoteMessageFromPeerWithBlock(t, p, proposal.GetBlock())
			require.NoError(t, lock.AddVoteMessage(vote))
		}

		assert.NoError(t, c.Commit(height, 0))
		newBlock, ok := bc.Top()
		require.True(t, ok)
		assert.Equal(t, proposal.GetBlock(), newBlock)
	})

	t.Run("invalid commit case", func(t *testing.T) {
		assert.Error(t, errors.Cause(c.Commit(height, 0)), ErrConsensusCommit.Error())
	})
}
