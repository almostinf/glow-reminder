package bot

import (
	"time"

	"github.com/almostinf/glow-reminder/config"
)

type Config struct {
	Token         string
	PollerTimeout time.Duration
}

func FromAppConfig(appCfg *config.AppConfig) Config {
	return Config{
		Token:         appCfg.Bot.Token,
		PollerTimeout: appCfg.Bot.PoolerTimeout,
	}
}
