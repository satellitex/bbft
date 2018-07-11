package usecase

import (
	"github.com/satellitex/bbft/model"
)

type ClientGateReceiver interface {
	Gate(tx model.Transaction) error
}

type ClientGateReceiverUsecase struct {
	validator model.StatelessValidator
}

func (c *ClientGateReceiverUsecase) Gate(tx model.Transaction) error {
	return nil
}
