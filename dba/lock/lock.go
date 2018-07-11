package lock

import (
	"github.com/satellitex/bbft/model"
	"sync"
)

type Lock interface {
	AddVoteMessage(vote model.VoteMessage) error
	GetLockedProposal() (model.Proposal, bool)
}

type LockOnMemory struct {
	mutex *sync.Mutex
}

func NewLockOnMemory() Lock {
	return &LockOnMemory{new(sync.Mutex)}
}

func (lock *LockOnMemory) AddVoteMessage(vote model.VoteMessage) error {
	return nil
}

func (lock *LockOnMemory) GetLockedProposal() (model.Proposal, bool) {
	return nil, false
}
