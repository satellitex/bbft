package consensus

import (
	"github.com/satellitex/bbft/dba/lock"
	"github.com/satellitex/bbft/model"
)

type ConsensusReceiver interface {
	Proposal(proposal model.Proposal) error
	Vote(vote model.VoteMessage) error
	PreCommit(preCommit model.VoteMessage) error
}

type ConsensusReceieverUsecase struct {
	lock   lock.Lock
	sender model.ConsensusSender
}

func (c *ConsensusReceieverUsecase) Proposal(proposal model.Proposal) error {
	return nil
}

func (c *ConsensusReceieverUsecase) Vote(vote model.VoteMessage) error {
	return nil
}

func (c *ConsensusReceieverUsecase) PreCommit(preCommit model.VoteMessage) error {
	return nil
}