package dba

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
	"sync"
)

var (
	ErrBlockChainCommit = errors.New("Failed Blockchain Commit")
)

type BlockChain interface {
	Top() (model.Block, bool)
	Commit(block model.Block) error
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

func (b *BlockChainOnMemory) GetIndex(block model.Block) (int64, bool) {
	hash, err := block.GetHash()
	if err != nil {
		return -1, false
	}
	id, ok := b.hashIndex[string(hash)]
	if ok {
		return id, true
	}
	return -1, false
}

func (b *BlockChainOnMemory) Top() (model.Block, bool) {
	defer b.m.Unlock()
	b.m.Lock()

	res, ok := b.db[b.counter-1]
	if !ok {
		return nil, false
	}
	return res, true
}

func (b *BlockChainOnMemory) Commit(block model.Block) error {
	defer b.m.Unlock()
	b.m.Lock()

	if _, ok := b.GetIndex(block); !ok {
		hash, err := block.GetHash()
		if err != nil {
			return errors.Wrapf(ErrBlockChainCommit, err.Error())
		}
		b.hashIndex[string(hash)] = b.counter
		b.db[b.counter] = block
		b.counter += 1
	} else {
		return errors.Wrapf(ErrBlockChainCommit,
			"Already exist block %x", block.GetHash())
	}
	return nil
}
