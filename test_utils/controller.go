package test_utils

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/satellitex/bbft/config"
	"github.com/satellitex/bbft/convertor"
	"github.com/stretchr/testify/require"
	"go.uber.org/multierr"
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

func ValidContext(t *testing.T, conf *config.BBFTConfig, prt proto.Message) context.Context {
	ctx, err := convertor.NewContextByProtobufDebug(conf, prt)
	require.NoError(t, err)
	return ctx
}

func MultiValidateStatusCode(t *testing.T, err error, code codes.Code) {
	require.Error(t, err)
	multiErr := multierr.Errors(err)
	for _, e := range multiErr {
		ValidateStatusCode(t, e, code)
	}
}
