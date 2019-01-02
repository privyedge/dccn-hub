package config

import (
	"fmt"
	"sync"

	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"

	"github.com/Ankr-network/refactor/util"
	"github.com/Ankr-network/refactor/util/db"
	"github.com/Ankr-network/refactor/app_dccn_usermgr/token"
)


type Config struct {
	DBConfig   utildb.Config
	TokenConfig token.Config

	DB         string
	Collection string

	TTL int
	Interval int

	once       sync.Once
	update     chan struct{}
	done       chan struct{}
	watcher config.Watcher
}

func New(path string) (*Config, error) {

	err := config.Load(
		file.NewSource(
			file.WithPath(path),
		))

	if err != nil {
		return nil, err
	}

	// w, err := config.Watch("DB")
	// if err != nil {
	// 	util.WriteLog(err.Error())
	// 	return nil, err
	// }

	conf := &Config{
		update: make(chan struct{}, 1),
		done: make(chan struct{}, 1),
		// watcher: w,
	}

	if err = config.Scan(&conf); err != nil {
		return nil, err
	}

	return conf, nil
}

// Watch watches the config's update.
func (p *Config) Watch(path string, relaod func()) {
	for {
		select {
		case <- p.done:
			return
		default:
			v, err := p.watcher.Next()
			if err != nil {
				util.WriteLog(err.Error())
				// ignores the errors
				break
			}

			tmpConf := Config{}
			err = v.Scan(&tmpConf)
			if err != nil {
				util.WriteLog(err.Error())
				// ignores the errors
				break
			}

			fmt.Printf("refresh config: %+v", tmpConf)
			p.Refresh(tmpConf)

			// Sends refresh signal.
			// p.update <- struct{}{}
			relaod()
		}
	}
}

// refresh update the value with new
func (p *Config) Refresh(tmp Config)  {
	util.WriteLog("Relaod configaration")
	p.DBConfig = tmp.DBConfig
}

// Finalize ensures all resources released
func (p *Config) Finalize()  {
	p.once.Do(func() {
		close(p.update)
		close(p.done)
		p.watcher.Stop()
	})
}
