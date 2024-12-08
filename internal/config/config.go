package config

import "github.com/RacoonMediaServer/rms-packages/pkg/configuration"

// Configuration represents entire service configuration
type Configuration struct {
	AllowedDomains map[string]Domain `json:"allowed-domains"`
	Debug          configuration.Debug
	Http           configuration.Http
	Base           string
}

type Domain struct {
	ContentType []string `json:"content-type"`
}

var config Configuration

// Load open and parses configuration file
func Load(configFilePath string) error {
	return configuration.Load(configFilePath, &config)
}

// Config returns loaded configuration
func Config() Configuration {
	return config
}
