package dba

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
	"sync"
)

type Lock interface {
	RegisterProposal(model.Proposal) error
	AddVoteMessage(vote model.VoteMessage) (bool, error)
	GetLockedProposal() (model.Proposal, bool)
}

type LockOnMemory struct {
	peerService        PeerService
	lockedProposal     model.Proposal
	registerdProposals map[string]model.Proposal
	acceptedPrposal    map[string]int
	mutex              *sync.Mutex
}

var (
	ErrSetLockedProposal = errors.Errorf("Failed SetLocked Proposal")
)

func NewLockOnMemory(peerService PeerService) Lock {
	return &LockOnMemory{
		peerService,
		nil, make(map[string]model.Proposal),
		make(map[string]int), new(sync.Mutex),
	}
}

func (lock *LockOnMemory) getAllowFailed() int {
	return (lock.peerService.Size() - 1) / 3
}

func (lock *LockOnMemory) getRequiredAccepet() int {
	return lock.getAllowFailed()*2 + 1
}

func (lock *LockOnMemory) setLockedProposal(proposal model.Proposal) (bool, error) {
	if proposal == nil {
		return false, errors.Errorf("set proposal is nil")
	}
	if lock.lockedProposal == nil {
		lock.lockedProposal = proposal
		return true, nil
	} else {
		lockedHash, err := lock.lockedProposal.GetBlock().GetHash()
		if err != nil {
			return false, errors.Wrapf(model.ErrBlockGetHash, "locked Hash: "+err.Error())
		}
		newHash, err := proposal.GetBlock().GetHash()
		if err != nil {
			return false, errors.Wrapf(model.ErrBlockGetHash, "new Hash: "+err.Error())
		}
		if !bytes.Equal(lockedHash, newHash) &&
			proposal.GetRound() > lock.lockedProposal.GetRound() {
			lock.lockedProposal = proposal
			return true, nil
		}
	}
	return false, nil
}

func (lock *LockOnMemory) RegisterProposal(proposal model.Proposal) error {
	defer lock.mutex.Unlock()
	lock.mutex.Lock()

	if proposal == nil {
		return errors.Wrapf(model.ErrInvalidProposal, "Proposal is nil")
	}
	hash, err := proposal.GetBlock().GetHash()
	if err != nil {
		return errors.Wrapf(model.ErrBlockGetHash, err.Error())
	}
	lock.registerdProposals[string(hash)] = proposal
	return nil
}

func (lock *LockOnMemory) AddVoteMessage(vote model.VoteMessage) (bool, error) {
	defer lock.mutex.Unlock()
	lock.mutex.Lock()

	if vote == nil {
		return false, errors.Wrapf(model.ErrInvalidVoteMessage, "VoteMessage is nil")
	}
	hash := string(vote.GetBlockHash())
	lock.acceptedPrposal[hash] += 1
	if lock.acceptedPrposal[hash] >= lock.getRequiredAccepet() {
		if ok, err := lock.setLockedProposal(lock.registerdProposals[hash]); !ok {
			if err != nil {
				return false, errors.Wrapf(ErrSetLockedProposal, err.Error())
			}
		} else {
			return true, nil
		}
	}
	return false, nil
}

func (lock *LockOnMemory) GetLockedProposal() (model.Proposal, bool) {
	if lock.lockedProposal == nil {
		return nil, false
	}
	return lock.lockedProposal, true
}
