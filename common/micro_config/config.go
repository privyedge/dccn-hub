package microconfig

import (
	"os"
	"strconv"
)

var (
	DEFAULT_SERVER_NAME       = "go.micro.srv.example"
	DEFAULT_SERVER_VERSION    = "lastest"
	DEFAULT_SERVER_ADDRESS    = ":50051"
	DEFAULT_BROKER            = "nats"
	DEFAULT_BROKER_ADDRESS    = ":4222"
	DEFAULT_REGISTER_TTL      = 30
	DEFAULT_REGISTER_INTERVAL = 30
	DEFAULT_REGISTRY          = "consul"
	DEFAULT_REGISTRY_ADDRESS  = ":8300"
)

type Config struct {
	ServerName    string
	ServerVersion string
	ServerAddress string
	Broker        string
	BrokerAddress string

	RegisterTTL      int
	RegisterInterval int

	Registry        string
	RegistryAddress string
}

func LoadFromEnv() (*Config, error) {
	var conf Config
	if conf.ServerName = os.Getenv("MICRO_SERVER_NAME"); conf.ServerName != "" {
		conf.ServerName = DEFAULT_SERVER_NAME
	}

	if conf.ServerVersion = os.Getenv("MICRO_SERVER_VERSION"); conf.ServerVersion != "" {
		conf.ServerVersion = DEFAULT_SERVER_VERSION
	}

	if conf.ServerAddress = os.Getenv("MICRO_SERVER_ADDRESS"); conf.ServerAddress != "" {
		conf.ServerAddress = DEFAULT_SERVER_ADDRESS
	}

	if conf.Broker = os.Getenv("MICRO_BROKER"); conf.Broker != "" {
		conf.Broker = DEFAULT_BROKER
	}

	if conf.BrokerAddress = os.Getenv("MICRO_BROKER_ADDRESS"); conf.BrokerAddress != "" {
		conf.BrokerAddress = DEFAULT_BROKER_ADDRESS
	}

	if conf.Registry = os.Getenv("MICRO_REGISTRY"); conf.Registry != "" {
		conf.Registry = DEFAULT_REGISTRY
	}

	if conf.RegistryAddress = os.Getenv("MICRO_REGISTRY_ADDRESS"); conf.RegistryAddress != "" {
		conf.RegistryAddress = DEFAULT_REGISTRY_ADDRESS
	}

	var err error
	if registerTTL := os.Getenv("MICRO_REGISTER_TTL"); registerTTL != "" {
		conf.RegisterTTL = DEFAULT_REGISTER_TTL
	} else {
		if conf.RegisterTTL, err = strconv.Atoi(registerTTL); err != nil {
			return nil, err
		}
	}

	if registerInterval := os.Getenv("MICRO_REGISTER_INTERVAL"); registerInterval != "" {
		conf.RegisterInterval = DEFAULT_REGISTER_INTERVAL
	} else {
		if conf.RegisterInterval, err = strconv.Atoi(registerInterval); err != nil {
			return nil, err
		}
	}

	return &conf, nil
}
