package logger

import "github.com/almostinf/glow-reminder/config"

type Config struct {
	Level string
}

func FromAppConfig(appCfg *config.AppConfig) Config {
	return Config{
		Level: appCfg.Log.Level,
	}
}
