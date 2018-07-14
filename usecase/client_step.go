package usecase

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
)

type ClientGateReceiver interface {
	Gate(tx model.Transaction) error
}

type ClientGateReceiverUsecase struct {
	slv model.StatelessValidator
	factory   model.ModelFactory
	sender    model.ConsensusSender
}

func NewClientGateReceiverUsecase(validator model.StatelessValidator, factory model.ModelFactory, sender model.ConsensusSender) ClientGateReceiver {
	return &ClientGateReceiverUsecase{
		slv: validator,
		factory:   factory,
		sender:    sender,
	}
}

func (c *ClientGateReceiverUsecase) Gate(tx model.Transaction) error {
	if tx == nil {
		return errors.Wrapf(model.ErrInvalidTransaction, "tx is nil")
	}

	if err := tx.Verify(); err != nil {
		return errors.Wrapf(model.ErrTransactionVerify, err.Error())
	}

	err := c.sender.Propagate(tx)
	if err != nil {
		return errors.Wrapf(model.ErrConsensusSenderPropagate, err.Error())
	}
	return nil
}
