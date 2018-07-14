package convertor

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
)

var ErrVoteMessageVerify = errors.Errorf("Failed VoteMessage Verify")

type VoteMessage struct {
	*bbft.VoteMessage
}

func (v *VoteMessage) GetSignature() model.Signature {
	return &Signature{v.Signature}
}

func (v *VoteMessage) Sign(pubKey []byte, privKey []byte) error {
	signature, err := Sign(privKey, v.GetBlockHash())
	if err != nil {
		return errors.Wrapf(ErrCryptoSign, err.Error())
	}
	if err := Verify(pubKey, v.GetBlockHash(), signature); err != nil {
		return errors.Wrapf(ErrCryptoVerify, err.Error())
	}
	v.Signature = &bbft.Signature{Pubkey: pubKey, Signature: signature}
	return nil
}

func (v *VoteMessage) Verify() error {
	if v.Signature == nil {
		return errors.Wrapf(model.ErrInvalidSignature, "VoteMessage.Signature is nil")
	}
	if err := Verify(v.Signature.Pubkey, v.GetBlockHash(), v.Signature.Signature); err != nil {
		return errors.Wrapf(ErrCryptoVerify, err.Error())
	}
	return nil
}
