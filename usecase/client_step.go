package usecase

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
	"log"
)

type ClientGateReceiver interface {
	Gate(tx model.Transaction) error
}

type ClientGateReceiverUsecase struct {
	slv    model.StatelessValidator
	sender model.ConsensusSender
}

func NewClientGateReceiverUsecase(validator model.StatelessValidator, sender model.ConsensusSender) ClientGateReceiver {
	return &ClientGateReceiverUsecase{
		slv:    validator,
		sender: sender,
	}
}

func (c *ClientGateReceiverUsecase) Gate(tx model.Transaction) error {
	if err := c.slv.TxValidate(tx); err != nil { // InvalidArgument (code = 3)
		return errors.Wrapf(model.ErrStatelessTxValidate, err.Error())
	}
	err := c.sender.Propagate(tx)
	if err != nil {
		log.Println(model.ErrConsensusSenderPropagate, err)
	}
	return nil
}
