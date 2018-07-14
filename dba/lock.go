package dba

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
	"sync"
)

type Lock interface {
	RegisterProposal(model.Proposal) error
	AddVoteMessage(vote model.VoteMessage) bool
	GetLockedProposal() (model.Proposal, bool)
}

type LockOnMemory struct {
	peerService        PeerService
	lockedProposal     model.Proposal
	registerdProposals map[string]model.Proposal
	acceptedPrposal    map[string]int
	mutex              *sync.Mutex
}

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

func (lock *LockOnMemory) setLockedProposal(proposal model.Proposal) {
	if lock.lockedProposal == nil {
		lock.lockedProposal = proposal
	} else {
		if !bytes.Equal(model.MustGetHash(lock.lockedProposal.GetBlock()),
			model.MustGetHash(proposal.GetBlock())) &&
			proposal.GetRound() > lock.lockedProposal.GetRound() {
			lock.lockedProposal = proposal
		}
	}
}

func (lock *LockOnMemory) RegisterProposal(proposal model.Proposal) error {
	defer lock.mutex.Unlock()
	lock.mutex.Lock()
	hash, err := proposal.GetBlock().GetHash()
	if err != nil {
		return errors.Wrapf(model.ErrBlockGetHash, err.Error())
	}
	lock.registerdProposals[string(hash)] = proposal
	return nil
}

func (lock *LockOnMemory) AddVoteMessage(vote model.VoteMessage) bool {
	defer lock.mutex.Unlock()
	lock.mutex.Lock()

	hash := string(vote.GetBlockHash())
	lock.acceptedPrposal[hash] += 1
	if lock.acceptedPrposal[hash] >= lock.getRequiredAccepet() {
		lock.setLockedProposal(lock.registerdProposals[hash])
		return true
	}

	return false
}

func (lock *LockOnMemory) GetLockedProposal() (model.Proposal, bool) {
	if lock.lockedProposal == nil {
		return nil, false
	}
	return lock.lockedProposal, true
}
