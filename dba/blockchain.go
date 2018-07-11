package dba

import (
	"database/sql"
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
	"sync"
)

var (
	ErrBlockChainPush = errors.New("Failed Blockchain Push")
)

type BlockChain interface {
	Top() (model.Block, bool)
	Push(block model.Block) error
}

type BlockChainSQL struct {
	db *sql.DB
}

type BlockChainOnMemory struct {
	db        map[int64]model.Block
	hashIndex map[string]int64
	counter   int64
	m         *sync.Mutex
}

func NewBlockChainSQL(db *sql.DB) BlockChain {
	return &BlockChainSQL{db}
}

func (b *BlockChainSQL) Top() (model.Block, bool) {
	return nil, false
}

func (b *BlockChainSQL) Push(block model.Block) error {
	return nil
}

func NewBlockChainOnMemory() BlockChain {
	return &BlockChainOnMemory{
		make(map[int64]model.Block),
		make(map[string]int64),
		0,
		new(sync.Mutex),
	}
}

func (b *BlockChainOnMemory) getIndex(block model.Block) (int64, bool) {
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

func (b *BlockChainOnMemory) Push(block model.Block) error {
	defer b.m.Unlock()
	b.m.Lock()

	if _, ok := b.getIndex(block); !ok {
		hash, err := block.GetHash()
		if err != nil {
			return errors.Wrapf(ErrBlockChainPush, err.Error())
		}
		b.hashIndex[string(hash)] = b.counter
		b.db[b.counter] = block
		b.counter += 1
	} else {
		return errors.Wrapf(ErrBlockChainPush,
			"Already exist block %x", block.GetHash())
	}
	return nil
}
