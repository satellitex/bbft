package model

type ModelFactory interface {
	NewBlock(height int64, preBlockHash []byte, createdTime int64, txs []Transaction, signature Signature) Block
	NewProposal(block Block, round int64) Proposal
	NewVoteMessage(hash []byte, signature Signature) VoteMessage
	NewSignature(pubkey []byte, signature []byte) Signature
}
