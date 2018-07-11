package convertor

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
)

type Block struct {
	*bbft.Block
	c model.Cryptor
}

type Propsoal struct {
	*bbft.Proposal
}

type BlockHeader struct {
	*bbft.Block_Header
	c model.Cryptor
}

func (b *Block) GetHash() ([]byte, error) {
	// TODO TransactionとHeraderを分割してHashする
	return CalcHashFromProto(b, b.c)
}

func (b *Block) Verify() bool {
	hash, err := b.GetHash()
	if err != nil {
		return false
	}
	return b.c.Verify(b.Signature.Pubkey, hash, b.Signature.Signature)
}

func (b *Block) GetHeader() model.BlockHeader {
	return &BlockHeader{b.Header, b.c}
}

func (h *BlockHeader) GetHash() ([]byte, error) {
	return CalcHashFromProto(h, h.c)
}
