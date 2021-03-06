package test_utils

import (
	"github.com/satellitex/bbft/config"
	"github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"testing"
)

type TxGateSender struct {
	client bbft.TxGateClient
	t      *testing.T
}

func (s *TxGateSender) Write(ctx context.Context, tx model.Transaction) error {
	ptx := tx.(*convertor.Transaction).Transaction
	res, err := s.client.Write(ctx, ptx)
	if err == nil {
		require.Equal(s.t, &bbft.TxResponse{}, res)
	}
	return err
}

func NewTxGateSender(t *testing.T, conf *config.BBFTConfig) *TxGateSender {
	conn, err := grpc.Dial(conf.Host+":"+conf.Port, grpc.WithInsecure())
	require.NoError(t, err)
	return &TxGateSender{
		bbft.NewTxGateClient(conn),
		t,
	}
}
