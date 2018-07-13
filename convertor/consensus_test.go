package convertor_test

import (
	"github.com/pkg/errors"
	. "github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/model"
	. "github.com/satellitex/bbft/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVoteMessage_Sign(t *testing.T) {
	t.Run("success valid key and exist hash", func(t *testing.T) {
		validPub, validPri := NewKeyPair()
		vote := NewModelFactory().NewVoteMessage(RandomByte())

		err := vote.Sign(validPub, validPri)
		assert.NoError(t, err)
	})
	t.Run("success valid key and nil hash", func(t *testing.T) {
		validPub, validPri := NewKeyPair()
		vote := NewModelFactory().NewVoteMessage(nil)

		err := vote.Sign(validPub, validPri)
		assert.NoError(t, err)
	})
	t.Run("failed invalid key and exist hash", func(t *testing.T) {
		invalid, _ := NewKeyPair()
		vote := NewModelFactory().NewVoteMessage(RandomByte())

		err := vote.Sign(invalid, invalid)
		assert.EqualError(t, errors.Cause(err), ErrCryptoSign.Error())
	})
	t.Run("failed invalid key and nil hash", func(t *testing.T) {
		invalid, _ := NewKeyPair()
		vote := NewModelFactory().NewVoteMessage(nil)

		err := vote.Sign(invalid, invalid)
		assert.Error(t, errors.Cause(err), ErrCryptoVerify.Error())
	})
	t.Run("failed invalid signed key", func(t *testing.T) {
		vote := NewModelFactory().NewVoteMessage(nil)

		err := vote.Sign(nil, nil)
		assert.Error(t, errors.Cause(err), ErrCryptoSign.Error())
	})
}

func TestVoteMessage_Verify(t *testing.T) {
	t.Run("failed nil signature", func(t *testing.T) {
		vote := NewModelFactory().NewVoteMessage(nil)
		vote.(*VoteMessage).Signature = nil

		assert.EqualError(t, errors.Cause(vote.Verify()), model.ErrInvalidSignature.Error())
	})
	t.Run("failed invalid Sign signature", func(t *testing.T) {
		invalid, _ := NewKeyPair()
		vote := NewModelFactory().NewVoteMessage(RandomByte())

		err := vote.Sign(invalid, invalid)
		require.Error(t, err)

		assert.EqualError(t, errors.Cause(vote.Verify()), ErrCryptoVerify.Error())
	})

}
