package convertor

import (
	"github.com/satellitex/bbft/model/proto"
	"github.com/satellitex/bbft/proto"
)

type Transaction struct {
	*bbft.Transaction
}

type TransactionPayload struct {
	*bbft.Transaction_Payload
}

func (t *Transaction) GetPayload() proto.TransactionPayload {
	return &TransactionPayload{t.Payload}
}

func (p *TransactionPayload) GetMessage() string {
	return p.Todo
}
