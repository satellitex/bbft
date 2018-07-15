package usecase

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/dba"
	"github.com/satellitex/bbft/model"
	"go.uber.org/multierr"
)

var (
	ErrAlradyReceivedSameObject = errors.New("Failed Alrady Received Same Object")
	ErrVoteNotInPeerService     = errors.New("Failed vote's pubkey doesn't exist in peerService")
)

type ConsensusReceiver interface {
	Propagate(tx model.Transaction) error
	Propose(proposal model.Proposal) error
	Vote(vote model.VoteMessage) error
	PreCommit(preCommit model.VoteMessage) error
}

type ConsensusReceieverUsecase struct {
	queue  dba.ProposalTxQueue
	ps     dba.PeerService
	lock   dba.Lock
	pool   dba.ReceiverPool
	bc     dba.BlockChain
	slv    model.StatelessValidator
	sender model.ConsensusSender
}

func NewConsensusReceiverUsecase(queue dba.ProposalTxQueue, ps dba.PeerService, lock dba.Lock, pool dba.ReceiverPool, bc dba.BlockChain, slv model.StatelessValidator, sender model.ConsensusSender) ConsensusReceiver {
	return &ConsensusReceieverUsecase{
		queue:  queue,
		ps:     ps,
		lock:   lock,
		pool:   pool,
		bc:     bc,
		slv:    slv,
		sender: sender,
	}
}

func (c *ConsensusReceieverUsecase) Propagate(tx model.Transaction) error {
	if err := c.slv.TxValidate(tx); err != nil { // InvalidArgument (code = 3)
		return errors.Wrapf(model.ErrStatelessTxValidate, err.Error())
	}
	if c.pool.IsExistPropagate(tx) { // AlreadyExist (code = 6)
		return errors.Wrapf(ErrAlradyReceivedSameObject, "tx: %#v", tx)
	}

	// After parallel
	errs := make(chan error)
	go func() { // === main calc ===
		if err := c.queue.Push(tx); err != nil {
			errs <- errors.Wrapf(dba.ErrProposalTxQueuePush, err.Error())
		}
		errs <- nil
	}()
	go func() {
		if err := c.pool.SetPropagate(tx); err != nil {
			errs <- errors.Wrapf(dba.ErrReceiverPoolSet, err.Error())
		}
		errs <- nil
	}()
	go func() {
		if err := c.sender.Propagate(tx); err != nil {
			errs <- errors.Wrapf(model.ErrConsensusSenderPropagate, err.Error())
		}
		errs <- nil
	}()
	var result error
	result = multierr.Append(result, <-errs)
	result = multierr.Append(result, <-errs)
	result = multierr.Append(result, <-errs)
	return result
}

func (c *ConsensusReceieverUsecase) Propose(proposal model.Proposal) error {
	if proposal == nil { // InvalidArgument (code = 3)
		return errors.Wrapf(model.ErrInvalidProposal, "proposal is nil")
	}
	if err := c.slv.BlockValidate(proposal.GetBlock()); err != nil { // InvalidArgument (code = 3)
		return errors.Wrapf(model.ErrStatelessBlockValidate, err.Error())
	}
	if c.pool.IsExistPropose(proposal) { // AlreadyExist (code = 6)
		return errors.Wrapf(ErrAlradyReceivedSameObject, "proposal: %#v", proposal)
	}

	// After parallel
	errs := make(chan error)
	go func() { // === main calc ===
		if err := c.lock.RegisterProposal(proposal); err != nil {
			errs <- errors.Wrapf(dba.ErrLockRegisteredProposal, err.Error())
		}
		errs <- nil
	}()
	go func() {
		if err := c.pool.SetPropose(proposal); err != nil {
			errs <- errors.Wrapf(dba.ErrReceiverPoolSet, err.Error())
		}
		errs <- nil
	}()
	go func() {
		if err := c.sender.Propose(proposal); err != nil {
			errs <- errors.Wrapf(model.ErrConsensusSenderPropagate, err.Error())
		}
		errs <- nil
	}()
	var result error
	result = multierr.Append(result, <-errs)
	result = multierr.Append(result, <-errs)
	result = multierr.Append(result, <-errs)
	return result
}

func (c *ConsensusReceieverUsecase) Vote(vote model.VoteMessage) error {
	if vote == nil { // InvalidArgument (code = 3)
		return errors.Wrapf(model.ErrInvalidVoteMessage, "vote is nil")
	}
	if err := vote.Verify(); err != nil { // InvalidArgument (code = 3)
		return errors.Wrapf(model.ErrVoteMessageVerify, err.Error())
	}
	if _, ok := c.ps.GetPeer(vote.GetSignature().GetPubkey()); !ok { // InvalidArgument (code = 3)
		return errors.Wrapf(ErrVoteNotInPeerService, "pubkey: %x", vote.GetSignature().GetPubkey())
	}
	if c.pool.IsExistVote(vote) { // AlreadyExist (code = 6)
		return errors.Wrapf(ErrAlradyReceivedSameObject, "vote: %#v", vote)
	}

	// after parallel
	errs := make(chan error)
	go func() {
		if err := c.lock.AddVoteMessage(vote); err != nil {
			errs <- errors.Wrapf(dba.ErrLockAddVoteMessage, err.Error())
		}
		errs <- nil
	}()
	go func() {
		if err := c.pool.SetVote(vote); err != nil {
			errs <- errors.Wrap(dba.ErrReceiverPoolSet, err.Error())
		}
		errs <- nil
	}()
	go func() {
		if err := c.sender.Vote(vote); err != nil {
			errs <- errors.Wrapf(model.ErrConsensusSenderPropagate, err.Error())
		}
		errs <- nil
	}()
	var result error
	result = multierr.Append(result, <-errs)
	result = multierr.Append(result, <-errs)
	result = multierr.Append(result, <-errs)
	return result
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
