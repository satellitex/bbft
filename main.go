package main

import (
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/satellitex/bbft/config"
	"github.com/satellitex/bbft/controller"
	"github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/dba"
	. "github.com/satellitex/bbft/grpc"
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
	"github.com/satellitex/bbft/usecase"
	"google.golang.org/grpc"
	"net"
)

func DemoGenesisCommit(conf *config.BBFTConfig, factory model.ModelFactory, bc dba.BlockChain, ps dba.PeerService) {
	genesisBlock, err := factory.NewBlock(0, nil, 0, nil)
	if err != nil {
		panic("DemoGenesisCommit: " + err.Error())
	}
	ps.AddPeer(factory.NewPeer(conf.Host+":"+conf.Port, conf.PublicKey))
	bc.Commit(genesisBlock)
}

func main() {

	fmt.Println("=========================== boot bbft ===========================")

	config.Init()
	conf := config.GetConfig()
	conf.PublicKey, conf.SecretKey = convertor.NewKeyPair()

	l, err := net.Listen("tcp", ":"+conf.Port)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Succcess New Listen")

	ps := dba.NewPeerServiceOnMemory()
	author := convertor.NewAuthor(ps)

	queue := dba.NewProposalTxQueueOnMemory(conf)
	lock := dba.NewLockOnMemory(ps, conf)
	pool := dba.NewReceiverPoolOnMemory(conf)
	bc := dba.NewBlockChainOnMemory()
	slv := convertor.NewStatelessValidator()
	sender := NewGrpcConsensusSender(conf, ps)
	receivChan := usecase.NewReceiveChannel(conf)

	consensusReceiver := usecase.NewConsensusReceiverUsecase(queue, ps, lock, pool, bc, slv, sender, receivChan)
	clientRceiver := usecase.NewClientGateReceiverUsecase(slv, sender)
	fmt.Println("Success New Receivers")

	s := grpc.NewServer([]grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_validator.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	}...)
	fmt.Println("Success New Server")

	bbft.RegisterConsensusGateServer(s, controller.NewConsensusController(consensusReceiver, author))
	bbft.RegisterTxGateServer(s, controller.NewClientGateController(clientRceiver, author))
	fmt.Println("Success New Register Endpoint")

	fmt.Println("Set Up!!")

	sfv := convertor.NewStatefulValidator(bc)
	factory := convertor.NewModelFactory()

	consensus := usecase.NewConsensusStepUsecase(conf, bc, ps, lock, queue, sender, slv, sfv, factory, receivChan)
	DemoGenesisCommit(conf, factory, bc, ps)

	// Consensus Run!!
	go func() {
		consensus.Run()
	}()

	if err := s.Serve(l); err != nil {
		fmt.Println("Failed to server grpc: ", err.Error())
	}

}
