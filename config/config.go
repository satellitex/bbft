package config

import _ "github.com/kelseyhightower/envconfig"

type BBFTConfig struct {
	Host                   string `default:"localhost"`
	Port                   string `default:"50053"`
	SecretKey              string `default:"secret_key"`
	QueueLimits            int
	LockedRegisteredLimits int
	LockedVotedLimits      int
}

var config BBFTConfig

func GetConfig() *BBFTConfig {
	return &config
}

func GetTestConfig() *BBFTConfig {
	testConfig := config
	testConfig.QueueLimits = 100
	testConfig.LockedRegisteredLimits = 100
	testConfig.LockedVotedLimits = 500
	return &testConfig
}
