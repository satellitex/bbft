package convertor

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
)

type GrpcConsensusSender struct {
	client bbft.ConsensusGateClient
}

func NewConsensusSender() model.ConsensusSender {
	return &GrpcConsensusSender{nil}
}

func (s *GrpcConsensusSender) Propagate(tx model.Transaction) error {
	if _, ok := tx.(*Transaction); !ok {
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
}

func NewMockConsensusSender() model.ConsensusSender {
	return &MockConsensusSender{}
}

func (s *MockConsensusSender) Propagate(tx model.Transaction) error {
	if _, ok := tx.(*Transaction); !ok {
		return errors.Wrapf(model.ErrInvalidTransaction, "tx can not cast convertor.Transaction: %#v", tx)
	}

	return nil
}

func (s *MockConsensusSender) Propose(proposal model.Proposal) error {
	if _, ok := proposal.(*Proposal); !ok {
		return errors.Wrapf(model.ErrInvalidProposal, "proposal can not cast convertor.Proposal: %#v", proposal)
	}

	return nil
}

func (s *MockConsensusSender) Vote(vote model.VoteMessage) error {
	if _, ok := vote.(*VoteMessage); !ok {
		return errors.Wrapf(model.ErrInvalidVoteMessage, "vote can not cast to convertor.VoteMessage %#v", vote)
	}
	return nil
}

func (s *MockConsensusSender) PreCommit(vote model.VoteMessage) error {
	if _, ok := vote.(*VoteMessage); !ok {
		return errors.Wrapf(model.ErrInvalidVoteMessage, "vote can not cast to convertor.VoteMessage %#v", vote)
	}

	return nil
}
