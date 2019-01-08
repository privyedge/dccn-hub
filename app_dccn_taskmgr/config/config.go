package config

import (
	dbcommon "github.com/Ankr-network/dccn-hub/common/db"

	config "github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
)

type Config struct {
	DBConfig dbcommon.Config `json:"db_config,omitempty"`

	AuthSrvName string `json:"auth_srv_name"`
	SrvName     string `json:"srv_name"`
	Version     string `json:"version"`
	TTL         int    `json:"ttl,omitempty"`
	Interval    int    `json:"interval,omitempty"`

	TopicPubNewTask      string `json:"topic_pub_new_task"`
	TopicPubCancelTask   string `json:"topic_pub_cancel_task"`
	TopicSubResultTask   string `json:"topic_sub_result_task"`
	TopicSubFeedbackTask string `json:"topic_sub_feedback_task"`
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
