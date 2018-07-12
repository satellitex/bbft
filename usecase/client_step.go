package usecase

import (
	"github.com/satellitex/bbft/model"
)

type ClientGateReceiver interface {
	Gate(tx model.Transaction) error
}

type ClientGateReceiverUsecase struct {
	validator model.StatelessValidator
	factory   model.ModelFactory
}

func (c *ClientGateReceiverUsecase) Gate(tx model.Transaction) error {
	return nil
}
