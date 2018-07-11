package controller

import (
	"github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/proto"
	"github.com/satellitex/bbft/usecase"
	"golang.org/x/net/context"
)

type ClientGateController struct {
	receiver usecase.ClientGateReceiver
}

func (c *ClientGateController) Write(_ context.Context, tx *bbft.Transaction) (*bbft.TxResponse, error) {
	transaction := &convertor.Transaction{tx}
	err := c.receiver.Gate(transaction)
	if err != nil {
		return nil, err
	}
	return &bbft.TxResponse{}, nil
}
