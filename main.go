package main

import (
	"encoding/base64"
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
	"os"
	"time"
)

func DecodeString64(s string) []byte {
	ret, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err.Error())
	}
	return ret
}

func DemoGenesisCommit(conf *config.BBFTConfig, factory model.ModelFactory, bc dba.BlockChain, ps dba.PeerService) {

	// myself addPeer
	conf.PublicKey = DecodeString64(conf.Demo.PublicKey)
	conf.SecretKey = DecodeString64(conf.Demo.SecretKey)

	ps.AddPeer(factory.NewPeer(conf.Host+":"+conf.Demo.Port1, DecodeString64(conf.Demo.Pubkey1)))
	ps.AddPeer(factory.NewPeer(conf.Host+":"+conf.Demo.Port2, DecodeString64(conf.Demo.Pubkey2)))
	ps.AddPeer(factory.NewPeer(conf.Host+":"+conf.Demo.Port3, DecodeString64(conf.Demo.Pubkey3)))
	ps.AddPeer(factory.NewPeer(conf.Host+":"+conf.Demo.Port4, DecodeString64(conf.Demo.Pubkey4)))

	genesisBlock, err := factory.NewBlock(0, nil, 0, nil)
	if err != nil {
		panic("DemoGenesisCommit: " + err.Error())
	}

	bc.Commit(genesisBlock)
}

func OnceNodeGenesis(conf *config.BBFTConfig, factory model.ModelFactory, bc dba.BlockChain, ps dba.PeerService) {

	// myself addPeer
	conf.PublicKey, conf.SecretKey = convertor.NewKeyPair()
	ps.AddPeer(factory.NewPeer(conf.Host+":"+conf.Port, conf.PublicKey))

	genesisBlock, err := factory.NewBlock(0, nil, 0, nil)
	if err != nil {
		panic("DemoGenesisCommit: " + err.Error())
	}

	bc.Commit(genesisBlock)
}

func main() {

	fmt.Println("=========================== boot bbft ===========================")

	config.Init()
	conf := config.GetConfig()

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

	if os.Getenv("DEMO") != "" {
		time.Sleep(time.Second*2)
		DemoGenesisCommit(conf, factory, bc, ps)
	} else {
		OnceNodeGenesis(conf, factory, bc, ps)
	}

	// Consensus Run!!
	go func() {
		consensus.Run()
	}()

	if err := s.Serve(l); err != nil {
		fmt.Println("Failed to server grpc: ", err.Error())
	}

}
