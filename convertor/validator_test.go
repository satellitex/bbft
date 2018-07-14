package convertor_test

import (
	"github.com/pkg/errors"
	. "github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/model"
	. "github.com/satellitex/bbft/test_utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStatelessValidator_Validate(t *testing.T) {
	slv := NewStatelessValidator()
	t.Run("success valid key and valid txs", func(t *testing.T) {
		block := ValidSignedBlock(t)
		assert.NoError(t, slv.Validate(block))
	})
	t.Run("failed valid key and inValid txs", func(t *testing.T) {
		block := InvalidSingedBlock(t)
		MultiErrorInCheck(t, slv.Validate(block), model.ErrTransactionVerify)
	})
	t.Run("failed invalid key and valid block", func(t *testing.T) {
		block := ValidErrSignedBlock(t)
		MultiErrorInCheck(t, slv.Validate(block), model.ErrBlockVerify)
	})
	t.Run("failed invalid key and invalid block", func(t *testing.T) {
		block := InvalidErrSignedBlock(t)
		MultiErrorInCheck(t, errors.Cause(slv.Validate(block)), model.ErrTransactionVerify)
	})
	t.Run("failed nil block", func(t *testing.T) {
		assert.EqualError(t, errors.Cause(slv.Validate(nil)), model.ErrInvalidBlock.Error())
	})
}
