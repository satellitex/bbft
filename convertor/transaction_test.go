package convertor_test

import (
	"github.com/pkg/errors"
	. "github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
	. "github.com/satellitex/bbft/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransaction_GetHash(t *testing.T) {
	txs := make([]model.Transaction, 50)
	for id, _ := range txs {
		txs[id] = RandomValidTx(t)
	}
	for id, a := range txs {
		for jd, b := range txs {
			if id != jd {
				assert.NotEqual(t, GetHash(t, a), GetHash(t, b))
			} else {
				assert.Equal(t, GetHash(t, a), GetHash(t, b))
			}
		}
	}
}

func TestTransaction_Verfy(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tx := RandomValidTx(t)
		assert.NoError(t, tx.Verify())
	})
	t.Run("failed invalid signature", func(t *testing.T) {
		tx := RandomInvalidTx(t)
		assert.EqualError(t, errors.Cause(tx.Verify()), ErrCryptoVerify.Error())
	})
	t.Run("failed not signed", func(t *testing.T) {
		tx, err := NewTxModelBuilder().Message(RandomStr()).Build()
		require.NoError(t, err)
		assert.EqualError(t, errors.Cause(tx.Verify()), ErrInvalidSignatures.Error())
	})
	t.Run("failed nil signature", func(t *testing.T) {
		tx, err := NewTxModelBuilder().Message(RandomStr()).Build()
		require.NoError(t, err)
		tx.(*Transaction).Signatures = make([]*bbft.Signature, 5)
		assert.EqualError(t, errors.Cause(tx.Verify()), model.ErrInvalidSignature.Error())
	})
}
