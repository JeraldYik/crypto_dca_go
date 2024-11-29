package config

var config *Config

func MustInit() {
	Set(initConfig())
	timeInit(Get(), nil)
	addTimeRelatedConfigs(config)
}

func Get() *Config {
	return config
}

func Set(c *Config) {
	config = c
}
