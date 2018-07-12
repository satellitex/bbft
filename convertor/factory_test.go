package convertor

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	rnd "math/rand"
	"strconv"
	"testing"
)

func randomStr() string {
	return strconv.FormatUint(rnd.Uint64(), 36)
}

func randomTx(t *testing.T) model.Transaction {
	tx, err := NewTxModelBuilder().
		Message(randomStr()).
		Signature(randomSig()).
		build()
	require.NoError(t, err)
	return tx
}

func randomTxs(t *testing.T) []model.Transaction {
	txs := make([]model.Transaction, 30)
	for id, _ := range txs {
		txs[id] = randomTx(t)
	}
	return txs
}

func randomByte() []byte {
	b, _ := NewKeyPair()
	return b
}

func randomSig() model.Signature {
	pub, sig := NewKeyPair()
	return (&ModelFactory{}).NewSignature(pub, sig)
}

type Hasher interface {
	GetHash() ([]byte, error)
}

func getHash(t *testing.T, hasher Hasher) []byte {
	hash, err := hasher.GetHash()
	require.NoError(t, err)
	return hash
}

func randomBlock(t *testing.T) model.Block {
	block, err := NewModelFactory().NewBlock(rnd.Int63(), randomByte(), rnd.Int63(), randomTxs(t), randomSig())
	require.NoError(t, err)
	return block
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
			randomTxs(t),
			randomSig(),
		},
		{
			"case 2",
			nil,
			999999999999,
			[]byte(""),
			0,
			randomTxs(t),
			randomSig(),
		},
		{
			"signature nil case",
			ErrModelFactoryNewBlock,
			0,
			nil,
			999999999999,
			randomTxs(t),
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
		{
			"minus number is no problem case",
			nil,
			-1,
			nil,
			-1,
			nil,
			randomSig(),
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			block, err := NewModelFactory().NewBlock(c.expectedHeight, c.expectedHash, c.expectedCreatedTime, c.expectedTxs, c.expectedSig)
			if c.expectedError != nil {
				assert.EqualError(t, errors.Cause(err), c.expectedError.Error())
				return
			}
			assert.NoError(t, err)
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

func TestProposalFactory(t *testing.T) {
	for _, c := range []struct {
		name          string
		expectedError error
		expectedBlock model.Block
		expectedRound int64
	}{
		{
			"case 1",
			nil,
			randomBlock(t),
			rnd.Int63(),
		},
		{
			"case 2",
			nil,
			randomBlock(t),
			rnd.Int63(),
		},
		{
			"case 3",
			nil,
			randomBlock(t),
			rnd.Int63(),
		},
		{
			"block nil case",
			ErrModelFactoryNewProposal,
			nil,
			rnd.Int63(),
		},
		{
			"round -1 case",
			nil,
			randomBlock(t),
			-1,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			proposal, err := NewModelFactory().NewProposal(c.expectedBlock, c.expectedRound)
			if c.expectedError != nil {
				assert.EqualError(t, errors.Cause(err), c.expectedError.Error())
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, getHash(t, c.expectedBlock), getHash(t, proposal.GetBlock()))
			assert.Equal(t, c.expectedRound, proposal.GetRound())
		})
	}

}

func TestVoteMessageFactory(t *testing.T) {
	for _, c := range []struct {
		name          string
		expectedError error
		expectedHash  []byte
		expectedSig   model.Signature
	}{
		{
			"case 1",
			nil,
			randomByte(),
			randomSig(),
		},
		{
			"case 2",
			nil,
			randomByte(),
			randomSig(),
		},
		{
			"case 3",
			nil,
			randomByte(),
			randomSig(),
		},
		{
			"case 4",
			nil,
			randomByte(),
			randomSig(),
		},
		{
			"sig nil case",
			ErrModelFactoryNewVoteMessage,
			randomByte(),
			nil,
		},
		{
			"hash nil case no problem",
			nil,
			nil,
			randomSig(),
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			vote, err := NewModelFactory().NewVoteMessage(c.expectedHash, c.expectedSig)
			if c.expectedError != nil {
				assert.EqualError(t, errors.Cause(err), c.expectedError.Error())
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, c.expectedHash, vote.GetBlockHash())
			assert.Equal(t, c.expectedSig.GetPubkey(), vote.GetSignature().GetPubkey())
			assert.Equal(t, c.expectedSig.GetSignature(), vote.GetSignature().GetSignature())
		})
	}

}

func TestSignatureFactory(t *testing.T) {
	for _, c := range []struct {
		name        string
		expectedPub []byte
		expectedSig []byte
	}{
		{
			"case 1",
			randomByte(),
			randomByte(),
		},
		{
			"case 2",
			randomByte(),
			randomByte(),
		},
		{
			"case 3",
			randomByte(),
			randomByte(),
		},
		{
			"case 4",
			randomByte(),
			randomByte(),
		},
		{
			"case 5",
			randomByte(),
			randomByte(),
		},
		{
			"pub nil case no problem",
			nil,
			randomByte(),
		},
		{
			"sig nil case no problem",
			randomByte(),
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			sig := NewModelFactory().NewSignature(c.expectedPub, c.expectedSig)
			assert.Equal(t, c.expectedPub, sig.GetPubkey())
			assert.Equal(t, c.expectedSig, sig.GetSignature())
		})
	}
}

func TestTxModelBuilder(t *testing.T) {
	validPub, validPriv := NewKeyPair()
	for _, c := range []struct {
		name              string
		expectedError     error
		expectedMessage   string
		expectedSignature model.Signature
		expectedPubkey    []byte
		expectedPrivKey   []byte
	}{
		{
			"case 1",
			nil,
			randomStr(),
			randomSig(),
			validPub,
			validPriv,
		},
		{
			"case 2",
			nil,
			randomStr(),
			randomSig(),
			validPub,
			validPriv,
		},
		{
			"case 3",
			nil,
			randomStr(),
			randomSig(),
			validPub,
			validPriv,
		},
		{
			"empty string case is valid",
			nil,
			"",
			randomSig(),
			validPub,
			validPriv,
		},
		{
			"signature nil case",
			ErrTxModelBuild,
			randomStr(),
			nil,
			validPub,
			validPriv,
		},
		{
			"pubkey nil case",
			ErrTxModelBuild,
			randomStr(),
			randomSig(),
			nil,
			validPriv,
		},
		{
			"privkey nil case",
			ErrTxModelBuild,
			randomStr(),
			randomSig(),
			validPub,
			nil,
		},
		{
			"all ng case",
			ErrTxModelBuild,
			randomStr(),
			nil,
			nil,
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			tx, err := NewTxModelBuilder().Message(c.expectedMessage).
				Signature(c.expectedSignature).
				Sign(c.expectedPubkey, c.expectedPrivKey).
				build()
			if c.expectedError != nil {
				assert.EqualError(t, errors.Cause(err), c.expectedError.Error())
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, c.expectedMessage, tx.GetPayload().GetMessage())
			assert.Equal(t, c.expectedSignature.GetPubkey(), tx.GetSignatures()[0].GetPubkey())
			assert.Equal(t, c.expectedSignature.GetSignature(), tx.GetSignatures()[0].GetSignature())
			assert.Equal(t, c.expectedPubkey, tx.GetSignatures()[1].GetPubkey())
			signature, err := Sign(c.expectedPrivKey, getHash(t, tx))
			require.NoError(t, err)
			assert.Equal(t, signature, tx.GetSignatures()[1].GetSignature())
		})
	}
}
