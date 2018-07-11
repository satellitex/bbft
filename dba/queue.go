package dba

import (
	"github.com/satellitex/bbft/model"
	"sync"
)

type ProposalTxQueue interface {
	Push(tx model.ProposalTx) error
	Pop() (model.ProposalTx, bool)
}

type ProposalTxQueueOnMemory struct {
	mutex *sync.Mutex
}

func NewProposalTxQueueOnMemory() ProposalTxQueue {
	return &ProposalTxQueueOnMemory{
		new(sync.Mutex),
	}
}

func (q *ProposalTxQueueOnMemory) Push(tx model.ProposalTx) error {
	return nil
}

func (q *ProposalTxQueueOnMemory) Pop() (model.ProposalTx, bool) {
	return nil, false
}
