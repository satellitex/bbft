package usecase

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/dba"
	"github.com/satellitex/bbft/model"
)

var (
	ErrAlradyReceivedSameObject = errors.New("Failed Alrady Received Same Object")
)

type ConsensusReceiver interface {
	Propagate(tx model.Transaction) error
	Propose(proposal model.Proposal) error
	Vote(vote model.VoteMessage) error
	PreCommit(preCommit model.VoteMessage) error
}

type ConsensusReceieverUsecase struct {
	queue  dba.ProposalTxQueue
	lock   dba.Lock
	pool   dba.ReceiverPool
	sender model.ConsensusSender
}

func NewConsensusReceiverUsecase(queue dba.ProposalTxQueue, lock dba.Lock, pool dba.ReceiverPool, sender model.ConsensusSender) ConsensusReceiver {
	return &ConsensusReceieverUsecase{
		queue:  queue,
		lock:   lock,
		pool:   pool,
		sender: sender,
	}
}

func (c *ConsensusReceieverUsecase) Propagate(tx model.Transaction) error {
	if tx == nil { // InvalidArgument (code = 3)
		return errors.Wrapf(model.ErrInvalidTransaction, "tx is nil")
	}
	if err := tx.Verify(); err != nil { // InvalidArgument (code = 3)
		return errors.Wrapf(model.ErrTransactionGetHash, err.Error())
	}
	if c.pool.IsExistPropagate(tx) { // AlreadyExist (code = 6)
		return errors.Wrapf(ErrAlradyReceivedSameObject, "tx: %#v", tx)
	}
	if err := c.queue.Push(tx); err != nil { // ResourceExhausted (code = 8)
		return errors.Wrapf(dba.ErrProposalTxQueuePush, err.Error())
	}
	if err := c.pool.SetPropagate(tx); err != nil {
		return errors.Wrapf(dba.ErrReceiverPoolSet, err.Error())
	}
	if err := c.sender.Propagate(tx); err != nil {
		return errors.Wrapf(model.ErrConsensusSenderPropagate, err.Error())
	}
	return nil
}

func (c *ConsensusReceieverUsecase) Propose(proposal model.Proposal) error {
	if proposal == nil {
		return errors.Wrapf(model.ErrInvalidProposal, "proposal is nil")
	}

	if err := c.sender.Propose(proposal); err != nil {
		return errors.Wrapf(model.ErrConsensusSenderPropagate, err.Error())
	}
	return nil
}

func (c *ConsensusReceieverUsecase) Vote(vote model.VoteMessage) error {
	if vote == nil {
		return errors.Wrapf(model.ErrInvalidVoteMessage, "vote is nil")
	}

	if err := c.sender.Vote(vote); err != nil {
		return errors.Wrapf(model.ErrConsensusSenderPropagate, err.Error())
	}
	return nil
}

func (c *ConsensusReceieverUsecase) PreCommit(preCommit model.VoteMessage) error {
	if preCommit == nil {
		return errors.Wrapf(model.ErrInvalidVoteMessage, "preCommit is nil")
	}

	if err := c.sender.PreCommit(preCommit); err != nil {
		return errors.Wrapf(model.ErrConsensusSenderPropagate, err.Error())
	}
	return nil
}
