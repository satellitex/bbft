package convertor_test

import (
	"github.com/pkg/errors"
	. "github.com/satellitex/bbft/convertor"
	. "github.com/satellitex/bbft/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVoteMessage_SignAndVerify(t *testing.T) {
	t.Run("success valid key and exist hash", func(t *testing.T) {
		validPub, validPri := NewKeyPair()
		vote := NewModelFactory().NewVoteMessage(RandomByte())

		err := vote.Sign(validPub, validPri)
		require.NoError(t, err)

		assert.NoError(t, vote.Verify())
	})
	t.Run("success valid key and nil hash", func(t *testing.T) {
		validPub, validPri := NewKeyPair()
		vote := NewModelFactory().NewVoteMessage(nil)

		err := vote.Sign(validPub, validPri)
		require.NoError(t, err)

		assert.NoError(t, vote.Verify())
	})
	t.Run("failed invalid key and exist hash", func(t *testing.T) {
		invalid, _ := NewKeyPair()
		vote := NewModelFactory().NewVoteMessage(RandomByte())

		err := vote.Sign(invalid, invalid)
		require.Error(t, err)

		assert.EqualError(t, errors.Cause(vote.Verify()), ErrVoteMessageVerify.Error())
	})
	t.Run("failed invalid key and nil hash", func(t *testing.T) {
		invalid, _ := NewKeyPair()
		vote := NewModelFactory().NewVoteMessage(nil)

		err := vote.Sign(invalid, invalid)
		require.Error(t, err)

		assert.EqualError(t, errors.Cause(vote.Verify()), ErrVoteMessageVerify.Error())
	})
	t.Run("failed nil signature", func(t *testing.T) {
		vote := NewModelFactory().NewVoteMessage(nil)
		vote.(*VoteMessage).Signature = nil

		assert.EqualError(t, errors.Cause(vote.Verify()), ErrVoteMessageVerify.Error())
	})
}
