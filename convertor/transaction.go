package convertor

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
)

var ErrTransactionVerify = errors.Errorf("Failed Transaction Verify")

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

func (t *Transaction) Verify() error {
	hash, err := t.GetHash()
	if err != nil {
		return errors.Wrapf(ErrTransactionVerify, err.Error())
	}
	if len(t.Signatures) == 0 {
		return errors.Wrapf(ErrTransactionVerify, "Signature length is 0")
	}
	for i, signature := range t.Signatures {
		if signature == nil {
			return errors.Wrapf(ErrTransactionVerify, "%d-th Signature is nil", i)
		}
		if err := Verify(signature.Pubkey, hash, signature.Signature); err != nil {
			return errors.Wrapf(ErrTransactionVerify, err.Error())
		}
	}
	return nil
}

func (p *TransactionPayload) GetMessage() string {
	return p.Todo
}
