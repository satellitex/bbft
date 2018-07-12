package convertor

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVoteMessage_SignAndVerify(t *testing.T) {
	t.Run("success valid key and exist hash", func(t *testing.T) {
		validPub, validPri := NewKeyPair()
		vote := NewModelFactory().NewVoteMessage(randomByte())

		err := vote.Sign(validPub, validPri)
		require.NoError(t, err)

		assert.True(t, vote.Verify())
	})
	t.Run("success valid key and nil hash", func(t *testing.T) {
		validPub, validPri := NewKeyPair()
		vote := NewModelFactory().NewVoteMessage(nil)

		err := vote.Sign(validPub, validPri)
		require.NoError(t, err)

		assert.True(t, vote.Verify())
	})
	t.Run("failed invalid key and exist hash", func(t *testing.T) {
		invalid, _ := NewKeyPair()
		vote := NewModelFactory().NewVoteMessage(randomByte())

		err := vote.Sign(invalid, invalid)
		require.Error(t, err)

		assert.False(t, vote.Verify())
	})
	t.Run("failed invalid key and nil hash", func(t *testing.T) {
		invalid, _ := NewKeyPair()
		vote := NewModelFactory().NewVoteMessage(nil)

		err := vote.Sign(invalid, invalid)
		require.Error(t, err)

		assert.False(t, vote.Verify())
	})
}
