package postgres

import (
	"time"

	"github.com/almostinf/glow-reminder/config"
)

type Config struct {
	URL          string
	MaxPoolSize  int
	ConnAttempts int
	ConnTimeout  time.Duration
}

func FromAppConfig(appCfg *config.AppConfig) Config {
	return Config{
		URL:          appCfg.PG.URL,
		MaxPoolSize:  appCfg.PG.PoolMax,
		ConnAttempts: appCfg.PG.ConnAttempts,
		ConnTimeout:  appCfg.PG.ConnTimeout,
	}
}
