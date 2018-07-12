package convertor

import (
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
)

type Block struct {
	*bbft.Block
}

func ConvertBlock(block *bbft.Block) *Block {
	return &Block{block}
}

type Proposal struct {
	*bbft.Proposal
}

func ConvertProposal(proposal *bbft.Proposal) *Proposal {
	return &Proposal{proposal}
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
	// TODO TransactionとHeraderを分割してHashする
	return CalcHashFromProto(b)
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
