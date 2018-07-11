package infrastructure

import (
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
)

type GrpcConsensusSender struct {
	client bbft.ConsensusGateClient
}

func NewConsensusSender() model.ConsensusSender {
	return &GrpcConsensusSender{nil}
}

func (s *GrpcConsensusSender) Propagate(ptx model.ProposalTx) error {
	return nil
}

func (s *GrpcConsensusSender) Propose(proposal model.Proposal) error {
	return nil
}

func (s *GrpcConsensusSender) Vote(vote model.VoteMessage) error {
	return nil
}

func (s *GrpcConsensusSender) PreCommit(vote model.VoteMessage) error {
	return nil
}
