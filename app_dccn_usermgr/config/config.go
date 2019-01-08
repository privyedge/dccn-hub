package config

import (
	"github.com/Ankr-network/dccn-hub/app_dccn_usermgr/token"
	dbcommon "github.com/Ankr-network/dccn-hub/common/db"

	config "github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
)

type Config struct {
	DBConfig    dbcommon.Config `json:"db_config,omitempty"`
	TokenConfig token.Config    `json:"token_config,omitempty"`

	SrvName  string `json:"srv_name"`
	Version  string `json:"version"`
	TTL      int    `json:"ttl,omitempty"`
	Interval int    `json:"interval,omitempty"`
}

func New(path string) (*Config, error) {

	err := config.Load(
		file.NewSource(
			file.WithPath(path),
		))

	if err != nil {
		return nil, err
	}

	var conf Config
	if err = config.Scan(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
