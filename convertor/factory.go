package convertor

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
)

type ModelFactory struct{}

var (
	ErrModelFactoryNewBlock       = errors.New("Failed Create model.Block")
	ErrModelFactoryNewProposal    = errors.New("Failed Create model.Proposal")
	ErrModelFactoryNewVoteMessage = errors.New("Failed Create model.VoteMessage")
)

func NewModelFactory() model.ModelFactory {
	return &ModelFactory{}
}

func (_ *ModelFactory) NewBlock(height int64, preBlockHash []byte, createdTime int64, txs []model.Transaction, signature model.Signature) (model.Block, error) {
	sig, ok := signature.(*Signature)
	if !ok {
		return nil, errors.Wrapf(ErrModelFactoryNewBlock,
			"Can not cast Signature model: %#v.", signature)
	}

	ptxs := make([]*bbft.Transaction, len(txs))
	for id, tx := range txs {
		tmp, ok := tx.(*Transaction)
		if !ok {
			return nil, errors.Wrapf(ErrModelFactoryNewBlock,
				"Can not cast Transaction model: %#v.", tx)
		}
		ptxs[id] = tmp.Transaction
	}
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
	}, nil
}

func (_ *ModelFactory) NewProposal(block model.Block, round int64) (model.Proposal, error) {
	b, ok := block.(*Block)
	if !ok {
		return nil, errors.Wrapf(ErrModelFactoryNewProposal,
			"Can not cast Block model: %#v.", block)
	}
	return &Proposal{
		&bbft.Proposal{
			Block: b.Block,
			Round: round,
		},
	}, nil
}

func (_ *ModelFactory) NewVoteMessage(hash []byte, signature model.Signature) (model.VoteMessage, error) {
	sigtmp, ok := signature.(*Signature)
	if !ok {
		return nil, errors.Wrapf(ErrModelFactoryNewVoteMessage,
			"Can not cast Signature model: %#v.", signature)
	}
	return &VoteMessage{
		&bbft.VoteMessage{
			BlockHash: hash,
			Signature: sigtmp.Signature,
		},
	}, nil
}

func (_ *ModelFactory) NewSignature(pubkey []byte, signature []byte) model.Signature {
	return &Signature{
		&bbft.Signature{
			Pubkey:    pubkey,
			Signature: signature,
		},
	}
}

type TxModelBuilder struct {
	*Transaction
}

func NewTxModelBuilder() *TxModelBuilder {
	return &TxModelBuilder{
		&Transaction{
			&bbft.Transaction{
				Payload:    &bbft.Transaction_Payload{},
				Signatures: make([]*bbft.Signature, 0, 5),
			},
		}}
}

func (b *TxModelBuilder) Signature(sig model.Signature) *TxModelBuilder {
	signature, _ := sig.(*Signature)
	b.Signatures = append(b.Signatures, signature.Signature)
	return b
}

func (b *TxModelBuilder) Sign(pubkey []byte, privateKey []byte) *TxModelBuilder {
	hash, err := b.GetHash()
	if err != nil {
		fmt.Printf("Error Sign : %s", err.Error())
	}
	b.Signatures = append(b.Signatures,
		&bbft.Signature{
			Pubkey:    pubkey,
			Signature: Sign(privateKey, hash),
		})
	return b
}

func (b *TxModelBuilder) Message(msg string) *TxModelBuilder {
	b.Payload.Todo = msg
	return b
}

func (b *TxModelBuilder) build() model.Transaction {
	return b.Transaction
}
