package convertor

import (
	"github.com/satellitex/bbft/model/proto"
	"github.com/satellitex/bbft/proto"
)

type Block struct {
	*bbft.Block
}

type BlockHeader struct {
	*bbft.Block_Header
}

func (b *Block) GetHeader() proto.BlockHeader {
	return &BlockHeader{b.Header}
}
