package config

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	conf, err := Get("../testdata/config.toml")
	if err != nil {
		t.Error(err.Error())
	}

	fmt.Printf("Load config %+v", conf)
}

