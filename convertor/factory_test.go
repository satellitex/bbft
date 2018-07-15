package convertor_test

import (
	"github.com/pkg/errors"
	. "github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/model"
	. "github.com/satellitex/bbft/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/multierr"
	"math/rand"
	"testing"
)

func TestBlockFactory(t *testing.T) {
	for _, c := range []struct {
		name                string
		expectedError       error
		expectedHeight      int64
		expectedHash        []byte
		expectedCreatedTime int64
		expectedTxs         []model.Transaction
	}{
		{
			"case 1",
			nil,
			10,
			[]byte("preBlockHash"),
			5,
			RandomTxs(t),
		},
		{
			"case 2",
			nil,
			999999999999,
			[]byte(""),
			0,
			RandomTxs(t),
		},
		{
			"hash nil case no problem",
			nil,
			0,
			nil,
			999999999999,
			RandomTxs(t),
		},
		{
			"tx nil case",
			model.ErrInvalidTransaction,
			100,
			nil,
			111,
			make([]model.Transaction, 2),
		},
		{
			"txs nil case no problem",
			nil,
			100,
			nil,
			111,
			nil,
		},
		{
			"minus number is no problem case",
			nil,
			-1,
			nil,
			-1,
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			block, err := NewModelFactory().NewBlock(c.expectedHeight, c.expectedHash, c.expectedCreatedTime, c.expectedTxs)
			if c.expectedError != nil {
				assert.EqualError(t, errors.Cause(err), c.expectedError.Error())
				return
			}
			assert.NoError(t, err)
			for id, tx := range block.GetTransactions() {
				assert.Equal(t, GetHash(t, c.expectedTxs[id]), GetHash(t, tx))
			}
			assert.Equal(t, c.expectedHeight, block.GetHeader().GetHeight())
			assert.Equal(t, c.expectedCreatedTime, block.GetHeader().GetCreatedTime())
		})
	}

}

func TestProposalFactory(t *testing.T) {
	for _, c := range []struct {
		name          string
		expectedError error
		expectedBlock model.Block
		expectedRound int32
	}{
		{
			"case 1",
			nil,
			RandomBlock(t),
			rand.Int31(),
		},
		{
			"case 2",
			nil,
			RandomBlock(t),
			rand.Int31(),
		},
		{
			"case 3",
			nil,
			RandomBlock(t),
			rand.Int31(),
		},
		{
			"block nil case",
			model.ErrInvalidBlock,
			nil,
			rand.Int31(),
		},
		{
			"round -1 case",
			nil,
			RandomBlock(t),
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
			assert.Equal(t, GetHash(t, c.expectedBlock), GetHash(t, proposal.GetBlock()))
			assert.Equal(t, c.expectedRound, proposal.GetRound())
		})
	}
}

func TestVoteMessageFactory(t *testing.T) {
	for _, c := range []struct {
		name          string
		expectedError error
		expectedHash  []byte
	}{
		{
			"case 1",
			nil,
			RandomByte(),
		},
		{
			"case 2",
			nil,
			RandomByte(),
		},
		{
			"case 3",
			nil,
			RandomByte(),
		},
		{
			"case 4",
			nil,
			RandomByte(),
		},
		{
			"hash nil case no problem",
			nil,
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			vote := NewModelFactory().NewVoteMessage(c.expectedHash)
			assert.Equal(t, c.expectedHash, vote.GetBlockHash())
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
			RandomByte(),
			RandomByte(),
		},
		{
			"case 2",
			RandomByte(),
			RandomByte(),
		},
		{
			"case 3",
			RandomByte(),
			RandomByte(),
		},
		{
			"case 4",
			RandomByte(),
			RandomByte(),
		},
		{
			"case 5",
			RandomByte(),
			RandomByte(),
		},
		{
			"pub nil case no problem",
			nil,
			RandomByte(),
		},
		{
			"sig nil case no problem",
			RandomByte(),
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
			RandomStr(),
			RandomInvalidSig(),
			validPub,
			validPriv,
		},
		{
			"case 2",
			nil,
			RandomStr(),
			RandomInvalidSig(),
			validPub,
			validPriv,
		},
		{
			"case 3",
			nil,
			RandomStr(),
			RandomInvalidSig(),
			validPub,
			validPriv,
		},
		{
			"empty string case is valid",
			nil,
			"",
			RandomInvalidSig(),
			validPub,
			validPriv,
		},
		{
			"signature nil case",
			model.ErrInvalidSignature,
			RandomStr(),
			nil,
			validPub,
			validPriv,
		},
		{
			"pubkey nil case",
			ErrCryptoVerify,
			RandomStr(),
			RandomInvalidSig(),
			nil,
			validPriv,
		},
		{
			"privkey nil case",
			ErrCryptoSign,
			RandomStr(),
			RandomInvalidSig(),
			validPub,
			nil,
		},
		{
			"all ng case",
			multierr.Combine(model.ErrInvalidSignature, ErrCryptoSign, ErrCryptoSign),
			RandomStr(),
			nil,
			nil,
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			tx, err := NewTxModelBuilder().Message(c.expectedMessage).
				Signature(c.expectedSignature).
				Sign(c.expectedPubkey, c.expectedPrivKey).
				Build()
			if c.expectedError != nil {
				MultiErrorCheck(t, errors.Cause(err), c.expectedError)
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, c.expectedMessage, tx.GetPayload().GetMessage())
			assert.Equal(t, c.expectedSignature.GetPubkey(), tx.GetSignatures()[0].GetPubkey())
			assert.Equal(t, c.expectedSignature.GetSignature(), tx.GetSignatures()[0].GetSignature())
			assert.Equal(t, c.expectedPubkey, tx.GetSignatures()[1].GetPubkey())
			signature, err := Sign(c.expectedPrivKey, GetHash(t, tx))
			require.NoError(t, err)
			assert.Equal(t, signature, tx.GetSignatures()[1].GetSignature())
		})
	}
}

func TestNewPeer(t *testing.T) {
	for _, c := range []struct {
		name    string
		address string
		pubkey  []byte
	}{
		{
			"case 1",
			"111.111.111.111",
			RandomByte(),
		},
		{
			"case 2",
			RandomStr(),
			RandomByte(),
		},
		{
			"case 3",
			"localhost",
			nil,
		},
		{
			"case 4",
			"",
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			peer := NewModelFactory().NewPeer(c.address, c.pubkey)
			assert.Equal(t, c.address, peer.GetAddress())
			assert.Equal(t, c.pubkey, peer.GetPubkey())
		})
	}
}
