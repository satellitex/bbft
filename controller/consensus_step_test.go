package controller_test

import (
	"context"
	"github.com/satellitex/bbft/config"
	. "github.com/satellitex/bbft/controller"
	"github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/dba"
	"github.com/satellitex/bbft/proto"
	. "github.com/satellitex/bbft/test_utils"
	"github.com/satellitex/bbft/usecase"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"testing"
)

func NewTestConsensusController(t *testing.T) (*config.BBFTConfig, *ConsensusController) {

	testConfig := GetTestConfig()
	queue := dba.NewProposalTxQueueOnMemory(testConfig)
	ps := RandomPeerService(t, 3)
	lock := dba.NewLockOnMemory(ps, testConfig)
	pool := dba.NewReceiverPoolOnMemory(testConfig)
	bc := dba.NewBlockChainOnMemory()
	slv := convertor.NewStatelessValidator()
	sender := convertor.NewMockConsensusSender()
	receivChan := usecase.NewReceiveChannel(testConfig)
	receiver := usecase.NewConsensusReceiverUsecase(queue, ps, lock, pool, bc, slv, sender, receivChan)

	author := convertor.NewAuthor(ps)

	// add peer this peer
	ps.AddPeer(convertor.NewModelFactory().NewPeer(testConfig.Host, testConfig.PublicKey))

	return testConfig, NewConsensusController(receiver, author)

}

func TestConsensusController_Propagate(t *testing.T) {

	conf, ctrl := NewTestConsensusController(t)

	validTx := RandomValidTx(t).(*convertor.Transaction).Transaction
	inValidTx := RandomInvalidTx(t).(*convertor.Transaction).Transaction

	evilConf := *conf
	pk, sk := convertor.NewKeyPair()
	evilConf.PublicKey = pk
	evilConf.SecretKey = sk

	for _, c := range []struct {
		ctx  context.Context
		tx   *bbft.Transaction
		code codes.Code
	}{
		{
			ctx:  ValidContext(t, conf, validTx),
			tx:   validTx,
			code: codes.OK,
		},
		{
			context.TODO(),
			validTx,
			codes.Unauthenticated,
		},
		{
			ValidContext(t, &evilConf, validTx),
			validTx,
			codes.PermissionDenied,
		},
		{
			ValidContext(t, conf, inValidTx),
			inValidTx,
			codes.InvalidArgument,
		},
		{
			ValidContext(t, conf, inValidTx),
			nil,
			codes.Unauthenticated,
		},
		{
			ValidContext(t, conf, validTx),
			validTx,
			codes.AlreadyExists,
		},
	} {
		_, err := ctrl.Propagate(c.ctx, c.tx)
		if c.code != codes.OK {
			ValidateStatusCode(t, err, c.code)
		} else {
			assert.NoError(t, err)
		}
	}

}
