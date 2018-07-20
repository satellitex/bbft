package test_utils

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/satellitex/bbft/config"
	"github.com/satellitex/bbft/convertor"
	"os"
)

func GetTestConfig() *config.BBFTConfig {
	testConfig := &config.BBFTConfig{}

	envconfig.MustProcess("bbft", testConfig)

	validPub, validPriv := convertor.NewKeyPair()
	testConfig.PublicKey = validPub
	testConfig.SecretKey = validPriv

	if os.Getenv("CIRCLECI") != "" {
		testConfig.QueueLimits = 20
		testConfig.LockedRegisteredLimits = 20
		testConfig.LockedVotedLimits = 30
		testConfig.ReceivePropagateTxPoolLimits = 20
		testConfig.ReceiveProposeProposalPoolLimits = 20
		testConfig.ReceiveVoteVoteMessagePoolLimits = 20
		testConfig.ReceivePreCommitVoteMessagePoolLimits = 20
		testConfig.PreCommitFinderLimits = 20
		testConfig.NumberOfBlockHasTransactions = 20
	} else {
		testConfig.QueueLimits = 100
		testConfig.LockedRegisteredLimits = 100
		testConfig.LockedVotedLimits = 500
		testConfig.ReceivePropagateTxPoolLimits = 100
		testConfig.ReceiveProposeProposalPoolLimits = 100
		testConfig.ReceiveVoteVoteMessagePoolLimits = 100
		testConfig.ReceivePreCommitVoteMessagePoolLimits = 100
		testConfig.PreCommitFinderLimits = 100
		testConfig.NumberOfBlockHasTransactions = 100
	}
	return testConfig
}
