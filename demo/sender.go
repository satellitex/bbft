package demo

import (
	. "github.com/satellitex/bbft/test_utils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"sync"
	"testing"
)

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
