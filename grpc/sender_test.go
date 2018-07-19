package grpc_test

import (
	"context"
	"github.com/satellitex/bbft/config"
	"github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
	. "github.com/satellitex/bbft/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"sync"
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

func TestTxGateWrite(t *testing.T) {
	conf := GetTestConfig()

	sender := NewTxGateSender(t, conf)

	tx := RandomValidTx(t)
	err := sender.Write(context.TODO(), tx)
	assert.NoError(t, err)

	tx = RandomInvalidTx(t)
	err = sender.Write(context.TODO(), tx)
	ValidateStatusCode(t, err, codes.InvalidArgument)

	waiter := &sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		waiter.Add(1)
		go func() {
			err := sender.Write(context.TODO(), RandomValidTx(t))
			assert.NoError(t, err)
			waiter.Done()
		}()
	}
	waiter.Wait()
}
