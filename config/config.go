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
	QueueLimits                           int `default:"5000"`
	LockedRegisteredLimits                int `default:"1000"`
	LockedVotedLimits                     int `default:"3000"`
	ReceivePropagateTxPoolLimits          int `default:"5000"`
	ReceiveProposeProposalPoolLimits      int `default:"5000"`
	ReceiveVoteVoteMessagePoolLimits      int `default:"5000"`
	ReceivePreCommitVoteMessagePoolLimits int `default:"500"`
	PreCommitFinderLimits                 int `default:"500"`

	// Consensus Parameter
	NumberOfBlockHasTransactions int           `default:"200"`
	AllowedConnectDelayTime      time.Duration `default:"500ms"`
	ProposeMaxCalcTime           time.Duration `default:"500ms"`
	VoteMaxCalcTime              time.Duration `default:"1s"`
	PreCommitMaxCalcTime         time.Duration `default:"200ms"`
	CommitMaxCalcTime            time.Duration `default:"500ms"`

	Demo Demo
}

var config BBFTConfig

func Init() {
	envconfig.MustProcess("bbft", &config)
}

func GetConfig() *BBFTConfig {
	return &config
}
