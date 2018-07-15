package dba

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/config"
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
	registeredQueue    []string

	acceptedCounter map[string]int
	findedVote      map[string]model.VoteMessage
	votedQueue      []string

	lockedLimit int
	votedLimit  int
	mutex       *sync.Mutex
}

var (
	ErrValidInPeerService    = errors.Errorf("Failed Valid In PeerService")
	ErrValidLockedProposal   = errors.Errorf("Failed Valid Locked Proposal")
	ErrAlreadyAddVoteMessage = errors.Errorf("Failed Alrady add same VoteMessage")
)

func NewLockOnMemory(peerService PeerService, cnf *config.BBFTConfig) Lock {
	return &LockOnMemory{
		peerService,
		nil,
		make(map[string]model.Proposal), make([]string, cnf.LockedRegisteredLimits),
		make(map[string]int),
		make(map[string]model.VoteMessage), make([]string, cnf.LockedVotedLimits),
		cnf.LockedRegisteredLimits,
		cnf.LockedVotedLimits,
		new(sync.Mutex),
	}
}

func getAllowFailed(ps PeerService) int {
	return (ps.Size() - 1) / 3
}

func getRequiredAccepet(ps PeerService) int {
	return getAllowFailed(ps)*2 + 1
}

func validInPeerService(vote model.VoteMessage, ps PeerService) bool {
	if _, ok := ps.GetPeer(vote.GetSignature().GetPubkey()); ok {
		return true
	}
	return false
}

func validLockedProposal(proposal model.Proposal, lockedProposal model.Proposal) (bool, error) {
	if proposal == nil {
		return false, errors.Wrapf(model.ErrInvalidProposal, "set proposal is nil")
	}
	if lockedProposal == nil {
		return true, nil
	} else {
		lockedHash, err := lockedProposal.GetBlock().GetHash()
		if err != nil {
			return false, errors.Wrapf(model.ErrBlockGetHash, "locked Hash: "+err.Error())
		}
		newHash, err := proposal.GetBlock().GetHash()
		if err != nil {
			return false, errors.Wrapf(model.ErrBlockGetHash, "new Hash: "+err.Error())
		}
		if !bytes.Equal(lockedHash, newHash) &&
			proposal.GetRound() > lockedProposal.GetRound() {
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
	// register proposal ===
	lock.registerdProposals[string(hash)] = proposal
	lock.registeredQueue = append(lock.registeredQueue, string(hash))
	if len(lock.registeredQueue) >= lock.lockedLimit {
		delete(lock.registerdProposals, lock.registeredQueue[0])
		lock.registeredQueue = lock.registeredQueue[1:]
	}
	// ===
	return nil
}

func (lock *LockOnMemory) AddVoteMessage(vote model.VoteMessage) (bool, error) {
	defer lock.mutex.Unlock()
	lock.mutex.Lock()

	if vote == nil {
		return false, errors.Wrapf(model.ErrInvalidVoteMessage, "VoteMessage is nil")
	}

	//if ok := validInPeerService(vote, lock.peerService); ok {
	//	return false, errors.Wrapf(ErrValidInPeerService, "PeerService doesn't have pubkey: %x", vote.GetSignature().GetPubkey())
	//}

	hash := string(vote.GetBlockHash())
	pub := string(vote.GetSignature().GetPubkey())
	if _, ok := lock.findedVote[hash+pub]; ok {
		return false, errors.Wrapf(ErrAlreadyAddVoteMessage, "already add vote: %#v", vote)
	}

	// add vote ===
	lock.findedVote[hash+pub] = vote
	lock.votedQueue = append(lock.votedQueue, hash+pub)
	if len(lock.votedQueue) == lock.votedLimit { // shits Limits
		delete(lock.acceptedCounter, string(lock.findedVote[lock.votedQueue[0]].GetBlockHash()))
		delete(lock.findedVote, lock.votedQueue[0])
		lock.votedQueue = lock.votedQueue[1:]
	}
	lock.acceptedCounter[hash]++
	// ===

	if getRequiredAccepet(lock.peerService) <= lock.acceptedCounter[hash] {
		proposal := lock.registerdProposals[hash]
		if ok, err := validLockedProposal(proposal, lock.lockedProposal); ok {
			lock.lockedProposal = proposal
			return true, nil
		} else {
			if err != nil {
				return false, errors.Wrapf(ErrValidLockedProposal, err.Error())
			}
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
