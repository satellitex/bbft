package convertor_test

import (
	"github.com/pkg/errors"
	. "github.com/satellitex/bbft/convertor"
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
		assert.EqualError(t, errors.Cause(slv.Validate(block)), ErrStatelessValidate.Error())
	})
	t.Run("failed invalid key and valid block", func(t *testing.T) {
		block := ValidErrSignedBlock(t)
		assert.EqualError(t, errors.Cause(slv.Validate(block)), ErrStatelessValidate.Error())
	})
	t.Run("failed invalid key and invalid block", func(t *testing.T) {
		block := InvalidErrSignedBlock(t)
		assert.EqualError(t, errors.Cause(slv.Validate(block)), ErrStatelessValidate.Error())
	})
}
