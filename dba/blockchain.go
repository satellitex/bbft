package dba

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
	"sync"
)

type BlockChain interface {
	Top() (model.Block, bool)
	Commit(block model.Block) error
	VerifyCommit(block model.Block) error
}

type BlockChainOnMemory struct {
	db        map[int64]model.Block
	hashIndex map[string]int64
	counter   int64
	m         *sync.Mutex
}

func NewBlockChainOnMemory() BlockChain {
	return &BlockChainOnMemory{
		make(map[int64]model.Block),
		make(map[string]int64),
		0,
		new(sync.Mutex),
	}
}

func (b *BlockChainOnMemory) getIndex(hash []byte) (int64, bool) {
	id, ok := b.hashIndex[string(hash)]
	if ok {
		return id, true
	}
	return -1, false
}

func (b *BlockChainOnMemory) Top() (model.Block, bool) {
	res, ok := b.db[b.counter-1]
	if !ok {
		return nil, false
	}
	return res, true
}

var (
	ErrBlockChainVerifyCommitVerifyBlock = errors.New("Failed Verify Block")
	ErrBlockChainVerifyCommitHeightCheck = errors.New("Failed Block Height is not valid")
	ErrBlockChainVerifyCommit            = errors.New("Failed Blockchain Verify Commit")
	//TODO
)

func (b *BlockChainOnMemory) VerifyCommit(block model.Block) error {
	// Verify block
	if err := block.Verify(); err != nil {
		return errors.Wrapf(ErrBlockChainVerifyCommitVerifyBlock, err.Error())
	}

	// Height Check
	if height := block.GetHeader().GetHeight(); height != b.counter {
		return errors.Wrapf(ErrBlockChainVerifyCommitHeightCheck, "height: %d, expected %d", height, b.counter)
	}

	top, ok := b.Top()
	// First Commit is always OK
	if ok {
		// Must PreBlockHash == top.Hash
		if preHash := block.GetHeader().GetPreBlockHash(); !bytes.Equal(preHash, model.MustGetHash(top)) {
			return errors.Wrapf(ErrBlockChainVerifyCommit,
				"block preBlockHash is not valid\npreBlockHash: %x\nexpected: %x\n",
				preHash, model.MustGetHash(top))
		}
		// Must createdTime > top.createdTime
		if createdTime := block.GetHeader().GetCreatedTime(); createdTime <= top.GetHeader().GetCreatedTime() {
			return errors.Wrapf(ErrBlockChainVerifyCommit,
				"block CreatedTime is not valid\ncreatedTime: %d\npreBlockCreatedTime: %d",
				createdTime, top.GetHeader().GetCreatedTime())
		}
		// Already exist check
		if id, ok := b.getIndex(model.MustGetHash(block)); ok {
			return errors.Wrapf(ErrBlockChainVerifyCommit,
				"Already exist block %x is %d-th Block", model.MustGetHash(block), id)
		}
	}
	return nil
}

func (b *BlockChainOnMemory) Commit(block model.Block) error {
	defer b.m.Unlock()
	b.m.Lock()
	if err := b.VerifyCommit(block); err != nil {
		return errors.Wrapf(ErrBlockChainVerifyCommit, err.Error())
	}

	b.hashIndex[string(model.MustGetHash(block))] = b.counter
	b.db[b.counter] = block
	b.counter += 1
	return nil
}
