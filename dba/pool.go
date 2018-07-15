package dba

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/config"
	"github.com/satellitex/bbft/model"
	"strconv"
	"sync"
)

var (
	ErrHasherPoolSet   = errors.New("Failed Hasher Pool Set")
	ErrReceiverPoolSet = errors.New("Failed Reveiver Pool Set")
)

type hasherPool interface {
	set(hasher model.Hasher) error
	isExist(hasher model.Hasher) bool
}

type hasherPoolOnMemory struct {
	mp    map[string]model.Hasher
	q     []string
	limit int
	mutex *sync.Mutex
}

func newHasherPoolOnMemory(limit int) hasherPool {
	return &hasherPoolOnMemory{make(map[string]model.Hasher), make([]string, 0, limit), limit, new(sync.Mutex)}
}

func (t *hasherPoolOnMemory) set(hasher model.Hasher) error {
	defer t.mutex.Unlock()
	t.mutex.Lock()

	if hasher == nil {
		return errors.New("set object is nil")
	}
	hash, err := hasher.GetHash()
	if err != nil {
		return errors.New(err.Error())
	}
	t.mp[string(hash)] = hasher
	t.q = append(t.q, string(hash))
	if len(t.q) >= t.limit {
		delete(t.mp, t.q[0])
		t.q = t.q[1:]
	}
	return nil
}

func (t *hasherPoolOnMemory) isExist(hasher model.Hasher) bool {
	if hasher == nil {
		return false
	}
	hash, err := hasher.GetHash()
	if err != nil {
		return false
	}
	if _, ok := t.mp[string(hash)]; ok {
		return true
	}
	return false
}

type ReceiverPool interface {
	SetPropagate(tx model.Transaction) error
	SetPropose(proposal model.Proposal) error
	SetVote(vote model.VoteMessage) error
	SetPreCommit(preCommit model.VoteMessage) error

	IsExistPropagate(tx model.Transaction) bool
	IsExistPropose(proposal model.Proposal) bool
	IsExistVote(vote model.VoteMessage) bool
	IsExistPreCommit(preCommit model.VoteMessage) bool
}

type proposalHasher struct {
	hash []byte
}

func newProposalHasher(proposal model.Proposal) model.Hasher {
	return &proposalHasher{
		[]byte(strconv.FormatInt(proposal.GetBlock().GetHeader().GetHeight(), 16) + "::" + strconv.FormatInt(int64(proposal.GetRound()), 16)),
	}
}

func (p *proposalHasher) GetHash() ([]byte, error) {
	return p.hash, nil
}

type voteHasher struct {
	hash []byte
}

func newVoteHasher(vote model.VoteMessage) model.Hasher {
	return &voteHasher{append(vote.GetBlockHash(), vote.GetSignature().GetPubkey()...)}
}

func (v *voteHasher) GetHash() ([]byte, error) {
	return v.hash, nil
}

type ReceiverPoolOnMemory struct {
	txPool        hasherPool
	proposePool   hasherPool
	votePool      hasherPool
	preCommitPool hasherPool
}

func NewReceiverPoolOnMemory(conf *config.BBFTConfig) ReceiverPool {
	return &ReceiverPoolOnMemory{
		newHasherPoolOnMemory(conf.ReceivePropagateTxPoolLimits),
		newHasherPoolOnMemory(conf.ReceiveProposeProposalPoolLimits),
		newHasherPoolOnMemory(conf.ReceiveVoteVoteMessagePoolLimits),
		newHasherPoolOnMemory(conf.ReceivePreCommitVoteMessagePoolLimits),
	}
}

func (r *ReceiverPoolOnMemory) SetPropagate(tx model.Transaction) error {
	if tx == nil {
		return errors.Wrapf(model.ErrInvalidTransaction, "tx is nil")
	}
	if err := r.txPool.set(tx); err != nil {
		return errors.Wrapf(ErrHasherPoolSet, err.Error())
	}
	return nil
}

func (r *ReceiverPoolOnMemory) SetPropose(proposal model.Proposal) error {
	if proposal == nil {
		return errors.Wrapf(model.ErrInvalidProposal, "proposal is nil")
	}
	hasher := newProposalHasher(proposal)
	if err := r.proposePool.set(hasher); err != nil {
		return errors.Wrapf(ErrHasherPoolSet, err.Error())
	}
	return nil
}

func (r *ReceiverPoolOnMemory) SetVote(vote model.VoteMessage) error {
	if vote == nil {
		return errors.Wrapf(model.ErrInvalidVoteMessage, "vote is nil")
	}
	hasher := newVoteHasher(vote)
	if err := r.votePool.set(hasher); err != nil {
		return errors.Wrapf(ErrHasherPoolSet, err.Error())
	}
	return nil
}

func (r *ReceiverPoolOnMemory) SetPreCommit(preCommit model.VoteMessage) error {
	if preCommit == nil {
		return errors.Wrapf(model.ErrInvalidVoteMessage, "preCommit is nil")
	}
	hasher := newVoteHasher(preCommit)
	if err := r.preCommitPool.set(hasher); err != nil {
		return errors.Wrapf(ErrHasherPoolSet, err.Error())
	}
	return nil
}

func (r *ReceiverPoolOnMemory) IsExistPropagate(tx model.Transaction) bool {
	return r.txPool.isExist(tx)
}

func (r *ReceiverPoolOnMemory) IsExistPropose(proposal model.Proposal) bool {
	if proposal == nil {
		return false
	}
	return r.proposePool.isExist(newProposalHasher(proposal))
}

func (r *ReceiverPoolOnMemory) IsExistVote(vote model.VoteMessage) bool {
	if vote == nil {
		return false
	}
	return r.votePool.isExist(newVoteHasher(vote))
}

func (r *ReceiverPoolOnMemory) IsExistPreCommit(preCommit model.VoteMessage) bool {
	if preCommit == nil {
		return false
	}
	return r.preCommitPool.isExist(newVoteHasher(preCommit))
}
