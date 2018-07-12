package dba

import (
	"github.com/satellitex/bbft/model"
	"sync"
)

type ProposalTxQueue interface {
	Push(tx model.Transaction) error
	Pop() (model.Transaction, bool)
}

type ProposalTxQueueOnMemory struct {
	mutex *sync.Mutex
}

func NewProposalTxQueueOnMemory() ProposalTxQueue {
	return &ProposalTxQueueOnMemory{
		new(sync.Mutex),
	}
}

func (q *ProposalTxQueueOnMemory) Push(tx model.Transaction) error {
	return nil
}

func (q *ProposalTxQueueOnMemory) Pop() (model.Transaction, bool) {
	return nil, false
}
