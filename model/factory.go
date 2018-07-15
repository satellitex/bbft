package model

import "github.com/pkg/errors"

var (
	ErrNewBlock    = errors.Errorf("Failed Factory NewBlock")
	ErrNewProposal = errors.Errorf("Failed Factory NewProposal")
)

type ModelFactory interface {
	NewBlock(height int64, preBlockHash []byte, createdTime int64, txs []Transaction) (Block, error)
	NewProposal(block Block, round int32) (Proposal, error)
	NewVoteMessage(hash []byte) VoteMessage
	NewSignature(pubkey []byte, signature []byte) Signature
	NewPeer(address string, pubkey []byte) Peer
}

type Hasher interface {
	GetHash() ([]byte, error)
}

func MustGetHash(hasher Hasher) []byte {
	hash, err := hasher.GetHash()
	if err != nil {
		panic("Unexpected GetHash : " + err.Error())
	}
	return hash
}
