package convertor_test

import (
	"context"
	"fmt"
	. "github.com/satellitex/bbft/convertor"
	. "github.com/satellitex/bbft/test_utils"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestAuthor(t *testing.T) {
	conf := GetTestConfig()

	ps := RandomPeerService(t, 4)

	author := NewAuthoer(ps)

	t.Run("failed case, Not found conf peer in PeerService", func(t *testing.T) {
		proto := RandomProposal(t)
		ctx, err := NewContextByProtobuf(conf, proto.(*Proposal))
		assert.NoError(t, err)

		_, err = author.DefaultReceiveAuth(ctx)
		ValidateStatusCode(t, err, codes.Unauthenticated)
	})

	// add conf peer to peer service
	ps.AddPeer(NewModelFactory().NewPeer(conf.Host, conf.PublicKey))

	t.Run("success case", func(t *testing.T) {
		proto := RandomProposal(t)
		ctx, err := NewContextByProtobufDebug(conf, proto.(*Proposal))
		assert.NoError(t, err)
		fmt.Println(ctx)

		_, err = author.DefaultReceiveAuth(ctx)
		assert.NoError(t, err)

		_, err = author.ProtoAurhorize(ctx, proto.(*Proposal))

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
}
