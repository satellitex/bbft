package dba

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
	"sync"
)

var (
	ErrProposalTxQueueLimits = errors.Errorf("PropposalTxQueue run limit reached")
)

type ProposalTxQueue interface {
	Push(tx model.Transaction) error
	Pop() (model.Transaction, bool)
}

type ProposalTxQueueOnMemory struct {
	mutex *sync.Mutex
	limit int
	queue []model.Transaction
}

func NewProposalTxQueueOnMemory(limit int) ProposalTxQueue {
	return &ProposalTxQueueOnMemory{
		new(sync.Mutex),
		limit,
		make([]model.Transaction, 0, limit),
	}
}

func (q *ProposalTxQueueOnMemory) Push(tx model.Transaction) error {
	defer q.mutex.Unlock()
	q.mutex.Lock()

	if tx == nil {
		return errors.Wrapf(model.ErrInvalidTransaction, "push transaction is nil")
	}
	if len(q.queue) < q.limit {
		q.queue = append(q.queue, tx)
	} else {
		return errors.Wrapf(ErrProposalTxQueueLimits, "queue's max length: %d", q.limit)
	}
	return nil
}

func (q *ProposalTxQueueOnMemory) Pop() (model.Transaction, bool) {
	if len(q.queue) == 0 {
		return nil, false
	}
	front := q.queue[0]
	q.queue = q.queue[1:]
	return front, true
}
