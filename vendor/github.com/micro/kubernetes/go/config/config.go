// Package config implements go-config with env and k8s configmap
package config

import (
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/configmap"
	"github.com/micro/go-config/source/env"
)

// NewConfig returns config with env and k8s configmap setup
func NewConfig(opts ...config.Option) config.Config {
	cfg := config.NewConfig()
	cfg.Load(
		env.NewSource(),
		configmap.NewSource(),
	)
	return cfg
}
