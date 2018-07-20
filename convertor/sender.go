package convertor

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
)

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
