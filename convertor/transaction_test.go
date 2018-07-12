package convertor

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransaction_GetHash(t *testing.T) {
	txs := make([]model.Transaction, 50)
	for id, _ := range txs {
		txs[id] = randomValidTx(t)
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
		assert.NoError(t, tx.Verify())
	})
	t.Run("failed invalid signature", func(t *testing.T) {
		tx := randomInvalidTx(t)
		assert.EqualError(t, errors.Cause(tx.Verify()), ErrTransactionVerify.Error())
	})
	t.Run("failed not signed", func(t *testing.T) {
		tx, err := NewTxModelBuilder().Message(randomStr()).build()
		require.NoError(t, err)
		assert.EqualError(t, errors.Cause(tx.Verify()), ErrTransactionVerify.Error())
	})
	t.Run("failed nil signature", func(t *testing.T) {
		tx, err := NewTxModelBuilder().Message(randomStr()).build()
		require.NoError(t, err)
		tx.(*Transaction).Signatures = make([]*bbft.Signature, 5)
		assert.EqualError(t, errors.Cause(tx.Verify()), ErrTransactionVerify.Error())
	})
}
