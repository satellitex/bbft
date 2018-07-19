package convertor_test

import (
	"bytes"
	"context"
	. "github.com/satellitex/bbft/convertor"
	. "github.com/satellitex/bbft/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestAuthor(t *testing.T) {
	conf := GetTestConfig()

	ps := RandomPeerService(t, 4)

	author := NewAuthor(ps)

	t.Run("failed case, Not found conf peer in PeerService", func(t *testing.T) {
		proto := RandomProposal(t)
		ctx, err := NewContextByProtobufDebug(conf, proto.(*Proposal))
		assert.NoError(t, err)

		_, err = author.DefaultReceiveAuth(ctx)
		ValidateStatusCode(t, err, codes.PermissionDenied)

		_, err = author.ProtoAurhorize(ctx, proto.(*Proposal))
		ValidateStatusCode(t, err, codes.PermissionDenied)
	})

	// add conf peer to peer service
	ps.AddPeer(NewModelFactory().NewPeer(conf.Host, conf.PublicKey))

	t.Run("success case", func(t *testing.T) {
		proto := RandomProposal(t)
		ctx, err := NewContextByProtobufDebug(conf, proto.(*Proposal))
		assert.NoError(t, err)

		_, err = author.DefaultReceiveAuth(ctx)
		assert.NoError(t, err)

		_, err = author.ProtoAurhorize(ctx, proto.(*Proposal))
		assert.NoError(t, err)
	})

	t.Run("failed case, TODO", func(t *testing.T) {
		ctx := context.TODO()
		_, err := author.DefaultReceiveAuth(ctx)
		ValidateStatusCode(t, err, codes.Unauthenticated)
	})

	t.Run("failed case, unverified metadata", func(t *testing.T) {
		proto := RandomProposal(t)
		md := metadata.Pairs(HeaderAuthorizeSignature, NewAuthorSignatureStr([]byte("dummy")),
			HeaderAuthorizePubkey, NewAuthorPubKeyStr(conf.PublicKey))
		ctx := metadata.NewIncomingContext(context.Background(), md)

		_, err := author.DefaultReceiveAuth(ctx)
		assert.NoError(t, err)

		_, err = author.ProtoAurhorize(ctx, proto.(*Proposal))
		ValidateStatusCode(t, err, codes.Unauthenticated)
	})

	t.Run("test verify only leader", func(t *testing.T) {
		id := func() int32 {
			for i, p := range ps.GetPermutationPeers(0) {
				if bytes.Equal(p.GetPubkey(), conf.PublicKey) {
					return int32(i)
				}
			}
			return -1
		}()
		require.NotEqual(t, -1, id)
		proto := RandomProposalWithHeightRound(t, 0, id)
		ctx, err := NewContextByProtobufDebug(conf, proto.(*Proposal))
		require.NoError(t, err)

		// success
		_, err = author.VerifyOnlyLeader(ctx, proto)
		assert.NoError(t, err)

		evilConf := *conf
		epub, epri := NewKeyPair()
		evilConf.PublicKey = epub
		evilConf.SecretKey = epri
		ctx, err = NewContextByProtobufDebug(&evilConf, proto.(*Proposal))
		// failed not leader
		_, err = author.VerifyOnlyLeader(ctx, proto)
		ValidateStatusCode(t, err, codes.PermissionDenied)
	})
}
