package model

type ModelFactory interface {
	NewBlock(height int64, preBlockHash []byte, createdTime int64, txs []Transaction, signature Signature) (Block, error)
	NewProposal(block Block, round int64) (Proposal, error)
	NewVoteMessage(hash []byte, signature Signature) (VoteMessage, error)
	NewSignature(pubkey []byte, signature []byte) Signature
}
