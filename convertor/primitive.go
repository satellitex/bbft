package convertor

import (
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
)

type Signature struct {
	*bbft.Signature
}

func NewSignature(pubkey []byte, signature []byte) model.Signature {
	return &Signature{
		&bbft.Signature{
			Pubkey:    pubkey,
			Signature: signature,
		},
	}
}
