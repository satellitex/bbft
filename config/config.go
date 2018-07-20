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
	QueueLimits                           int `default:"100"`
	LockedRegisteredLimits                int `default:"100"`
	LockedVotedLimits                     int `default:"500"`
	ReceivePropagateTxPoolLimits          int `default:"1000"`
	ReceiveProposeProposalPoolLimits      int `default:"100"`
	ReceiveVoteVoteMessagePoolLimits      int `default:"100"`
	ReceivePreCommitVoteMessagePoolLimits int `default:"100"`
	PreCommitFinderLimits                 int `default:"100"`

	// Consensus Parameter
	NumberOfBlockHasTransactions int           `default:"100"`
	AllowedConnectDelayTime      time.Duration `default:"500ms"`
	ProposeMaxCalcTime           time.Duration `default:"1s"`
	VoteMaxCalcTime              time.Duration `default:"2s"`
	PreCommitMaxCalcTime         time.Duration `default:"500ms"`
	CommitMaxCalcTime            time.Duration `default:"1s"`

	Demo Demo
}

var config BBFTConfig

func Init() {
	envconfig.MustProcess("bbft", &config)
}

func GetConfig() *BBFTConfig {
	return &config
}
