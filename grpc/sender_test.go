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
	bc.Commit(RandomCommitableBlock(t, bc))

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

type TxGateSender struct {
	client bbft.TxGateClient
	t      *testing.T
}

func (s *TxGateSender) Write(ctx context.Context, tx model.Transaction) error {
	ptx := tx.(*convertor.Transaction).Transaction
	res, err := s.client.Write(ctx, ptx)
	if err == nil {
		require.Equal(s.t, &bbft.TxResponse{}, res)
	}
	return err
}

func NewTxGateSender(t *testing.T, conf *config.BBFTConfig) *TxGateSender {
	conn, err := grpc.Dial(conf.Host+":"+conf.Port, grpc.WithInsecure())
	require.NoError(t, err)
	return &TxGateSender{
		bbft.NewTxGateClient(conn),
		t,
	}
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

	fmt.Printf("%x", ps.GetPeers()[0].GetPubkey())
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
				assert.EqualError(t, errors.Cause(err), model.ErrInvalidTransaction.Error())
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
	conf := GetTestConfig()
	ps := RandomPeerService(t, 3)
	ps.AddPeer(RandomPeerFromConf(conf))

	server := NewTestGrpcServer()

	go func() {
		//WIP
		SetUpTestServer(t, conf, ps, server)
	}()

	leaderId := func() int32 {
		for i, p := range ps.GetPermutationPeers(1) {
			if bytes.Equal(conf.PublicKey, p.GetPubkey()) {
				return int32(i)
			}
		}
		return -1
	}()
	require.NotEqual(t, -1, leaderId)

	validProposal := RandomProposalWithHeightRound(t, 1, leaderId)
	invalidProposal := RandomInvalidProposal(t).(*convertor.Proposal)

	evilConf := *conf
	pk, sk := convertor.NewKeyPair()
	evilConf.PublicKey = pk
	evilConf.SecretKey = sk
	sender := NewGrpcConsensusSender(conf, ps)
	evilSender := NewGrpcConsensusSender(&evilConf, ps)

	notLeaderConf := evilConf
	notLeaderConf.PublicKey = ps.GetPermutationPeers(1)[(leaderId+1)%2].GetPubkey()
	notLeaderConf.PublicKey = ps.GetPermutationPeers(1)[(leaderId+1)%2].(*PeerWithPriv).PrivKey

	notLeaderSender := NewGrpcConsensusSender(&notLeaderConf, ps)

	fmt.Printf("%x", ps.GetPeers()[0].GetPubkey())
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
			"failed case, authenticated but not leader",
			validProposal,
			notLeaderSender,
			codes.PermissionDenied,
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
			model.ErrInvalidTransaction,
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
				assert.EqualError(t, errors.Cause(err), model.ErrInvalidTransaction.Error())
			} else if c.code != codes.OK {
				ValidateStatusCode(t, err, c.code)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	server.GracefulStop()
}
