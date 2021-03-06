package grpc_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/config"
	"github.com/satellitex/bbft/controller"
	"github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/dba"
	. "github.com/satellitex/bbft/grpc"
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
	. "github.com/satellitex/bbft/test_utils"
	"github.com/satellitex/bbft/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"net"
	"sync"
	"testing"
)

func SetUpTestServer(t *testing.T, conf *config.BBFTConfig, ps dba.PeerService, s *grpc.Server) {
	fmt.Println(conf.Port)

	l, err := net.Listen("tcp", ":"+conf.Port)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Succcess New Listen")

	author := convertor.NewAuthor(ps)

	queue := dba.NewProposalTxQueueOnMemory(conf)
	lock := dba.NewLockOnMemory(ps, conf)
	pool := dba.NewReceiverPoolOnMemory(conf)

	bc := dba.NewBlockChainOnMemory()

	slv := convertor.NewStatelessValidator()
	sender := convertor.NewMockConsensusSender() // WIP
	receivChan := usecase.NewReceiveChannel(conf)

	consensusReceiver := usecase.NewConsensusReceiverUsecase(queue, ps, lock, pool, bc, slv, sender, receivChan)
	clientRceiver := usecase.NewClientGateReceiverUsecase(slv, sender)
	fmt.Println("Success New Receivers")

	bbft.RegisterConsensusGateServer(s, controller.NewConsensusController(consensusReceiver, author))
	bbft.RegisterTxGateServer(s, controller.NewClientGateController(clientRceiver, author))
	fmt.Println("Success New Register Endpoint")

	if err := s.Serve(l); err != nil {
		fmt.Printf("Failed to server grpc: %s", err.Error())
	}
}

func NewTestGrpcServer() *grpc.Server {
	return grpc.NewServer([]grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_validator.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	}...)
}

func TestTxGateWrite(t *testing.T) {
	conf := GetTestConfig()
	ps := dba.NewPeerServiceOnMemory()
	server := NewTestGrpcServer()

	go func() {
		SetUpTestServer(t, conf, ps, server)
	}()

	sender := NewTxGateSender(t, conf)

	tx := RandomValidTx(t)
	err := sender.Write(context.TODO(), tx)
	assert.NoError(t, err)

	tx = RandomInvalidTx(t)
	err = sender.Write(context.TODO(), tx)
	ValidateStatusCode(t, err, codes.InvalidArgument)

	waiter := &sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		waiter.Add(1)
		go func() {
			err := sender.Write(context.TODO(), RandomValidTx(t))
			assert.NoError(t, err)
			waiter.Done()
		}()
	}
	waiter.Wait()

	server.GracefulStop()
}

