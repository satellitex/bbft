package client

import (
	"github.com/satellitex/bbft/dba/queue"
	"github.com/satellitex/bbft/model"
)

type ClientStepReceiver interface {
	Gate(tx model.Transaction) error
	Propagate(ptx model.ProposalTx) error
}

type ClientStepReceiverUsecase struct {
	queue     queue.ProposalTxQueue
	sender    model.ConsensusSender
	validator model.StatelessValidator
}

func (c *ClientStepReceiverUsecase) Gate(tx model.Transaction) error {
	return nil
}

func (c *ClientStepReceiverUsecase) Propagete(ptx model.Transaction) error {
	return nil
}

func (c *ClientStepReceiverUsecase) propagate() error {
	return nil
}

func (c *ClientStepReceiverUsecase) insertQueue() error {
	return nil
}
