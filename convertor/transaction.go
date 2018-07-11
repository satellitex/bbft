package convertor

import (
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
)

type Transaction struct {
	*bbft.Transaction
}

type TransactionPayload struct {
	*bbft.Transaction_Payload
}

func (t *Transaction) GetPayload() model.TransactionPayload {
	return &TransactionPayload{t.Payload}
}

func (t *Transaction) GetHash() ([]byte, error) {
	return CalcHashFromProto(t.Payload)
}

func (t *Transaction) GetSignatures() []model.Signature {
	ret := make([]model.Signature, len(t.Signatures))
	for i, sig := range t.Signatures {
		ret[i] = Signature{sig}
	}
	return ret
}

func (t *Transaction) Verify() bool {
	hash, err := t.GetHash()
	if err != nil {
		return false
	}
	for _, signature := range t.Signatures {
		if Verify(signature.Pubkey, hash, signature.Signature) == false {
			return false
		}
	}
	return true
}

func (p *TransactionPayload) GetMessage() string {
	return p.Todo
}
