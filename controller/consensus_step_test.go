package controller_test

import (
	"bytes"
	"context"
	"github.com/satellitex/bbft/config"
	. "github.com/satellitex/bbft/controller"
	"github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/dba"
	"github.com/satellitex/bbft/proto"
	. "github.com/satellitex/bbft/test_utils"
	"github.com/satellitex/bbft/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"testing"
)

func NewTestConsensusController(t *testing.T) (*config.BBFTConfig, dba.PeerService, *ConsensusController) {

	testConfig := GetTestConfig()
	queue := dba.NewProposalTxQueueOnMemory(testConfig)
	ps := RandomPeerService(t, 3)
	lock := dba.NewLockOnMemory(ps, testConfig)
	pool := dba.NewReceiverPoolOnMemory(testConfig)

	// bc is height = 1
	bc := dba.NewBlockChainOnMemory()
	bc.Commit(RandomCommitableBlock(t, bc))

	slv := convertor.NewStatelessValidator()
	sender := convertor.NewMockConsensusSender()
	receivChan := usecase.NewReceiveChannel(testConfig)
	receiver := usecase.NewConsensusReceiverUsecase(queue, ps, lock, pool, bc, slv, sender, receivChan)

	author := convertor.NewAuthor(ps)

	// add peer this peer
	ps.AddPeer(convertor.NewModelFactory().NewPeer(testConfig.Host, testConfig.PublicKey))

	return testConfig, ps, NewConsensusController(receiver, author)

}

