package convertor

import (
	"github.com/satellitex/bbft/crypto"
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
)

type Block struct {
	*bbft.Block
}

type Propsoal struct {
	*bbft.Proposal
}

type BlockHeader struct {
	*bbft.Block_Header
}

func (b *Block) GetHash() (crypto.HashPtr, error) {
	// TODO 各Transactionを分割でHashすることで高速化できる
	return crypto.CalcHashFromProto(b)
}

func (b *Block) GetHeader() model.BlockHeader {
	return &BlockHeader{b.Header}
}

func (h *BlockHeader) GetHash() (crypto.HashPtr, error) {
	return crypto.CalcHashFromProto(h)
}
