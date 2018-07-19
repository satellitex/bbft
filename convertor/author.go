package convertor

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/satellitex/bbft/config"
	"github.com/satellitex/bbft/dba"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	headerAuthorize = "authorization"
	signatureLabel  = "signature"
	pubkeyLabel     = "pubkey"
)

func newSignatureStr(signature []byte) string {
	return signatureLabel + " " + string(signature)
}

func newPubKeyStr(pubkey []byte) string {
	return pubkeyLabel + " " + string(pubkey)
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
	md := metadata.Pairs(headerAuthorize, newSignatureStr(signature),
		headerAuthorize, newPubKeyStr(conf.PublicKey))
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	return ctx, nil
}

type Author struct {
	ps dba.PeerService
}

func GetPubkey(ctx context.Context) ([]byte, error) {
	pubStr, err := grpc_auth.AuthFromMD(ctx, pubkeyLabel)
	if err != nil {
		fmt.Printf("Failed Auth FromMD pubkey(%s) Label: %s", pubkeyLabel, err.Error())
		return nil, err
	}
	return []byte(pubStr), nil
}

func GetSignature(ctx context.Context) ([]byte, error) {
	sigStr, err := grpc_auth.AuthFromMD(ctx, signatureLabel)
	if err != nil {
		fmt.Printf("Failed Auth FromMD signature(%s) Label: %s", signatureLabel, err.Error())
		return nil, err
	}
	return []byte(sigStr), nil
}

func (a *Author) DefaultReceiveAuth(ctx context.Context) (context.Context, error) {
	pubkey, err := GetPubkey(ctx)
	if err != nil {
		return ctx, err
	}
	if _, ok := a.ps.GetPeer(pubkey); !ok {
		return ctx, status.Errorf(codes.Unauthenticated, "Failed Auth Unknown Peer's pubkey: %x", pubkey)
	}
	return ctx, nil
}

func (a *Author) ProtoAurhorize(ctx context.Context, proto proto.Message) (context.Context, error) {
	signature, err := GetSignature(ctx)
	if err != nil {
		return ctx, err
	}
	pubkey, err := GetPubkey(ctx)
	if err != nil {
		return ctx, err
	}
	if _, ok := a.ps.GetPeer(pubkey); !ok {
		return ctx, status.Errorf(codes.Unauthenticated, "Failed Auth Unknown Peer's pubkey: %x", pubkey)
	}
	hash, err := CalcHashFromProto(proto)
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, err.Error())
	}
	if err := Verify(pubkey, hash, signature); err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, err.Error())
	}
	return ctx, nil
}
