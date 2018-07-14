package model

import "github.com/pkg/errors"

var (
	ErrVoteMessageGetBlockHash = errors.Errorf("Failed VoteMessage GetBlockHash")
	ErrVoteMessageVerify       = errors.Errorf("Failed VoteMessage Verify")
	ErrVoteMessageSign         = errors.Errorf("Failed VoteMessage Sign")
	ErrInvalidVoteMessage      = errors.Errorf("Failed Invalid VoteMessage")
)

type VoteMessage interface {
	GetBlockHash() []byte
	GetSignature() Signature
	Sign(pubKey []byte, privKey []byte) error
	Verify() error
}
