package config

import (
	"github.com/Ankr-network/dccn-hub/app_dccn_usermgr/token"
	dbcommon "github.com/Ankr-network/dccn-hub/common/db"
)

type Config struct {
	DB    dbcommon.Config
	Token token.Config
}
