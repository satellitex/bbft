package controller

import (
	"github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/proto"
	"github.com/satellitex/bbft/usecase/client"
	"golang.org/x/net/context"
)

type ClientGateController struct {
	receiver client.ClientGateReceiver
}

func (c *ClientGateController) Write(_ context.Context, tx *bbft.Transaction) (*bbft.TxResponse, error) {
	transaction := &convertor.Transaction{tx}
	err := c.receiver.Gate(transaction)
	if err != nil {
		return nil, err
	}
	return &bbft.TxResponse{}, nil
}