func TestConsensusController_Propagate(t *testing.T) {

	conf, _, ctrl := NewTestConsensusController(t)

	validTx := RandomValidTx(t).(*convertor.Transaction).Transaction
	inValidTx := RandomInvalidTx(t).(*convertor.Transaction).Transaction

	evilConf := *conf
	pk, sk := convertor.NewKeyPair()
	evilConf.PublicKey = pk
	evilConf.SecretKey = sk

	for _, c := range []struct {
		name string
		ctx  context.Context
		tx   *bbft.Transaction
		code codes.Code
	}{
		{
			"success case",
			ValidContext(t, conf, validTx),
			validTx,
			codes.OK,
		},
		{
			"failed case, unauthenticated context",
			context.TODO(),
			validTx,
			codes.Unauthenticated,
		},
		{
			"failed case, authenticated but not peer",
			ValidContext(t, &evilConf, validTx),
			validTx,
			codes.PermissionDenied,
		},
		{
			"failed case, invalid transaction",
			ValidContext(t, conf, inValidTx),
			inValidTx,
			codes.InvalidArgument,
		},
		{
			"failed case, unauthenticated",
			ValidContext(t, conf, inValidTx),
			nil,
			codes.Unauthenticated,
		},
		{
			"failed case, duplicate sending tx",
			ValidContext(t, conf, validTx),
			validTx,
			codes.AlreadyExists,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			_, err := ctrl.Propagate(c.ctx, c.tx)
			if c.code != codes.OK {
				ValidateStatusCode(t, err, c.code)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func TestConsensusController_Propose(t *testing.T) {

	conf, ps, ctrl := NewTestConsensusController(t)

	leaderId := func() int32 {
		for i, p := range ps.GetPermutationPeers(1) {
			if bytes.Equal(conf.PublicKey, p.GetPubkey()) {
				return int32(i)
			}
		}
		return -1
	}()
	require.NotEqual(t, -1, leaderId)

	validProposal := RandomProposalWithHeightRound(t, 1, leaderId).(*convertor.Proposal).Proposal
	unLeaderProposal := RandomProposalWithHeightRound(t, 1, (leaderId+1)%2).(*convertor.Proposal).Proposal
	invalidProposal := RandomInvalidProposal(t).(*convertor.Proposal).Proposal

	evilConf := *conf
	pk, sk := convertor.NewKeyPair()
	evilConf.PublicKey = pk
	evilConf.SecretKey = sk

	for _, c := range []struct {
		name     string
		ctx      context.Context
		proposal *bbft.Proposal
		code     codes.Code
	}{
		{
			"success case",
			ValidContext(t, conf, validProposal),
			validProposal,
			codes.OK,
		},
		{
			"failed case, unauthenticated context",
			context.TODO(),
			validProposal,
			codes.Unauthenticated,
		},
		{
			"failed case, authenticated but not peer",
			ValidContext(t, &evilConf, validProposal),
			validProposal,
			codes.PermissionDenied,
		},
		{
			"failed case, authenticated and peer but not leader",
			ValidContext(t, conf, unLeaderProposal),
			unLeaderProposal,
			codes.PermissionDenied,
		},
		{
			"failed case, invalid Proposal",
			ValidContext(t, conf, invalidProposal),
			invalidProposal,
			codes.InvalidArgument,
		},
		{
			"failed case, nil",
			ValidContext(t, conf, invalidProposal),
			nil,
			codes.Unauthenticated,
		},
		{
			"failed case, duplicate sent",
			ValidContext(t, conf, validProposal),
			validProposal,
			codes.AlreadyExists,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			_, err := ctrl.Propose(c.ctx, c.proposal)
			if c.code != codes.OK {
				ValidateStatusCode(t, err, c.code)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func TestConsensusController_Vote(t *testing.T) {

	conf, ps, ctrl := NewTestConsensusController(t)

	validVote := RandomVoteMessageFromPeer(t, ps.GetPeers()[0]).(*convertor.VoteMessage).VoteMessage
	unPeerValidVote := RandomVoteMessage(t).(*convertor.VoteMessage).VoteMessage

	evilConf := *conf
	pk, sk := convertor.NewKeyPair()
	evilConf.PublicKey = pk
	evilConf.SecretKey = sk

	for _, c := range []struct {
		name string
		ctx  context.Context
		vote *bbft.VoteMessage
		code codes.Code
	}{
		{
			"success case",
			ValidContext(t, conf, validVote),
			validVote,
			codes.OK,
		},
		{
			"failed case, unauthenticated context",
			context.TODO(),
			validVote,
			codes.Unauthenticated,
		},
		{
			"failed case, authenticated but not peer",
			ValidContext(t, &evilConf, validVote),
			validVote,
			codes.PermissionDenied,
		},
		{
			"failed case, unsigned vote",
			ValidContext(t, conf, unPeerValidVote),
			unPeerValidVote,
			codes.InvalidArgument,
		},
		{
			"failed case, nil",
			ValidContext(t, conf, validVote),
			nil,
			codes.Unauthenticated,
		},
		{
			"failed case, duplicate sent",
			ValidContext(t, conf, validVote),
			validVote,
			codes.AlreadyExists,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			_, err := ctrl.Vote(c.ctx, c.vote)
			if c.code != codes.OK {
				ValidateStatusCode(t, err, c.code)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func TestConsensusController_PreCommit(t *testing.T) {

	conf, ps, ctrl := NewTestConsensusController(t)

	validVote := RandomVoteMessageFromPeer(t, ps.GetPeers()[0]).(*convertor.VoteMessage).VoteMessage
	unPeerValidVote := RandomVoteMessage(t).(*convertor.VoteMessage).VoteMessage

	evilConf := *conf
	pk, sk := convertor.NewKeyPair()
	evilConf.PublicKey = pk
	evilConf.SecretKey = sk

	for _, c := range []struct {
		name string
		ctx  context.Context
		vote *bbft.VoteMessage
		code codes.Code
	}{
		{
			"success case",
			ValidContext(t, conf, validVote),
			validVote,
			codes.OK,
		},
		{
			"failed case, unauthenticated context",
			context.TODO(),
			validVote,
			codes.Unauthenticated,
		},
		{
			"failed case, authenticated but not peer",
			ValidContext(t, &evilConf, validVote),
			validVote,
			codes.PermissionDenied,
		},
		{
			"failed case, unsigned vote",
			ValidContext(t, conf, unPeerValidVote),
			unPeerValidVote,
			codes.InvalidArgument,
		},
		{
			"failed case, nil",
			ValidContext(t, conf, validVote),
			nil,
			codes.Unauthenticated,
		},
		{
			"failed case, duplicate sent",
			ValidContext(t, conf, validVote),
			validVote,
			codes.AlreadyExists,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			_, err := ctrl.PreCommit(c.ctx, c.vote)
			if c.code != codes.OK {
				ValidateStatusCode(t, err, c.code)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
