package dba_test

import (
	. "github.com/satellitex/bbft/dba"
	. "github.com/satellitex/bbft/test_utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
)

func testLock_RegisterProposal(t *testing.T, lock Lock) {
	t.Run("success register valid random proposal", func(t *testing.T) {
		err := lock.RegisterProposal(RandomProposal(t))
		assert.NoError(t, err)
	})

	t.Run("failed register nil", func(t *testing.T) {
		err := lock.RegisterProposal(nil)
		assert.EqualError(t, errors.Cause(err), model.ErrInvalidProposal.Error())
	})
}

func TestLockOnMemory_RegisterProposal(t *testing.T) {
	lock := NewLockOnMemory(NewPeerServiceOnMemory())
	testLock_RegisterProposal(t, lock)
}
