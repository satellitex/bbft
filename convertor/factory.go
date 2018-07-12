package convertor

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
	"go.uber.org/multierr"
)

type ModelFactory struct{}

var (
	ErrModelFactoryNewBlock       = errors.New("Failed Create model.Block")
	ErrModelFactoryNewProposal    = errors.New("Failed Create model.Proposal")
	ErrModelFactoryNewVoteMessage = errors.New("Failed Create model.VoteMessage")
	ErrTxModelBuild               = errors.Errorf("Failed build TxModelBuildeer has error")
)

func NewModelFactory() model.ModelFactory {
	return &ModelFactory{}
}

func (_ *ModelFactory) NewBlock(height int64, preBlockHash []byte, createdTime int64, txs []model.Transaction) (model.Block, error) {
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

func (_ *ModelFactory) NewVoteMessage(hash []byte) model.VoteMessage {
	return &VoteMessage{
		&bbft.VoteMessage{BlockHash: hash},
	}
}

func (_ *ModelFactory) NewSignature(pubkey []byte, signature []byte) model.Signature {
	return NewSignature(pubkey, signature)
}

type TxModelBuilder struct {
	*Transaction
	err error
}

func NewTxModelBuilder() *TxModelBuilder {
	return &TxModelBuilder{
		&Transaction{
			&bbft.Transaction{
				Payload:    &bbft.Transaction_Payload{},
				Signatures: make([]*bbft.Signature, 0, 32),
			},
		}, nil}
}

// Test 用 Verifyしない
func (b *TxModelBuilder) Signature(sig model.Signature) *TxModelBuilder {
	signature, ok := sig.(*Signature)
	if !ok {
		b.err = multierr.Append(b.err, errors.Errorf("Failed Signature: Can not cast Signature model: %#v.", sig))
		return b
	}
	b.Signatures = append(b.Signatures, signature.Signature)
	return b
}

func (b *TxModelBuilder) Sign(pubkey []byte, privateKey []byte) *TxModelBuilder {
	hash, err := b.GetHash()
	if err != nil {
		b.err = multierr.Append(b.err, errors.Errorf("Failed Sign: %s", err.Error()))
	}
	signature, err := Sign(privateKey, hash)
	if err != nil {
		b.err = multierr.Append(b.err, errors.Errorf("Failed Sign: %s", err.Error()))
	}
	if !Verify(pubkey, hash, signature) {
		b.err = multierr.Append(b.err, errors.Errorf("Failed Sign: Can not verify comb with pubKey and privateKey"))
	}
	b.Signatures = append(b.Signatures,
		&bbft.Signature{
			Pubkey:    pubkey,
			Signature: signature,
		})
	return b
}

func (b *TxModelBuilder) Message(msg string) *TxModelBuilder {
	b.Payload.Todo = msg
	return b
}

func (b *TxModelBuilder) build() (model.Transaction, error) {
	if b.err != nil {
		return nil, errors.Wrapf(ErrTxModelBuild, b.err.Error())
	}
	return b.Transaction, nil
}
