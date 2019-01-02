package utildb

import (
	"gopkg.in/mgo.v2"
	"sync"
	"time"
)

var once sync.Once

// Config uses to init a db connect
type Config struct {
	// Addrs holds the addresses for the seed servers.
	Addrs []string `json:"addrs"`
	// PoolLimit defines the per-server socket pool limit. Defaults to 4096.
	PoolLimit int `json:"pool_limit"`
	// Timeout is the amount of time to wait for a server to respond
	Timeout int `json:"timeout"`
}

// CreateDBConnection returns a db connection, it is recommended to use for once, and then use copy or clone to reuse it
// remembers to close after every copy() or clone()
func CreateDBConnection(conf Config) (s *mgo.Session, err error) {
	once.Do( func() {
		info := mgo.DialInfo{
			Addrs: conf.Addrs,
			Timeout: time.Duration(conf.Timeout) * time.Second,
			PoolLimit: conf.PoolLimit,
		}

		if s, err = mgo.DialWithInfo(&info); err != nil {
			return
		}
		s.SetMode(mgo.Monotonic, true)
		s.SetSafe(&mgo.Safe{})
	} )
	return
}
