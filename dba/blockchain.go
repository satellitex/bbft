package dba

import (
	"database/sql"
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/crypto"
	"github.com/satellitex/bbft/model/proto"
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

type BlockChainMock struct {
	db        map[int64]model.Block
	hashIndex map[crypto.Hash]int64
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

func NewBlockChainMock() BlockChain {
	return &BlockChainMock{
		make(map[int64]model.Block),
		make(map[crypto.Hash]int64),
		0,
		new(sync.Mutex),
	}
}

func (b *BlockChainMock) getIndex(block model.Block) (int64, bool) {
	hashPtr, err := block.GetHash()
	if err != nil {
		return -1, false
	}
	id, ok := b.hashIndex[*hashPtr]
	if ok {
		return id, true
	}
	return -1, false
}

func (b *BlockChainMock) Top() (model.Block, bool) {
	defer b.m.Unlock()
	b.m.Lock()

	res, ok := b.db[b.counter-1]
	if !ok {
		return nil, false
	}
	return res, true
}

func (b *BlockChainMock) Push(block model.Block) error {
	defer b.m.Unlock()
	b.m.Lock()

	if _, ok := b.getIndex(block); !ok {
		hashPtr, err := block.GetHash()
		if err != nil {
			return errors.Wrapf(ErrBlockChainPush, err.Error())
		}
		b.hashIndex[*hashPtr] = b.counter
		b.db[b.counter] = block
		b.counter += 1
	} else {
		return errors.Wrapf(ErrBlockChainPush,
			"Already exist block %x", block.GetHash())
	}
	return nil
}
