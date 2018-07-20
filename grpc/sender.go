package grpc

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/config"
	. "github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/dba"
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"sync"
)

type GrpcConnectionManager struct {
	clients map[string]bbft.ConsensusGateClient
}

func NewGrpcConnectManager() *GrpcConnectionManager {
	return &GrpcConnectionManager{
		make(map[string]bbft.ConsensusGateClient),
	}
}

func (m *GrpcConnectionManager) CreateConn(peer model.Peer) error {
	// TODO now Insecure...?
	gc, err := grpc.Dial(peer.GetAddress(), grpc.WithInsecure())
	if err != nil {
		return err
	}
	m.clients[peer.GetAddress()] = bbft.NewConsensusGateClient(gc)
	return nil
}

func (m *GrpcConnectionManager) GetConnectsToChannel(peers []model.Peer, ret chan bbft.ConsensusGateClient) {
	for _, p := range peers {
		if client, ok := m.clients[p.GetAddress()]; ok {
			ret <- client
		} else {
			if err := m.CreateConn(p); err != nil {
				log.Println("Error Connection to peer: ", p)
				return
			}
			ret <- m.clients[p.GetAddress()]
		}
	}
	close(ret)
}

type GrpcConsensusSender struct {
	conf    *config.BBFTConfig
	manager *GrpcConnectionManager
	ps      dba.PeerService
}

func NewGrpcConsensusSender(conf *config.BBFTConfig, ps dba.PeerService) model.ConsensusSender {
	sender := &GrpcConsensusSender{conf: conf, manager: NewGrpcConnectManager(), ps: ps}
	return sender
}

func (s *GrpcConsensusSender) broadCast(send func(bbft.ConsensusGateClient, chan error, *sync.WaitGroup)) error {
	// BroadCast to All Peer in PeerService
	clientChan := make(chan bbft.ConsensusGateClient)
	peers := s.ps.GetPeers()
	go func() {
		s.manager.GetConnectsToChannel(peers, clientChan)
	}()

	errChan := make(chan error)
	resultChan := make(chan error)
	waiter := &sync.WaitGroup{}
	go func() {
		var errs error
		for err := range errChan {
			errs = multierr.Append(errs, err)
		}
		resultChan <- errs
	}()
	for client := range clientChan {
		waiter.Add(1)
		go send(client, errChan, waiter)
	}
	waiter.Wait()
	close(errChan)
	return <-resultChan
}

func (s *GrpcConsensusSender) Propagate(tx model.Transaction) error {
	if proto, ok := tx.(*Transaction); ok {
		ctx, err := NewContextByProtobuf(s.conf, proto)
		if err != nil {
			return err
		}

		// BroadCast to All Peer in PeerService
		return s.broadCast(
			func(c bbft.ConsensusGateClient, errChan chan error, waiter *sync.WaitGroup) {
				if _, err := c.Propagate(ctx, proto.Transaction); err != nil {
					errChan <- err
				}
				waiter.Done()
			})
	} else {
		return errors.Wrapf(model.ErrInvalidTransaction, "tx can not cast convertor.Transaction: %#v", tx)
	}

	return nil
}

func (s *GrpcConsensusSender) Propose(proposal model.Proposal) error {
	evilBlock, _ := NewModelFactory().NewBlock(proposal.GetBlock().GetHeader().GetHeight(),
		proposal.GetBlock().GetHeader().GetPreBlockHash(),
		proposal.GetBlock().GetHeader().GetCreatedTime()+1,
		proposal.GetBlock().GetTransactions())
	evilProposal, _ := NewModelFactory().NewProposal(evilBlock, proposal.GetRound())
	evilProto, _ := evilProposal.(*Proposal)
	if proto, ok := proposal.(*Proposal); ok {
		ctx, err := NewContextByProtobuf(s.conf, proto)
		if err != nil {
			return err
		}

		// BroadCast to All Peer in PeerService
		return s.broadCast(
			func(c bbft.ConsensusGateClient, errChan chan error, waiter *sync.WaitGroup) {
				if rand.Int()%2 == 0 {
					log.Println("Evil Propose Pettern NormalProto!!!!!")
					if _, err := c.Propose(ctx, proto.Proposal); err != nil {
						errChan <- err
					}
				} else {
					log.Println("Evil Propose Pettern EvilProto!!!!!")
					if _, err := c.Propose(ctx, evilProto.Proposal); err != nil {
						errChan <- err
					}
				}
				waiter.Done()
			})
	} else {
		return errors.Wrapf(model.ErrInvalidProposal, "proposal can not cast convertor.Proposal: %#v", proposal)
	}
	return nil
}

func (s *GrpcConsensusSender) Vote(vote model.VoteMessage) error {
	if proto, ok := vote.(*VoteMessage); ok {
		ctx, err := NewContextByProtobuf(s.conf, proto)
		if err != nil {
			return err
		}

		// BroadCast to All Peer in PeerService
		return s.broadCast(
			func(c bbft.ConsensusGateClient, errChan chan error, waiter *sync.WaitGroup) {
				if _, err := c.Vote(ctx, proto.VoteMessage); err != nil {
					errChan <- err
				}
				waiter.Done()
			})
	} else {
		return errors.Wrapf(model.ErrInvalidVoteMessage, "vote can not cast to convertor.VoteMessage %#v", vote)
	}
	return nil
}

func (s *GrpcConsensusSender) PreCommit(vote model.VoteMessage) error {
	if proto, ok := vote.(*VoteMessage); ok {
		ctx, err := NewContextByProtobuf(s.conf, proto)
		if err != nil {
			return err
		}

		// BroadCast to All Peer in PeerService
		return s.broadCast(
			func(c bbft.ConsensusGateClient, errChan chan error, waiter *sync.WaitGroup) {
				if _, err := c.PreCommit(ctx, proto.VoteMessage); err != nil {
					errChan <- err
				}
				waiter.Done()
			})
	} else {
		return errors.Wrapf(model.ErrInvalidVoteMessage, "vote can not cast to convertor.VoteMessage %#v", vote)
	}
	return nil
}
