package convertor

import (
	"github.com/satellitex/bbft/crypto"
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
)

type Transaction struct {
	*bbft.Transaction
}

type TransactionPayload struct {
	*bbft.Transaction_Payload
}

func (t *Transaction) GetHash() (crypto.HashPtr, error) {
	return crypto.CalcHashFromProto(t)
}

func (t *Transaction) GetPayload() model.TransactionPayload {
	return &TransactionPayload{t.Payload}
}

func (p *TransactionPayload) GetMessage() string {
	return p.Todo
}
