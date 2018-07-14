package dba_test

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/config"
	. "github.com/satellitex/bbft/dba"
	"github.com/satellitex/bbft/model"
	. "github.com/satellitex/bbft/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func testProposalTxQueue(t *testing.T, queue ProposalTxQueue) {
	t.Run("Success, random valid tx push", func(t *testing.T) {
		// Invalid Tx no problem
		txs := []model.Transaction{
			RandomValidTx(t),
			RandomInvalidTx(t),
			RandomValidTx(t),
			RandomInvalidTx(t),
			RandomValidTx(t),
			RandomInvalidTx(t),
		}

		for _, tx := range txs {
			err := queue.Push(tx)
			assert.NoError(t, err)
		}

		for _, tx := range txs {
			front, ok := queue.Pop()
			assert.True(t, ok)
			assert.Equal(t, tx, front)
		}
	})

	t.Run("Failed, nil tx push", func(t *testing.T) {
		err := queue.Push(nil)
		assert.EqualError(t, errors.Cause(err), model.ErrInvalidTransaction.Error())
	})

	t.Run("Empty pop, return nil", func(t *testing.T) {
		front, ok := queue.Pop()
		assert.False(t, ok)
		assert.Equal(t, nil, front)
	})

	t.Run("Failed, over limits random valid tx push", func(t *testing.T) {
		limit := config.GetTestConfig().QueueLimits
		txs := make([]model.Transaction, 0, limit)
		for i := 0; i < limit; i++ {
			txs = append(txs, RandomValidTx(t))
			require.NoError(t, queue.Push(txs[i]))
		}

		err := queue.Push(RandomValidTx(t))
		assert.EqualError(t, errors.Cause(err), ErrProposalTxQueueLimits.Error())

		for i := 0; i < limit; i++ {
			front, ok := queue.Pop()
			assert.True(t, ok)
			assert.Equal(t, txs[i], front)
		}

		_, ok := queue.Pop()
		assert.False(t, ok)
	})
}

func TestProposalTxQueueOnMemory(t *testing.T) {
	config := config.GetTestConfig()

	queue := NewProposalTxQueueOnMemory(config.QueueLimits)
	testProposalTxQueue(t, queue)
}
