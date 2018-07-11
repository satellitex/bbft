package client

import (
	"github.com/satellitex/bbft/dba/queue"
	"github.com/satellitex/bbft/model"
)

type ClientStepUsecase struct {
	tx     model.ProposalTx
	queue  queue.ProposalTxQueue
	sender model.ConsensusSender
	validator model.StatelessValidator
}

func (c *ClientStepUsecase) Compute() error {
	return nil
}

func (c *ClientStepUsecase) propagate() error {
	return nil
}

func (c *ClientStepUsecase) insertQueue() error {
	return nil
}
