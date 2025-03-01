package scheduler

import (
	"time"

	"github.com/almostinf/glow-reminder/config"
)

type Config struct {
	CycleDuration time.Duration
}

func FromAppConfig(appCfg *config.AppConfig) Config {
	return Config{
		CycleDuration: appCfg.Scheduler.CycleDuration,
	}
}
