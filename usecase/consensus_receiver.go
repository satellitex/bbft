package usecase

import (
	"github.com/satellitex/bbft/dba"
	"github.com/satellitex/bbft/model"
)

type ConsensusReceiver interface {
	Propagate(ptx model.Transaction) error
	Propose(proposal model.Proposal) error
	Vote(vote model.VoteMessage) error
	PreCommit(preCommit model.VoteMessage) error
}

type ConsensusReceieverUsecase struct {
	queue  dba.ProposalTxQueue
	lock   dba.Lock
	sender model.ConsensusSender
}

func (c *ConsensusReceieverUsecase) Propagate(proposalTx model.Transaction) error {
	return nil
}

func (c *ConsensusReceieverUsecase) Propose(proposal model.Proposal) error {
	return nil
}

func (c *ConsensusReceieverUsecase) Vote(vote model.VoteMessage) error {
	return nil
}

func (c *ConsensusReceieverUsecase) PreCommit(preCommit model.VoteMessage) error {
	return nil
}
