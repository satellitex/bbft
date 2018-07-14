package usecase_test

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/model"
	. "github.com/satellitex/bbft/test_utils"
	. "github.com/satellitex/bbft/usecase"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClientGateReceiverUsecase_Gate(t *testing.T) {
	validator := convertor.NewStatelessValidator()
	sender := convertor.NewMockConsensusSender()

	gate := NewClientGateReceiverUsecase(validator, sender)

	t.Run("success case", func(t *testing.T) {
		tx := RandomValidTx(t)
		err := gate.Gate(tx)
		assert.NoError(t, err)
		assert.Equal(t, sender.(*convertor.MockConsensusSender).Tx, tx)
	})

	t.Run("failed case", func(t *testing.T) {
		err := gate.Gate(RandomInvalidTx(t))
		assert.EqualError(t, errors.Cause(err), model.ErrStatelessTxValidate.Error())
	})
}
