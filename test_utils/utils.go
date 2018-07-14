package test_utils

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"
	"testing"
)

func MultiErrorCheck(t *testing.T, err error, expectedErr error) {
	multiErr := multierr.Errors(err)
	expectedMultiErr := multierr.Errors(expectedErr)

	for id, e := range multiErr {
		assert.EqualError(t, errors.Cause(e), expectedMultiErr[id].Error())
	}
}

func MultiErrorInCheck(t *testing.T, err error, expectedErr error) {
	multiErr := multierr.Errors(err)
	multiCausedErr := make([]error, 0, len(multiErr))

	for _, e := range multiErr {
		multiCausedErr = append(multiCausedErr, errors.Cause(e))
	}
	assert.Contains(t, multiCausedErr, expectedErr)
}
