package test_utils

import (
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func StatusCodeEquals(err error, code codes.Code) bool {
	status, ok := status.FromError(err)
	if !ok {
		return false
	}

	return status.Code() == code
}

func ValidateStatusCode(t *testing.T, err error, code codes.Code) {
	require.Error(t, err)
	if !StatusCodeEquals(err, code) {
		t.Errorf("Validate Status Code Error %v, but want %v", err, code)
	}
}
