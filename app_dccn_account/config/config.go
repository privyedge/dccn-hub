package config

import (
	"github.com/Ankr-network/refactor/util"
	"github.com/Ankr-network/refactor/util/db"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
	"sync"
)


type ReloadFunc func()

type Config struct {
	once sync.Once
	DB utildb.DBConfig
	DBName string
	Collection string
	update chan struct{}
	done chan struct{}
}

func New(path string) (*Config, error) {

	err := config.Load(
		file.NewSource(
			file.WithPath(path),
		))

	if err != nil {
		return nil, err
	}

	conf := &Config{
		update: make(chan struct{}),
		done: make(chan struct{}),
	}

	if err = config.Scan(&conf); err != nil {
		return nil, err
	}

	return conf, nil
}

// Watch watches the config's update.
func (p *Config) Watch(path string, relaod func()) {
	w, err := config.Watch(path)
	if err != nil {
		util.WriteLog(err.Error())
		return
	}
	for {
		select {
		case <- p.done:
			return
		default:
			v, err := w.Next()
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

			p.Refresh(tmpConf)

			// Sends refresh signal.
			p.update <- struct{}{}
			relaod()
		}
	}
}

// refresh update the value with new
func (p *Config) Refresh(tmp Config)  {
	p.DB = tmp.DB
}

func (p *Config) WatchNew() <- chan struct{} {
	return p.update
}

// Finalize ensures all resources released
func (p *Config) Finalize()  {
	p.once.Do(func() {
		close(p.update)
		close(p.done)
	})
}
