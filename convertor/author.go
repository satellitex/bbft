package convertor

import (
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/satellitex/bbft/config"
	"github.com/satellitex/bbft/dba"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
)

var (
	HeaderAuthorizeSignature = "authorization_sig-bin"
	HeaderAuthorizePubkey    = "authorization_pub-bin"
)

func NewAuthorSignatureStr(signature []byte) string {
	return string(signature)
}

func NewAuthorPubKeyStr(pubkey []byte) string {
	return string(pubkey)
}

func NewContextByProtobuf(conf *config.BBFTConfig, proto proto.Message) (context.Context, error) {
	hash, err := CalcHashFromProto(proto)
	if err != nil {
		return nil, err
	}
	signature, err := Sign(conf.SecretKey, hash)
	if err != nil {
		return nil, err
	}
	md := metadata.Pairs(HeaderAuthorizeSignature, NewAuthorSignatureStr(signature),
		HeaderAuthorizePubkey, NewAuthorPubKeyStr(conf.PublicKey))
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	return ctx, nil
}

func NewContextByProtobufDebug(conf *config.BBFTConfig, proto proto.Message) (context.Context, error) {
	hash, err := CalcHashFromProto(proto)
	if err != nil {
		return nil, err
	}
	signature, err := Sign(conf.SecretKey, hash)
	if err != nil {
		return nil, err
	}
	md := metadata.Pairs(HeaderAuthorizeSignature, NewAuthorSignatureStr(signature),
		HeaderAuthorizePubkey, NewAuthorPubKeyStr(conf.PublicKey))
	ctx := metadata.NewIncomingContext(context.Background(), md)
	return ctx, nil
}

type Author struct {
	ps dba.PeerService
}

func NewAuthor(ps dba.PeerService) *Author {
	return &Author{ps}
}

func AuthParamFromMD(ctx context.Context, header string) (string, error) {
	val := metautils.ExtractIncoming(ctx).Get(header)
	if val == "" {
		return "", status.Errorf(codes.Unauthenticated, "Request unauthenticated header with "+header)

	}
	return val, nil
}

func (a *Author) GetPubkey(ctx context.Context) ([]byte, error) {
	pubStr, err := AuthParamFromMD(ctx, HeaderAuthorizePubkey)
	if err != nil {
		log.Println("Failed Auth FromMD pubkey: ", err)
		return nil, err
	}
	return []byte(pubStr), nil
}

func (a *Author) GetSignature(ctx context.Context) ([]byte, error) {
	sigStr, err := AuthParamFromMD(ctx, HeaderAuthorizeSignature)
	if err != nil {
		log.Println("Failed Auth FromMD signature: ", err)
		return nil, err
	}
	return []byte(sigStr), nil
}

func (a *Author) DefaultReceiveAuth(ctx context.Context) (context.Context, error) {
	pubkey, err := a.GetPubkey(ctx)
	if err != nil {
		return ctx, err
	}
	if _, ok := a.ps.GetPeer(pubkey); !ok {
		return ctx, status.Errorf(codes.PermissionDenied, "Failed Auth Unknown Peer's pubkey: %x", pubkey)
	}
	return ctx, nil
}

func (a *Author) ProtoAurhorize(ctx context.Context, proto proto.Message) (context.Context, error) {
	signature, err := a.GetSignature(ctx)
	if err != nil {
		return ctx, err
	}
	pubkey, err := a.GetPubkey(ctx)
	if err != nil {
		return ctx, err
	}
	hash, err := CalcHashFromProto(proto)
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, err.Error())
	}
	if err := Verify(pubkey, hash, signature); err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, err.Error())
	}
	if _, ok := a.ps.GetPeer(pubkey); !ok {
		return ctx, status.Errorf(codes.PermissionDenied, "Failed Auth Unknown Peer's pubkey: %x", pubkey)
	}
	return ctx, nil
}
