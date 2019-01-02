package config

import (
	"testing"
)

var (
	conf *Config
	err error
	path = "config.json"
)

func TestNew(t *testing.T) {
	conf, err = New(path)
	if err != nil {
		t.Error(err.Error())
	}

	t.Logf("Load config %+v", conf)
}

// TODO:
// func TestWatch(t *testing.T) {
// 	TestNew(t)
//
// 	conf.Watch(path, func() {
// 		fmt.Printf("Refresh Config %+v", conf)
// 	})
//
// 	time.Sleep(18 * time.Second)
// 	conf.done <- struct{}{}
// }
