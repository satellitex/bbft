package config

import (
	_ "github.com/kelseyhightower/envconfig"
)

type BBFTConfig struct {
	Host                                  string `default:"localhost"`
	Port                                  string `default:"50053"`
	SecretKey                             string `default:"secret_key"`
	QueueLimits                           int
	LockedRegisteredLimits                int
	LockedVotedLimits                     int
	ReceivePropagateTxPoolLimits          int
	ReceiveProposeProposalPoolLimits      int
	ReceiveVoteVoteMessagePoolLimits      int
	ReceivePreCommitVoteMessagePoolLimits int
}

var config BBFTConfig

func GetConfig() *BBFTConfig {
	return &config
}
