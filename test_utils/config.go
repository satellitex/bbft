package test_utils

import (
	"github.com/satellitex/bbft/config"
	"os"
)

func GetTestConfig() *config.BBFTConfig {
	testConfig := config.GetConfig()

	if os.Getenv("CIRCLECI") != "" {
		testConfig.QueueLimits = 20
		testConfig.LockedRegisteredLimits = 20
		testConfig.LockedVotedLimits = 30
		testConfig.ReceivePropagateTxPoolLimits = 20
		testConfig.ReceiveProposeProposalPoolLimits = 20
		testConfig.ReceiveVoteVoteMessagePoolLimits = 20
		testConfig.ReceivePreCommitVoteMessagePoolLimits = 20
		testConfig.PreCommitFinderLimits = 20
	} else {
		testConfig.QueueLimits = 100
		testConfig.LockedRegisteredLimits = 100
		testConfig.LockedVotedLimits = 500
		testConfig.ReceivePropagateTxPoolLimits = 100
		testConfig.ReceiveProposeProposalPoolLimits = 100
		testConfig.ReceiveVoteVoteMessagePoolLimits = 100
		testConfig.ReceivePreCommitVoteMessagePoolLimits = 100
		testConfig.PreCommitFinderLimits = 100
	}
	return testConfig
}
