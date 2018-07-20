package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type BBFTConfig struct {
	Host                                  string `default:"localhost"`
	Port                                  string `default:"50053"`
	PublicKey                             []byte
	SecretKey                             []byte
	QueueLimits                           int `default:"1000"`
	LockedRegisteredLimits                int `default:"1000"`
	LockedVotedLimits                     int `default:"5000"`
	ReceivePropagateTxPoolLimits          int `default:"3000"`
	ReceiveProposeProposalPoolLimits      int `default:"100"`
	ReceiveVoteVoteMessagePoolLimits      int `default:"100"`
	ReceivePreCommitVoteMessagePoolLimits int `default:"100"`
	PreCommitFinderLimits                 int `default:"100"`

	// Consensus Parameter
	NumberOfBlockHasTransactions int           `default:"1000"`
	AllowedConnectDelayTime      time.Duration `default:"200ms"`
	ProposeMaxCalcTime           time.Duration `default:"300ms"`
	VoteMaxCalcTime              time.Duration `default:"300ms"`
	PreCommitMaxCalcTime         time.Duration `default:"100ms"`
	CommitMaxCalcTime            time.Duration `default:"200ms"`

	Demo Demo
}

var config BBFTConfig

func Init() {
	envconfig.MustProcess("bbft", &config)
}

func GetConfig() *BBFTConfig {
	return &config
}
