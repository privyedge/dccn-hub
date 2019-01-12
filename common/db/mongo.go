package dbcommon

import (
	"os"
	"strconv"
	"sync"
	"time"

	mgo "gopkg.in/mgo.v2"
)

var once sync.Once

var (
	DEFAULT_DB           = "dccn"
	DEFAULTT_COLLECTIOIN = "user"
	DEFAULT_HOST         = "localhost:27017"
	DEFAULT_POOL_LIMIT   = 4096
	DEFAULT_TIMEOUT      = 30
)

// Config uses to init a db connect
type Config struct {
	// DB db name
	DB string `json:"db"`
	// Collection db table
	Collection string `json:"collection"`
	// Host holds the addresses for the seed servers.
	Host string `json:"host"`
	// PoolLimit defines the per-server socket pool limit. Defaults to 4096.
	PoolLimit int `json:"pool_limit"`
	// Timeout is the amount of time to wait for a server to respond
	Timeout int `json:"timeout"`
}

// CreateDBConnection returns a db connection, it is recommended to use for once, and then use copy or clone to reuse it
// remembers to close after every copy() or clone()
func CreateDBConnection(conf Config) (s *mgo.Session, err error) {
	once.Do(func() {
		info := mgo.DialInfo{
			Addrs:     []string{conf.Host},
			Timeout:   time.Duration(conf.Timeout) * time.Second,
			PoolLimit: conf.PoolLimit,
		}

		if s, err = mgo.DialWithInfo(&info); err != nil {
			return
		}
		s.SetMode(mgo.Monotonic, true)
		s.SetSafe(&mgo.Safe{})
	})
	return
}

// LoadFromEnv Load DB Config from env.
func LoadFromEnv() (Config, error) {
	var conf Config

	if conf.Host = os.Getenv("DB_HOST"); conf.Host == "" {
		conf.Host = DEFAULT_HOST
	}

	if conf.DB = os.Getenv("DB"); conf.DB == "" {
		conf.DB = DEFAULT_DB
	}

	if conf.Collection = os.Getenv("DB"); conf.Collection == "" {
		conf.Collection = DEFAULTT_COLLECTIOIN
	}

	var err error

	if poolLimit := os.Getenv("DB_POOL_LIMIT"); poolLimit == "" {
		conf.PoolLimit = DEFAULT_POOL_LIMIT
	} else {
		conf.PoolLimit, err = strconv.Atoi(poolLimit)
		if err != nil {
			return Config{}, err
		}
	}

	if timeout := os.Getenv("DB_TIMEOUT"); timeout == "" {
		conf.Timeout = DEFAULT_TIMEOUT
	} else {
		conf.Timeout, err = strconv.Atoi(timeout)
	}
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}
