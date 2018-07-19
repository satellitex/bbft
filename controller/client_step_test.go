package controller_test

import (
	"context"
	. "github.com/satellitex/bbft/controller"
	"github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/dba"
	"github.com/satellitex/bbft/proto"
	. "github.com/satellitex/bbft/test_utils"
	"github.com/satellitex/bbft/usecase"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"testing"
)

func NewTestClientGateController(t *testing.T) *ClientGateController {
	ps := dba.NewPeerServiceOnMemory()
	sender := convertor.NewMockConsensusSender()
	receiver := usecase.NewClientGateReceiverUsecase(
		convertor.NewStatelessValidator(),
		sender,
	)
	author := convertor.NewAuthor(ps)
	return NewClientGateController(receiver, author)
}

func TestClientGateController_Write(t *testing.T) {
	ctrl := NewTestClientGateController(t)

	for _, c := range []struct {
		ctx  context.Context
		tx   *bbft.Transaction
		code codes.Code
	}{
		{
			context.TODO(),
			RandomValidTx(t).(*convertor.Transaction).Transaction,
			codes.OK,
		},
		{
			context.TODO(),
			nil,
			codes.InvalidArgument,
		},
		{
			context.TODO(),
			RandomInvalidTx(t).(*convertor.Transaction).Transaction,
			codes.InvalidArgument,
		},
	} {
		_, err := ctrl.Write(c.ctx, c.tx)
		if c.code != codes.OK {
			ValidateStatusCode(t, err, c.code)
		} else {
			assert.NoError(t, err)
		}
	}
}
