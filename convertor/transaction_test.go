package convertor

import (
	"github.com/satellitex/bbft/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func randomValidTx(t *testing.T) model.Transaction {
	validPub, validPriv := NewKeyPair()
	tx, err := NewTxModelBuilder().
		Message(randomStr()).
		Sign(validPub, validPriv).
		build()
	require.NoError(t, err)
	return tx
}

func TestTransaction_GetHash(t *testing.T) {
	txs := make([]model.Transaction, 50)
	for id, _ := range txs {
		txs[id] = randomTx(t)
	}
	for id, a := range txs {
		for jd, b := range txs {
			if id != jd {
				assert.NotEqual(t, getHash(t, a), getHash(t, b))
			} else {
				assert.Equal(t, getHash(t, a), getHash(t, b))
			}
		}
	}
}

func TestTransaction_Verfy(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tx := randomValidTx(t)
		assert.True(t, tx.Verify())
	})
	t.Run("failed", func(t *testing.T) {
		tx := randomTx(t)
		assert.False(t, tx.Verify())
	})
}