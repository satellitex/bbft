package convertor

import (
	"github.com/satellitex/bbft/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	rnd "math/rand"
	"strconv"
	"testing"
)

func randomTx() model.Transaction {
	return NewTxModelBuilder().Message(strconv.FormatUint(rnd.Uint64(), 36)).Signature(randomSig()).build()
}

func randomTxs() []model.Transaction {
	txs := make([]model.Transaction, 30)
	for id, _ := range txs {
		txs[id] = randomTx()
	}
	return txs
}

func randomSig() model.Signature {
	pub, sig := NewKeyPair()
	return (&ModelFactory{}).NewSignature(pub, sig)
}

func getHash(t *testing.T, transaction model.Transaction) []byte {
	hash, err := transaction.GetHash()
	require.NoError(t, err)
	return hash
}

func TestBlockFactory(t *testing.T) {
	for _, c := range []struct {
		name                string
		expectedError       error
		expectedHeight      int64
		expectedHash        []byte
		expectedCreatedTime int64
		expectedTxs         []model.Transaction
		expectedSig         model.Signature
	}{
		{
			"case 1",
			nil,
			10,
			[]byte("preBlockHash"),
			5,
			randomTxs(),
			randomSig(),
		},
		{
			"case 2",
			nil,
			999999999999,
			[]byte(""),
			0,
			randomTxs(),
			randomSig(),
		},
		{
			"signature nil case",
			ErrModelFactoryNewBlock,
			0,
			nil,
			999999999999,
			randomTxs(),
			nil,
		},
		{
			"tx nil case",
			ErrModelFactoryNewBlock,
			100,
			nil,
			111,
			make([]model.Transaction, 2),
			randomSig(),
		},
		{
			"txs nil case",
			nil,
			100,
			nil,
			111,
			nil,
			randomSig(),
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			block, err := NewModelFactory().NewBlock(c.expectedHeight, c.expectedHash, c.expectedCreatedTime, c.expectedTxs, c.expectedSig)
			if c.expectedError != nil {
				assert.Error(t, err)
				return
			} else {
				assert.NoError(t, err)
			}
			for id, tx := range block.GetTransactions() {
				assert.Equal(t, getHash(t, c.expectedTxs[id]), getHash(t, tx))
			}
			assert.Equal(t, c.expectedHeight, block.GetHeader().GetHeight())
			assert.Equal(t, c.expectedCreatedTime, block.GetHeader().GetCreatedTime())
			assert.Equal(t, c.expectedSig.GetSignature(), block.GetSignature().GetSignature())
			assert.Equal(t, c.expectedSig.GetPubkey(), block.GetSignature().GetPubkey())
		})
	}

}
