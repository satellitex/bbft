package test_utils

import "github.com/satellitex/bbft/config"

func GetTestConfig() *config.BBFTConfig {
	testConfig := config.GetConfig()
	testConfig.QueueLimits = 100
	testConfig.LockedRegisteredLimits = 100
	testConfig.LockedVotedLimits = 500
	testConfig.ReceivePropagateTxPoolLimits = 100
	testConfig.ReceiveProposeProposalPoolLimits = 100
	testConfig.ReceiveVoteVoteMessagePoolLimits = 100
	testConfig.ReceivePreCommitVoteMessagePoolLimits = 100
	return testConfig
}
