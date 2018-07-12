package usecase

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/dba"
	"github.com/satellitex/bbft/model"
)

type ConsensusStep interface {
	Run()
	Propose() error
	Vote() error
	PreCommit() error
}

type ConsensusStepUsecase struct {
	bc                 dba.BlockChain
	lock               dba.Lock
	queue              dba.ProposalTxQueue
	sender             model.ConsensusSender
	statelessValidator model.StatelessValidator
	statefulValidator  model.StatefulValidator
	factory            model.ModelFactory
}

var (
	ErrConsensusProposal  = errors.Errorf("Failed This peer Proposal")
	ErrConsensusVote      = errors.Errorf("Failed This peer Vote")
	ErrConsensusPreCommit = errors.Errorf("Failed This peer PreCommit")
	ErrConsensusCommit    = errors.Errorf("Failed This peer ConsensusCommit")
)

func (c *ConsensusStepUsecase) Run() {
	for {
		for {
			c.Propose()
			c.Vote()
			c.PreCommit()
		}
		c.Commit()
	}
}

func (c *ConsensusStepUsecase) Propose() error {
	return nil
}

func (c *ConsensusStepUsecase) Vote() error {
	return nil
}

func (c *ConsensusStepUsecase) PreCommit() error {
	return nil
}

func (c *ConsensusStepUsecase) Commit() error {
	proposal, ok := c.lock.GetLockedProposal()
	if !ok {
		return errors.Wrapf(ErrConsensusCommit,
			"Not Founbd Locked Proposal")
	}
	block := proposal.GetBlock()
	c.bc.Commit(block)
	return nil
}
