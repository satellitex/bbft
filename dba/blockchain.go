package dba

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
	"sync"
)

type BlockChain interface {
	Top() (model.Block, bool)
	FindTx(hash []byte) (model.Transaction, bool)
	// Commit is allowed only Commitable Block, ohterwise panic
	Commit(block model.Block)
	VerifyCommit(block model.Block) error
}

type BlockChainOnMemory struct {
	db        map[int64]model.Block
	tx        map[string]model.Transaction
	hashIndex map[string]int64
	counter   int64
	m         *sync.Mutex
}

func NewBlockChainOnMemory() BlockChain {
	return &BlockChainOnMemory{
		make(map[int64]model.Block),
		make(map[string]model.Transaction),
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
	ErrBlockChainVerifyCommitInvalidHeight       = errors.New("Failed Invalid Height of Block")
	ErrBlockChainVerifyCommitInvalidPreBlockHash = errors.New("Failed Invalid PreBlockHash of Block")
	ErrBlockChainVerifyCommitInvalidCreatedTime  = errors.New("Failed Invalid CreatedTime of Block")
	ErrBlockChainVerifyCommitAlreadyExist        = errors.New("Failed Alraedy Exist Block")
	ErrBlockChainVerifyCommit                    = errors.New("Failed Blockchain Verify Commit")
)

func (b *BlockChainOnMemory) VerifyCommit(block model.Block) error {
	b.m.Lock()
	defer b.m.Unlock()

	if block == nil {
		return errors.Wrapf(model.ErrInvalidBlock, "block is nil")
	}

	// Height Check
	if height := block.GetHeader().GetHeight(); height != b.counter {
		return errors.Wrapf(ErrBlockChainVerifyCommitInvalidHeight, "height: %d, expected %d", height, b.counter)
	}

	top, ok := b.Top()
	// First Commit is always OK
	if ok {
		// Must PreBlockHash == top.Hash
		if preHash := block.GetHeader().GetPreBlockHash(); !bytes.Equal(preHash, model.MustGetHash(top)) {
			return errors.Wrapf(ErrBlockChainVerifyCommitInvalidPreBlockHash,
				"block preBlockHash is not valid\npreBlockHash: %x\nexpected: %x\n",
				preHash, model.MustGetHash(top))
		}
		// Must createdTime > top.createdTime
		if createdTime := block.GetHeader().GetCreatedTime(); createdTime <= top.GetHeader().GetCreatedTime() {
			return errors.Wrapf(ErrBlockChainVerifyCommitInvalidCreatedTime,
				"block CreatedTime is not valid\ncreatedTime: %d\npreBlockCreatedTime: %d",
				createdTime, top.GetHeader().GetCreatedTime())
		}
		// Already exist check
		hash, err := block.GetHash()
		if err != nil {
			return errors.Wrapf(model.ErrBlockGetHash, err.Error())
		}
		if id, ok := b.getIndex(hash); ok {
			return errors.Wrapf(ErrBlockChainVerifyCommitAlreadyExist,
				"Already exist block %x is %d-th Block", model.MustGetHash(block), id)
		}
	}
	return nil
}

func (b *BlockChainOnMemory) Commit(block model.Block) {
	b.m.Lock()
	defer b.m.Unlock()

	if block == nil {
		panic("commit block is nil")
	}
	b.hashIndex[string(model.MustGetHash(block))] = b.counter
	b.db[b.counter] = block
	b.counter += 1

	for _, tx := range block.GetTransactions() {
		if tx == nil {
			panic("commit transaction is nil")
		}
		b.tx[string(model.MustGetHash(tx))] = tx
	}
}

func (b *BlockChainOnMemory) FindTx(hash []byte) (model.Transaction, bool) {
	b.m.Lock()
	defer b.m.Unlock()

	tx, ok := b.tx[string(hash)]
	if !ok {
		return nil, false
	}
	return tx, true
}
