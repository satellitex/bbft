package dba_test

import (
	"github.com/pkg/errors"
	. "github.com/satellitex/bbft/dba"
	"github.com/satellitex/bbft/model"
	. "github.com/satellitex/bbft/test_utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"sync"
)

func testReceiverPoolOnMemory(t *testing.T, pool ReceiverPool) {

	t.Run("success set propagete tx, and check isExist", func(t *testing.T) {
		obj := RandomValidTx(t)
		err := pool.SetPropagate(obj)
		assert.NoError(t, err)
		// ok case
		assert.True(t, pool.IsExistPropagate(obj))
		// false case
		assert.False(t, pool.IsExistPropagate(RandomValidTx(t)))
	})

	t.Run("failed set propagete nil", func(t *testing.T) {
		err := pool.SetPropagate(nil)
		assert.EqualError(t, errors.Cause(err), model.ErrInvalidTransaction.Error())
	})

	t.Run("success set propose proposal, and check isExist", func(t *testing.T) {
		obj := RandomProposal(t)
		err := pool.SetPropose(obj)
		assert.NoError(t, err)
		// ok case
		assert.True(t, pool.IsExistPropose(obj))
		// false case
		assert.False(t, pool.IsExistPropose(RandomProposal(t)))
	})

	t.Run("failed set propagete nil", func(t *testing.T) {
		err := pool.SetPropose(nil)
		assert.EqualError(t, errors.Cause(err), model.ErrInvalidProposal.Error())
	})

	t.Run("success set vote voteMessage, and check isExist", func(t *testing.T) {
		obj := RandomVoteMessage(t)
		err := pool.SetVote(obj)
		assert.NoError(t, err)
		// ok case
		assert.True(t, pool.IsExistVote(obj))
		// false case
		assert.False(t, pool.IsExistVote(RandomVoteMessage(t)))
	})

	t.Run("failed set vote nil", func(t *testing.T) {
		err := pool.SetVote(nil)
		assert.EqualError(t, errors.Cause(err), model.ErrInvalidVoteMessage.Error())
	})

	t.Run("success set preCommit voteMessage, and check isExist", func(t *testing.T) {
		obj := RandomVoteMessage(t)
		err := pool.SetPreCommit(obj)
		assert.NoError(t, err)
		// ok case
		assert.True(t, pool.IsExistPreCommit(obj))
		// false case
		assert.False(t, pool.IsExistPreCommit(RandomVoteMessage(t)))
	})

	t.Run("failed set preCommit nil", func(t *testing.T) {
		err := pool.SetPreCommit(nil)
		assert.EqualError(t, errors.Cause(err), model.ErrInvalidVoteMessage.Error())
	})

	t.Run("success paralell set", func(t *testing.T) {
		waiter := &sync.WaitGroup{}
		for i := 0; i < GetTestConfig().ReceivePropagateTxPoolLimits*2; i++ {
			waiter.Add(2)
			tx := RandomValidTx(t)
			go func() {
				assert.NoError(t, pool.SetPropagate(tx))
				waiter.Done()
			}()
			go func() {
				pool.IsExistPropagate(tx)
				waiter.Done()
			}()
		}
		waiter.Wait()
	})
}

func TestReceiverPoolOnMemory(t *testing.T) {
	pool := NewReceiverPoolOnMemory(GetTestConfig())
	testReceiverPoolOnMemory(t, pool)
}
