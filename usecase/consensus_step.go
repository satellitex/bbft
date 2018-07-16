package usecase

import (
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/config"
	"github.com/satellitex/bbft/dba"
	"github.com/satellitex/bbft/model"
)

type ConsensusStep interface {
	Run()
	Propose() error
	Vote() error
	PreCommit() error
}

// [Height][Round] = Proposal を管理するgi
type ProposalFinder struct {
	field map[int64]map[int32]model.Proposal
}

func NewProposalFinder() *ProposalFinder {
	return &ProposalFinder{
		make(map[int64]map[int32]model.Proposal),
	}
}

func (f *ProposalFinder) Find(height int64, round int32) (model.Proposal, bool) {
	if fie, ok := f.field[height]; ok {
		if p, ok := fie[round]; ok {
			return p, true
		}
	}
	return nil, false
}

func (f *ProposalFinder) Set(proposal model.Proposal) error {
	if proposal == nil {
		return errors.Wrap(model.ErrInvalidProposal, "proposal is nil")
	}
	height, round := proposal.GetBlock().GetHeader().GetHeight(), proposal.GetRound()
	if _, ok := f.field[height]; !ok {
		f.field[height] = make(map[int32]model.Proposal)
	}
	f.field[height][round] = proposal
	return nil
}

func (f *ProposalFinder) Clear(height int64) {
	for key, _ := range f.field {
		if key < height {
			delete(f.field, key)
		}
	}
}

// PreCommit を管理する
//
// PreCommit が 2/3 以上集まった時、collectedHash を上書きする。
// Get() 時に collectedHash が存在した場合、PreCommit が 2/3以上集まっているので Commit Phase に遷移する。
// その後、取得した collectedHash は再度 nil にする。
type PreCommitFinder struct {
	collectedHash []byte
	field         map[string]int
	queue         []string
	limit         int
	ps            dba.PeerService
}

func NewPreCommitFinder(ps dba.PeerService, conf config.BBFTConfig) *PreCommitFinder {
	return &PreCommitFinder{
		nil,
		make(map[string]int),
		make([]string, 0, conf.PreCommitFinderLimits),
		conf.PreCommitFinderLimits,
		ps,
	}
}

func (f *PreCommitFinder) Get() ([]byte, bool) {
	if f.collectedHash == nil {
		return nil, false
	}
	ret := f.collectedHash
	f.collectedHash = nil
	return ret, true
}

func (f *PreCommitFinder) Set(vote model.VoteMessage) error {
	if vote == nil {
		return errors.Wrapf(model.ErrInvalidVoteMessage, "vote is nil")
	}

	hashStr := string(vote.GetBlockHash())
	if _, ok := f.field[hashStr]; !ok {
		if len(f.queue) >= f.limit {
			delete(f.field, f.queue[0])
			f.queue = f.queue[1:]
		}
		f.queue = append(f.queue, hashStr)
	}
	f.field[hashStr]++

	if f.ps.GetNumberOfRequiredAcceptPeers() <= f.field[hashStr] {
		f.collectedHash = vote.GetBlockHash()
		f.field[hashStr] = math.MinInt32
	}
	return nil
}

type ConsensusStepUsecase struct {
	bc                 dba.BlockChain
	lock               dba.Lock
	queue              dba.ProposalTxQueue
	sender             model.ConsensusSender
	statelessValidator model.StatelessValidator
	statefulValidator  model.StatefulValidator
	channel            *ReceiveChannel
	factory            model.ModelFactory
}

var (
	ErrConsensusProposal  = errors.Errorf("Failed This peer Proposal")
	ErrConsensusVote      = errors.Errorf("Failed This peer Vote")
	ErrConsensusPreCommit = errors.Errorf("Failed This peer PreCommit")
	ErrConsensusCommit    = errors.Errorf("Failed This peer ConsensusCommit")
)

func (c *ConsensusStepUsecase) Run() {
	for {
		for {
			c.Propose()
			c.Vote()
			c.PreCommit()
		}
		c.Commit()
	}
}

func (c *ConsensusStepUsecase) Propose() error {
	return nil
}

func (c *ConsensusStepUsecase) Vote() error {
	return nil
}

func (c *ConsensusStepUsecase) PreCommit() error {
	return nil
}

func (c *ConsensusStepUsecase) Commit() error {
	proposal, ok := c.lock.GetLockedProposal(0)
	if !ok {
		return errors.Wrapf(ErrConsensusCommit,
			"Not Founbd Locked Proposal")
	}
	block := proposal.GetBlock()
	c.bc.Commit(block)
	return nil
}
