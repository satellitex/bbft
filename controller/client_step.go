package controller

import (
	"github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/proto"
	"github.com/satellitex/bbft/usecase"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ClientGateController struct {
	receiver usecase.ClientGateReceiver
	author   *convertor.Author
}

func NewClientGateController(receiver usecase.ClientGateReceiver, author *convertor.Author) *ClientGateController {
	return &ClientGateController{
		receiver: receiver,
		author:   author,
	}
}

func (c *ClientGateController) Write(ctx context.Context, tx *bbft.Transaction) (*bbft.TxResponse, error) {
	transaction := &convertor.Transaction{tx}

	ctx, err := c.author.ProtoAurhorize(ctx, transaction)
	if err != nil {
		return nil, err
	}

	err = c.receiver.Gate(transaction)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &bbft.TxResponse{}, nil
}
