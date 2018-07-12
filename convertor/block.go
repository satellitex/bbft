package convertor

import (
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
)

type Block struct {
	*bbft.Block
}

type Proposal struct {
	*bbft.Proposal
}

type BlockHeader struct {
	*bbft.Block_Header
}

func (b *Block) GetHeader() model.BlockHeader {
	return &BlockHeader{b.Header}
}

func (b *Block) GetTransactions() []model.Transaction {
	ret := make([]model.Transaction, len(b.Transactions))
	for id, tx := range b.Transactions {
		ret[id] = &Transaction{tx}
	}
	return ret
}

func (b *Block) GetSignature() model.Signature {
	return &Signature{b.Signature}
}

func (b *Block) GetHash() ([]byte, error) {
	//TODO tx.GetHash() は payload の hash なので signature を含んでない,毎回 sha256計算したほうが一気にやるよりはやそう？
	result, err := b.GetHeader().GetHash()
	if err != nil {
		return nil, err
	}
	for _, tx := range b.GetTransactions() {
		hash, err := tx.GetHash()
		if err != nil {
			return nil, err
		}
		result = append(result, hash...)
	}
	return CalcHash(result), nil
}

func (b *Block) Verify() bool {
	hash, err := b.GetHash()
	if err != nil {
		return false
	}
	return Verify(b.Signature.Pubkey, hash, b.Signature.Signature)
}

func (h *BlockHeader) GetHash() ([]byte, error) {
	return CalcHashFromProto(h)
}

func (p *Proposal) GetBlock() model.Block {
	return &Block{p.Block}
}
