package config

type BBFTConfig struct {
	Host      string `default:"localhost"`
	Port      string `default:"50053"`
	SecretKey string `default:"secret_key"`
}

var config BBFTConfig

func GetConfig() *BBFTConfig {
	return &config
}

func GetTestConfig() *BBFTConfig {
	testConfig := config
	return &testConfig
}
