// Package config provides the configuration
package config

import (
	"time"

	"github.com/spf13/viper"
)

var (
	// Port is the port key.
	Port string
	// ReadTimeout is the read timeout key.
	ReadTimeout time.Duration
	// WriteTimeout is the write timeout key.
	WriteTimeout time.Duration
	// MaxHeaderBytes is the max header bytes key.
	MaxHeaderBytes int
)

// InitEnvironment is the environment initialization function.
func InitEnvironment() {
	viper.AutomaticEnv()

	viper.SetDefault("port", "8080")
	viper.SetDefault("read_timeout", 10)
	viper.SetDefault("write_timeout", 10)
	viper.SetDefault("max_header_bytes", 1<<20)

	Port = viper.GetString("port")
	ReadTimeout = time.Duration(viper.GetInt("read_timeout")) * time.Second
	WriteTimeout = time.Duration(viper.GetInt("write_timeout")) * time.Second
	MaxHeaderBytes = viper.GetInt("max_header_bytes")
}
