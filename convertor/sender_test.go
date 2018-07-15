package convertor_test

import (
	"github.com/pkg/errors"
	. "github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/model"
	. "github.com/satellitex/bbft/test_utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testConseusSender_InputCheck(t *testing.T, sender model.ConsensusSender) {

	t.Run("success propagete", func(t *testing.T) {
		err := sender.Propagate(RandomValidTx(t))
		assert.NoError(t, err)
	})

	t.Run("success propose", func(t *testing.T) {
		err := sender.Propose(RandomProposal(t))
		assert.NoError(t, err)
	})

	t.Run("success vote", func(t *testing.T) {
		err := sender.Vote(RandomVoteMessage(t))
		assert.NoError(t, err)
	})

	t.Run("success precommit", func(t *testing.T) {
		err := sender.PreCommit(RandomVoteMessage(t))
		assert.NoError(t, err)
	})

	t.Run("failed propagete, input: nil", func(t *testing.T) {
		err := sender.Propagate(nil)
		assert.EqualError(t, errors.Cause(err), model.ErrInvalidTransaction.Error())
	})

	t.Run("failed propose, input: nil", func(t *testing.T) {
		err := sender.Propose(nil)
		assert.EqualError(t, errors.Cause(err), model.ErrInvalidProposal.Error())
	})

	t.Run("failed vote, input: nil", func(t *testing.T) {
		err := sender.Vote(nil)
		assert.EqualError(t, errors.Cause(err), model.ErrInvalidVoteMessage.Error())
	})

	t.Run("failed precommit, input: nil", func(t *testing.T) {
		err := sender.PreCommit(nil)
		assert.EqualError(t, errors.Cause(err), model.ErrInvalidVoteMessage.Error())
	})

}

func TestMockConsensusSender(t *testing.T) {
	sender := NewMockConsensusSender()
	testConseusSender_InputCheck(t, sender)
}
