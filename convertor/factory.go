package convertor

import (
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
)

type ModelFactory struct{}

func (_ *ModelFactory) NewBlock(height int64, preBlockHash []byte, createdTime int64, txs []model.Transaction, signature model.Signature) model.Block {
	ptxs := make([]*bbft.Transaction, len(txs))
	for id, tx := range txs {
		tmp, _ := tx.(*Transaction)
		ptxs[id] = tmp.Transaction
	}
	sig, _ := signature.(*Signature)
	return &Block{
		&bbft.Block{
			Header: &bbft.Block_Header{
				Height:       height,
				PreBlockHash: preBlockHash,
				CreatedTime:  createdTime,
			},
			Transactions: ptxs,
			Signature:    sig.Signature,
		},
	}
}

func (_ *ModelFactory) NewProposal(block model.Block, round int64) model.Proposal {
	btmp, _ := block.(*Block)
	return &Proposal{
		&bbft.Proposal{
			Block: btmp.Block,
			Round: round,
		},
	}
	return nil
}

func (_ *ModelFactory) NewVoteMessage(hash []byte, signature model.Signature) model.VoteMessage {
	sigtmp, _ := signature.(*Signature)
	return &VoteMessage{
		&bbft.VoteMessage{
			BlockHash: hash,
			Signature: sigtmp.Signature,
		},
	}
}

func (_ *ModelFactory) NewSignature(pubkey []byte, signature []byte) model.Signature {
	return &Signature{
		&bbft.Signature{
			Pubkey:    pubkey,
			Signature: signature,
		},
	}
}
