package convertor

import (
	"github.com/pkg/errors"
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
	//TODO 毎回 sha256計算したほうが一気にやるよりはやそう？
	result, err := b.GetHeader().GetHash()
	if err != nil {
		return nil, err
	}
	for _, tx := range b.GetTransactions() {
		proto, ok := tx.(*Transaction)
		if !ok {
			return nil, errors.Errorf("Can not cast Transaction model: %#v.", tx)
		}
		hash, err := CalcHashFromProto(proto)
		if err != nil {
			return nil, err
		}
		result = append(result, hash...)
	}
	return CalcHash(result), nil
}

func (b *Block) Verify() bool {
	for _, tx := range b.GetTransactions() {
		if !tx.Verify() {
			return false
		}
	}
	hash, err := b.GetHash()
	if err != nil {
		return false
	}
	return Verify(b.Signature.Pubkey, hash, b.Signature.Signature)
}

func (b *Block) Sign(pubKey []byte, privKey []byte) error {
	hash, err := b.GetHash()
	if err != nil {
		return err
	}
	signature, err := Sign(privKey, hash)
	if err != nil {
		return err
	}
	if !Verify(pubKey, hash, signature) {
		return errors.Errorf("Failed Verify")
	}
	b.Signature = &bbft.Signature{Pubkey: pubKey, Signature: signature}
	return nil
}

func (h *BlockHeader) GetHash() ([]byte, error) {
	return CalcHashFromProto(h)
}

func (p *Proposal) GetBlock() model.Block {
	return &Block{p.Block}
}
