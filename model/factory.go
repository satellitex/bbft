package model

type ModelFactory interface {
	NewBlock(height int64, preBlockHash []byte, createdTime int64, txs []Transaction) (Block, error)
	NewProposal(block Block, round int64) (Proposal, error)
	NewVoteMessage(hash []byte) VoteMessage
	NewSignature(pubkey []byte, signature []byte) Signature
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
