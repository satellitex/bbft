package model

import "github.com/pkg/errors"

var (
	ErrInvalidBlock = errors.Errorf("Failed Invalid Block")
	ErrBlockGetHash = errors.Errorf("Failed Block GetHash")
	ErrBlockVerify  = errors.Errorf("Failed Block Verify")
	ErrBlockSign    = errors.Errorf("Failed Block Sign")

	ErrInvalidBlockHeader = errors.Errorf("Failed Invalid BlockHeader")
	ErrBlockHeaderGetHash = errors.Errorf("Failed BlockHeader GetHash")

	ErrInvalidProposal = errors.Errorf("Failed Invalid Proposal")
)

type Block interface {
	GetHeader() BlockHeader
	GetTransactions() []Transaction
	GetSignature() Signature
	GetHash() ([]byte, error)
	Verify() error
	Sign(pubKey []byte, privKey []byte) error
}

type BlockHeader interface {
	GetHeight() int64
	GetPreBlockHash() []byte
	GetCreatedTime() int64
	GetHash() ([]byte, error)
}

type Proposal interface {
	GetBlock() Block
	GetRound() int64
}