func TestGrpcConsensusSender_Propagate(t *testing.T) {
	conf := GetTestConfig()
	ps := dba.NewPeerServiceOnMemory()
	ps.AddPeer(RandomPeerFromConf(conf)) // Just One Peer

	server := NewTestGrpcServer()

	go func() {
		SetUpTestServer(t, conf, ps, server)
	}()

	validTx := RandomValidTx(t)
	inValidTx := RandomInvalidTx(t)

	evilConf := *conf
	pk, sk := convertor.NewKeyPair()
	evilConf.PublicKey = pk
	evilConf.SecretKey = sk

	sender := NewGrpcConsensusSender(conf, ps)
	evilSender := NewGrpcConsensusSender(&evilConf, ps)

	for _, c := range []struct {
		name   string
		tx     model.Transaction
		sender model.ConsensusSender
		code   codes.Code
		err    error
	}{
		{
			"success case",
			validTx,
			sender,
			codes.OK,
			nil,
		},
		{
			"failed case, authenticated but not peer",
			validTx,
			evilSender,
			codes.PermissionDenied,
			nil,
		},
		{
			"failed case, invalid transaction",
			inValidTx,
			sender,
			codes.InvalidArgument,
			nil,
		},
		{
			"failed case, nil",
			nil,
			sender,
			codes.Unauthenticated,
			model.ErrInvalidTransaction,
		},
		{
			"failed case, duplicate sending tx",
			validTx,
			sender,
			codes.AlreadyExists,
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			err := c.sender.Propagate(c.tx)
			if c.err != nil {
				assert.EqualError(t, errors.Cause(err), c.err.Error())
			} else if c.code != codes.OK {
				ValidateStatusCode(t, err, c.code)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	server.GracefulStop()
}

func TestGrpcConsensusSender_Propose(t *testing.T) {
	confs := []*config.BBFTConfig{
		GetTestConfig(),
		GetTestConfig(),
		GetTestConfig(),
		GetTestConfig(),
	}
	confs[0].Port = "50053"
	confs[1].Port = "50054"
	confs[2].Port = "50055"
	confs[3].Port = "50056"

	ps := dba.NewPeerServiceOnMemory()
	for _, conf := range confs {
		ps.AddPeer(RandomPeerFromConf(conf))
	}

	servers := make([]*grpc.Server, 0, 4)
	for i, conf := range confs {
		servers = append(servers, NewTestGrpcServer())
		go func(conf *config.BBFTConfig, server *grpc.Server) {
			SetUpTestServer(t, conf, ps, server)
		}(conf, servers[i])
	}

	leaderId := func() int32 {
		for i, p := range ps.GetPermutationPeers(1) {
			if bytes.Equal(confs[0].PublicKey, p.GetPubkey()) {
				return int32(i)
			}
		}
		return -1
	}()
	require.NotEqual(t, -1, leaderId)

	leader := ps.GetPermutationPeers(1)[leaderId]
	validProposal := RandomProposalWithPeer(t, 1, leaderId, leader)
	unLeaderSignedProposal := RandomProposalWithHeightRound(t, 1, leaderId)
	invalidProposal := RandomInvalidProposalWithRound(t, 1, leaderId)

	evilConf := *confs[0]
	pk, sk := convertor.NewKeyPair()
	evilConf.PublicKey = pk
	evilConf.SecretKey = sk

	sender := NewGrpcConsensusSender(confs[0], ps)
	evilSender := NewGrpcConsensusSender(&evilConf, ps)

	for _, c := range []struct {
		name     string
		proposal model.Proposal
		sender   model.ConsensusSender
		code     codes.Code
		err      error
	}{
		{
			"success case",
			validProposal,
			sender,
			codes.OK,
			nil,
		},
		{
			"failed case, authenticated but not peer",
			validProposal,
			evilSender,
			codes.PermissionDenied,
			nil,
		},
		{
			"failed case, authenticated but not leader proposal",
			unLeaderSignedProposal,
			sender,
			codes.InvalidArgument,
			nil,
		},
		{
			"failed case, invalid proposal",
			invalidProposal,
			sender,
			codes.InvalidArgument,
			nil,
		},
		{
			"failed case, nil",
			nil,
			sender,
			codes.Unauthenticated,
			model.ErrInvalidProposal,
		},
		{
			"failed case, duplicate sent",
			validProposal,
			sender,
			codes.AlreadyExists,
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			err := c.sender.Propose(c.proposal)
			if c.err != nil {
				assert.EqualError(t, errors.Cause(err), c.err.Error())
			} else if c.code != codes.OK {
				MultiValidateStatusCode(t, err, c.code)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	for _, s := range servers {
		s.GracefulStop()
	}
}

func TestGrpcConsensusSender_Vote(t *testing.T) {
	confs := []*config.BBFTConfig{
		GetTestConfig(),
		GetTestConfig(),
		GetTestConfig(),
		GetTestConfig(),
	}
	confs[0].Port = "50053"
	confs[1].Port = "50054"
	confs[2].Port = "50055"
	confs[3].Port = "50056"

	ps := dba.NewPeerServiceOnMemory()
	for _, conf := range confs {
		ps.AddPeer(RandomPeerFromConf(conf))
	}

	servers := make([]*grpc.Server, 0, 4)
	for i, conf := range confs {
		servers = append(servers, NewTestGrpcServer())
		go func(conf *config.BBFTConfig, server *grpc.Server) {
			SetUpTestServer(t, conf, ps, server)
		}(conf, servers[i])
	}

	validVote := RandomVoteMessageFromPeer(t, ps.GetPeers()[0])
	unPeerValidVote := RandomVoteMessage(t)

	evilConf := *confs[0]
	pk, sk := convertor.NewKeyPair()
	evilConf.PublicKey = pk
	evilConf.SecretKey = sk

	sender := NewGrpcConsensusSender(confs[0], ps)
	evilSender := NewGrpcConsensusSender(&evilConf, ps)

	for _, c := range []struct {
		name   string
		vote   model.VoteMessage
		sender model.ConsensusSender
		code   codes.Code
		err    error
	}{
		{
			"success case",
			validVote,
			sender,
			codes.OK,
			nil,
		},
		{
			"failed case, authenticated but not peer",
			validVote,
			evilSender,
			codes.PermissionDenied,
			nil,
		},
		{
			"failed case, invalid vote(un peer signed)",
			unPeerValidVote,
			sender,
			codes.InvalidArgument,
			nil,
		},
		{
			"failed case, unsigned vote",
			convertor.NewModelFactory().NewVoteMessage(RandomByte()),
			sender,
			codes.InvalidArgument,
			nil,
		},
		{
			"failed case, nil",
			nil,
			sender,
			codes.Unauthenticated,
			model.ErrInvalidVoteMessage,
		},
		{
			"failed case, duplicate sent",
			validVote,
			sender,
			codes.AlreadyExists,
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			err := c.sender.Vote(c.vote)
			if c.err != nil {
				assert.EqualError(t, errors.Cause(err), c.err.Error())
			} else if c.code != codes.OK {
				MultiValidateStatusCode(t, err, c.code)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	for _, s := range servers {
		s.GracefulStop()
	}
}

func TestGrpcConsensusSender_PreCommit(t *testing.T) {
	confs := []*config.BBFTConfig{
		GetTestConfig(),
		GetTestConfig(),
		GetTestConfig(),
		GetTestConfig(),
	}
	confs[0].Port = "50053"
	confs[1].Port = "50054"
	confs[2].Port = "50055"
	confs[3].Port = "50056"

	ps := dba.NewPeerServiceOnMemory()
	for _, conf := range confs {
		ps.AddPeer(RandomPeerFromConf(conf))
	}

	servers := make([]*grpc.Server, 0, 4)
	for i, conf := range confs {
		servers = append(servers, NewTestGrpcServer())
		go func(conf *config.BBFTConfig, server *grpc.Server) {
			SetUpTestServer(t, conf, ps, server)
		}(conf, servers[i])
	}

	validVote := RandomVoteMessageFromPeer(t, ps.GetPeers()[0])
	unPeerValidVote := RandomVoteMessage(t)

	evilConf := *confs[0]
	pk, sk := convertor.NewKeyPair()
	evilConf.PublicKey = pk
	evilConf.SecretKey = sk

	sender := NewGrpcConsensusSender(confs[0], ps)
	evilSender := NewGrpcConsensusSender(&evilConf, ps)

	for _, c := range []struct {
		name   string
		vote   model.VoteMessage
		sender model.ConsensusSender
		code   codes.Code
		err    error
	}{
		{
			"success case",
			validVote,
			sender,
			codes.OK,
			nil,
		},
		{
			"failed case, authenticated but not peer",
			validVote,
			evilSender,
			codes.PermissionDenied,
			nil,
		},
		{
			"failed case, invalid vote(un peer signed)",
			unPeerValidVote,
			sender,
			codes.InvalidArgument,
			nil,
		},
		{
			"failed case, unsigned vote",
			convertor.NewModelFactory().NewVoteMessage(RandomByte()),
			sender,
			codes.InvalidArgument,
			nil,
		},
		{
			"failed case, nil",
			nil,
			sender,
			codes.Unauthenticated,
			model.ErrInvalidVoteMessage,
		},
		{
			"failed case, duplicate sent",
			validVote,
			sender,
			codes.AlreadyExists,
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			err := c.sender.PreCommit(c.vote)
			if c.err != nil {
				assert.EqualError(t, errors.Cause(err), c.err.Error())
			} else if c.code != codes.OK {
				MultiValidateStatusCode(t, err, c.code)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	for _, s := range servers {
		s.GracefulStop()
	}
}
