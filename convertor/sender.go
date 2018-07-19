package convertor

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/config"
	"github.com/satellitex/bbft/dba"
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
	"google.golang.org/grpc"
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
				fmt.Printf("Error Connection to peer: %#v", p)
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

func NewConsensusSender(conf *config.BBFTConfig, ps dba.PeerService) model.ConsensusSender {
	sender := &GrpcConsensusSender{conf: conf, manager: NewGrpcConnectManager(), ps: ps}
	return sender
}

func (s *GrpcConsensusSender) Propagate(tx model.Transaction) error {
	if proto, ok := tx.(*Transaction); ok {
		ctx, err := NewContextByProtobuf(s.conf, proto)
		if err != nil {
			return err
		}

		// BroadCast to All Peer in PeerService
		clientChan := make(chan bbft.ConsensusGateClient)
		s.manager.GetConnectsToChannel(s.ps.GetPeers(), clientChan)
		for client := range clientChan {
			go func() {
				if _, err := client.Propagate(ctx, proto.Transaction); err != nil {
					fmt.Printf("Failed Propagate Error : %s", err.Error())
				}
			}()
		}

	} else {
		return errors.Wrapf(model.ErrInvalidTransaction, "tx can not cast convertor.Transaction: %#v", tx)
	}

	return nil
}

func (s *GrpcConsensusSender) Propose(proposal model.Proposal) error {
	if _, ok := proposal.(*Proposal); !ok {
		return errors.Wrapf(model.ErrInvalidProposal, "proposal can not cast convertor.Proposal: %#v", proposal)
	}
	return nil
}

func (s *GrpcConsensusSender) Vote(vote model.VoteMessage) error {
	if _, ok := vote.(*VoteMessage); !ok {
		return errors.Wrapf(model.ErrInvalidVoteMessage, "vote can not cast to convertor.VoteMessage %#v", vote)
	}
	return nil
}

func (s *GrpcConsensusSender) PreCommit(vote model.VoteMessage) error {
	if _, ok := vote.(*VoteMessage); !ok {
		return errors.Wrapf(model.ErrInvalidVoteMessage, "vote can not cast to convertor.VoteMessage %#v", vote)
	}
	return nil
}

type MockConsensusSender struct {
	Tx               model.Transaction
	Proposal         model.Proposal
	VoteMessage      model.VoteMessage
	PreCommitMessage model.VoteMessage
}

func NewMockConsensusSender() model.ConsensusSender {
	return &MockConsensusSender{}
}

func (s *MockConsensusSender) Propagate(tx model.Transaction) error {
	if _, ok := tx.(*Transaction); !ok {
		return errors.Wrapf(model.ErrInvalidTransaction, "tx can not cast convertor.Transaction: %#v", tx)
	}
	s.Tx = tx
	return nil
}

func (s *MockConsensusSender) Propose(proposal model.Proposal) error {
	if _, ok := proposal.(*Proposal); !ok {
		return errors.Wrapf(model.ErrInvalidProposal, "proposal can not cast convertor.Proposal: %#v", proposal)
	}
	s.Proposal = proposal
	return nil
}

func (s *MockConsensusSender) Vote(vote model.VoteMessage) error {
	if _, ok := vote.(*VoteMessage); !ok {
		return errors.Wrapf(model.ErrInvalidVoteMessage, "vote can not cast to convertor.VoteMessage %#v", vote)
	}
	s.VoteMessage = vote
	return nil
}

func (s *MockConsensusSender) PreCommit(vote model.VoteMessage) error {
	if _, ok := vote.(*VoteMessage); !ok {
		return errors.Wrapf(model.ErrInvalidVoteMessage, "vote can not cast to convertor.VoteMessage %#v", vote)
	}
	s.PreCommitMessage = vote
	return nil
}
