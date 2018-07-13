package model

import "github.com/pkg/errors"

var ErrInvalidSignature = errors.Errorf("Failed Invalid Signature")

type Signature interface {
	GetPubkey() []byte
	GetSignature() []byte
}
