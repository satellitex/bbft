package dba

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/config"
	"github.com/satellitex/bbft/model"
	"sync"
)

// Lock は 2/3以上のAcceptedVoteを獲得したProposalを管理する
//
// 各 Height について、 2/3以上の AcceptedVote を獲得した Proposal が複数あるとき Round の大きい方の Lock を取る。
type Lock interface {
	// Proposal を登録する。
	RegisterProposal(model.Proposal) error

	// Vote を登録する。
	// Error Case )
	// 	1 ) Vote が nil の場合
	//  2 ) 既にそのVoteが登録されていた場合
	AddVoteMessage(vote model.VoteMessage) error
	// 高さ height における Lock を取得する。存在しなければ bool = false, otherwise true
	GetLockedProposal(height int64) (model.Proposal, bool)
	// ある高さ未満の Lock をすべて消す。
	Clean(height int64)
}

type LockOnMemory struct {
	peerService        PeerService
	lockedProposal     map[int64]model.Proposal
	registerdProposals map[string]model.Proposal
	registeredQueue    []string

	acceptedCounter map[string]int
	findedVote      map[string]model.VoteMessage
	votedQueue      []string

	mutex *sync.Mutex
}

var (
	ErrLockAddVoteMessage      = errors.New("Failed Lock Add VoteMessage")
	ErrLockRegisteredProposal  = errors.New("Faild Lock Registered Proposal")
	ErrAlreadyAddVoteMessage   = errors.New("Failed Alrady add same VoteMessage")
	ErrAlreadyRegisterProposal = errors.New("Failed Alrady register same Proposal")
)

func NewLockOnMemory(peerService PeerService, cnf *config.BBFTConfig) Lock {
	return &LockOnMemory{
		peerService,
		make(map[int64]model.Proposal),
		make(map[string]model.Proposal), make([]string, 0, cnf.LockedRegisteredLimits),
		make(map[string]int),
		make(map[string]model.VoteMessage), make([]string, 0, cnf.LockedVotedLimits),
		new(sync.Mutex),
	}
}

func getAllowFailed(ps PeerService) int {
	return (ps.Size() - 1) / 3
}

func getRequiredAccepet(ps PeerService) int {
	return getAllowFailed(ps)*2 + 1
}

func validLockedProposal(proposal model.Proposal, lockedProposal model.Proposal) bool {
	if proposal == nil {
		return false
	}
	if lockedProposal == nil {
		return true
	} else {
		if proposal.GetRound() > lockedProposal.GetRound() {
			return true
		}
	}
	return false
}

func (lock *LockOnMemory) checkAndLock(hash string) {
	if proposal, ok := lock.registerdProposals[hash]; ok {
		height := proposal.GetBlock().GetHeader().GetHeight()
		if getRequiredAccepet(lock.peerService) <= lock.acceptedCounter[hash] {
			if ok := validLockedProposal(proposal, lock.lockedProposal[height]); ok {
				lock.lockedProposal[height] = proposal
			}
		}
	}
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
	if _, ok := lock.registerdProposals[string(hash)]; ok {
		return errors.Wrapf(ErrAlreadyRegisterProposal, "alrady register proposal: %#v", proposal)
	}

	// === register proposal ===
	if len(lock.registeredQueue) >= cap(lock.registeredQueue) { // shifts Limits
		delete(lock.registerdProposals, lock.registeredQueue[0])
		lock.registeredQueue = lock.registeredQueue[1:]
	}
	lock.registerdProposals[string(hash)] = proposal
	lock.registeredQueue = append(lock.registeredQueue, string(hash))
	// =========================
	lock.checkAndLock(string(hash))
	return nil
}

func (lock *LockOnMemory) AddVoteMessage(vote model.VoteMessage) error {
	defer lock.mutex.Unlock()
	lock.mutex.Lock()

	if vote == nil {
		return errors.Wrapf(model.ErrInvalidVoteMessage, "VoteMessage is nil")
	}

	hash := string(vote.GetBlockHash())
	pub := string(vote.GetSignature().GetPubkey())
	if _, ok := lock.findedVote[hash+pub]; ok {
		return errors.Wrapf(ErrAlreadyAddVoteMessage, "already add vote: %#v", vote)
	}

	// === add vote ===
	if len(lock.votedQueue) >= cap(lock.votedQueue) { // shifts Limits
		delete(lock.acceptedCounter, string(lock.findedVote[lock.votedQueue[0]].GetBlockHash()))
		delete(lock.findedVote, lock.votedQueue[0])
		lock.votedQueue = lock.votedQueue[1:]
	}
	lock.findedVote[hash+pub] = vote
	lock.votedQueue = append(lock.votedQueue, hash+pub)
	lock.acceptedCounter[hash]++
	// ================

	lock.checkAndLock(hash)
	return nil
}

func (lock *LockOnMemory) GetLockedProposal(height int64) (model.Proposal, bool) {
	if ret, ok := lock.lockedProposal[height]; ok {
		return ret, true
	}
	return nil, false
}

func (lock *LockOnMemory) Clean(height int64) {
	for k, _ := range lock.lockedProposal {
		if k < height {
			delete(lock.lockedProposal, k)
		}
	}
}
