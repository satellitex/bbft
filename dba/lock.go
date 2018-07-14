package dba

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
	"sync"
)

type Lock interface {
	RegisterProposal(model.Proposal) error
	AddVoteMessage(vote model.VoteMessage) error
	GetLockedProposal() (model.Proposal, bool)
}

type LockOnMemory struct {
	lockedProposal     model.Proposal
	registerdProposals map[string]model.Proposal
	acceptedPrposal    map[string]int64
	mutex              *sync.Mutex
}

func NewLockOnMemory() Lock {
	return &LockOnMemory{nil, make(map[string]int64), new(sync.Mutex)}
}

func (lock *LockOnMemory) RegisterProposal(proposal model.Proposal) error {
	defer lock.mutex.Unlock()
	lock.mutex.Lock()
	hash, err := proposal.GetBlock().GetHash()
	if err != nil {
		return errors.Wrapf(model.ErrBlockGetHash, err)
	}
	lock.registerdProposals[proposal.GetBlock().GetHash()] = proposal
	return nil
}

func (lock *LockOnMemory) AddVoteMessage(vote model.VoteMessage) error {
	defer lock.mutex.Unlock()
	lock.mutex.Lock()
	return nil
}

func (lock *LockOnMemory) GetLockedProposal() (model.Proposal, bool) {
	if lock.lockedProposal == nil {
		return nil, false
	}
	return lock.lockedProposal, true
}
