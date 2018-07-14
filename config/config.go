package config

import _ "github.com/kelseyhightower/envconfig"

type BBFTConfig struct {
	Host        string `default:"localhost"`
	Port        string `default:"50053"`
	SecretKey   string `default:"secret_key"`
	QueueLimits int
	LockLimits int
}

var config BBFTConfig

func GetConfig() *BBFTConfig {
	return &config
}

func GetTestConfig() *BBFTConfig {
	testConfig := config
	testConfig.QueueLimits = 100
	testConfig.LockLimits = 100
	return &testConfig
}
