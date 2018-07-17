package usecase

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/config"
	"github.com/satellitex/bbft/dba"
	"github.com/satellitex/bbft/model"
	"math"
	"time"
)

type ConsensusStep interface {
	Run()
	Propose(height int64, round int32) error
	Vote(height int64, round int32) error
	PreCommit(height int64, round int32) error
	Commit(height int64, round int32) error
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

func NewPreCommitFinder(ps dba.PeerService, conf *config.BBFTConfig) *PreCommitFinder {
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

func UnixTime(t time.Time) int64 {
	return t.UnixNano()
}

func Now() int64 {
	return UnixTime(time.Now())
}

type ConsensusStepUsecase struct {
	conf    *config.BBFTConfig
	bc      dba.BlockChain
	ps      dba.PeerService
	lock    dba.Lock
	queue   dba.ProposalTxQueue
	sender  model.ConsensusSender
	slv     model.StatelessValidator
	sfv     model.StatefulValidator
	factory model.ModelFactory
	channel *ReceiveChannel

	proposalFinder    *ProposalFinder
	preCommitFinder   *PreCommitFinder
	thisRoundProposal model.Proposal
	roundStartTime    time.Duration
	roundCommitTime   time.Duration
	proposeTimeOut    time.Duration
	voteStartTimeOut  time.Duration
	preCommitTimeOut  time.Duration
}

var (
	ErrConsensusProposal  = errors.Errorf("Failed This peer Proposal")
	ErrConsensusVote      = errors.Errorf("Failed This peer Vote")
	ErrConsensusPreCommit = errors.Errorf("Failed This peer PreCommit")
	ErrConsensusCommit    = errors.Errorf("Failed This peer ConsensusCommit")
)

// Runnning Consensus Endless...
func (c *ConsensusStepUsecase) Run() {
	for {
		top, ok := c.bc.Top()
		if !ok {
			panic("Unexpected Error No BlockChain Top")
		}
		c.roundStartTime = time.Duration(top.GetHeader().GetCommitTime())
		height, round := top.GetHeader().GetHeight(), int32(-1)
		for {
			round++

			// each Phase TimeOut Calc
			c.proposeTimeOut = c.roundStartTime + c.conf.ProposeMaxCalcTime + c.conf.AllowedConnectDelayTime
			c.voteStartTimeOut = c.proposeTimeOut + c.conf.VoteMaxCalcTime + c.conf.AllowedConnectDelayTime
			c.preCommitTimeOut = c.voteStartTimeOut + c.conf.PreCommitMaxCalcTime + c.conf.AllowedConnectDelayTime
			c.roundCommitTime = c.preCommitTimeOut + c.conf.CommitMaxCalcTime

			if err := c.Propose(height, round); err != nil {
				fmt.Printf("Consensus ProposePhase Error!! height:%d, round:%d\n%s", height, round, err.Error())
			}

			if err := c.Vote(height, round); err != nil {
				fmt.Printf("Consensus VotePhase Error!! height:%d, round:%d\n%s", height, round, err.Error())
			}

			if err := c.PreCommit(height, round); err != nil {
				fmt.Printf("Consensus VotePhase Error!! height:%d, round:%d\n%s", height, round, err.Error())
			} else {
				break
			}
			c.roundStartTime = c.preCommitTimeOut
		}
		c.Commit(height, round)
	}
}

/*
  if Lock.emtpy():
    if Leader is me:
      N <- number of transactions to take out.
      cnt <- 0
      for tx in ProposalTxQueue:
        if validate(tx):
          Block.Transactions <- tx
          cnt ++
        if N <= cnt:
          break
      Block.Header.Height <- BlockChain.Top.Header.Height
      Block.Header.Hash <- BlockChain.Top.Hash
      Block.Header.CreatedTime <- now
      Block.Signature <- sign( sum_hash(hash(Block.Header),hash(sum_hash(Block.Transactions))), myself.privateKey )
      Proposal[Round].Block <- Block
      Proposal[Round].Round <- Round
      send_all_peers( Propose, Proposal[Round] )
      Vote( Height, Round )
    else:
      wait until receiveProposal(Round=Round) or ProposeTimeOut or not Lock.empty()
      Proposal[Round] <- receivedProposal(Round=Round)
  Vote( Height, Round )
*/
func (c *ConsensusStepUsecase) Propose(height int64, round int32) error {
	if _, ok := c.lock.GetLockedProposal(height); !ok {
		if bytes.Equal(c.ps.GetPermutationPeers(height)[round].GetPubkey(), c.conf.PublicKey) {
			// Leader is me
			txs := make([]model.Transaction, 0, c.conf.NumberOfBlockHasTransactions)
			for len(txs) < c.conf.NumberOfBlockHasTransactions {
				tx, ok := c.queue.Pop()
				if !ok { // ProposalTx is empty
					break
				}
				if err := c.slv.TxValidate(tx); err != nil {
					continue
				}
				txs = append(txs, tx)
			}
			top, ok := c.bc.Top()
			if !ok {
				return errors.New("Unexpected Error No BlockChain Top")
			}
			block, err := c.factory.NewBlock(height, model.MustGetHash(top), int64(c.roundCommitTime), txs)
			if err != nil {
				return err
			}
			block.Sign(c.conf.PublicKey, c.conf.SecretKey)
			propsal, err := c.factory.NewProposal(block, round)
			if err != nil {
				return err
			}

			if err = c.sender.Propose(propsal); err != nil {
				return err
			}
		} else {
			// Leader is not me
			timer := time.NewTimer(c.proposeTimeOut - time.Duration(Now()))
			select {
			case <-timer.C:
				break
			case proposal := <-c.channel.Propose:
				c.proposalFinder.Set(proposal)
				if c.thisRoundProposal, ok = c.proposalFinder.Find(height, round); ok {
					break
				}
			case <-c.channel.Vote:
				if _, ok := c.lock.GetLockedProposal(height); ok {
					break
				}
			case preCommit := <-c.channel.PreCommit:
				c.preCommitFinder.Set(preCommit)
			}
		}
	}
	return nil
}

func (c *ConsensusStepUsecase) Vote(height int64, round int32) error {
	return nil
}

func (c *ConsensusStepUsecase) PreCommit(height int64, round int32) error {
	return nil
}

func (c *ConsensusStepUsecase) Commit(height int64, round int32) error {
	proposal, ok := c.lock.GetLockedProposal(0)
	if !ok {
		return errors.Wrapf(ErrConsensusCommit,
			"Not Founbd Locked Proposal")
	}
	block := proposal.GetBlock()
	c.bc.Commit(block)
	return nil
}
