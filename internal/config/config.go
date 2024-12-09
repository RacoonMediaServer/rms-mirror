package config

import (
	"strings"

	"github.com/RacoonMediaServer/rms-packages/pkg/configuration"
)

// Configuration represents entire service configuration
type Configuration struct {
	AllowedDomains map[string]Domain `json:"allowed-domains"`
	Debug          configuration.Debug
	Http           configuration.Http
	Base           string
}

type Domain struct {
	ContentType []string `json:"content-type"`
	LimitMB     uint32   `json:"limit"`
}

func (d Domain) LimitBytes() int64 {
	return int64(d.LimitMB * 1024 * 1024)
}

func (d Domain) MakeAcceptHeader() string {
	if len(d.ContentType) == 0 {
		return "*/*"
	}

	return strings.Join(d.ContentType, ",")
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
