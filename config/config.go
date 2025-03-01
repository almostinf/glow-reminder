package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	PG struct {
		Dialect      string        `env-required:"true" yaml:"dialect" env:"DIALECT"`
		URL          string        `env-required:"true" yaml:"pg_url" env:"PG_URL"`
		PoolMax      int           `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		ConnAttempts int           `env-required:"true" yaml:"conn_attempts" env:"PG_CONN_ATTEMPTS"`
		ConnTimeout  time.Duration `env-required:"true" yaml:"conn_timeout" env:"PG_CONN_TIMEOUT"`
	}

	Bot struct {
		Token         string        `yaml:"token" env:"TOKEN"`
		PoolerTimeout time.Duration `env-required:"true" yaml:"pooler_timeout" env:"POOLER_TIMEOUT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level" env:"LOG_LEVEL"`
	}

	HTTP struct {
		Port uint32 `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		Host string `env-required:"true" yaml:"host" env:"HOST"`
	}

	Redis struct {
		Url string `env-required:"true" yaml:"redis_url" env:"REDIS_URL"`
	}

	Scheduler struct {
		CycleDuration time.Duration `env-required:"true" yaml:"cycle_duration" env:"CYCLE_DURATION"`
	}

	GlowReminderClient struct {
		Host string `env-required:"true" yaml:"host" env:"GLOW_REMINDER_CLIENT_HOST"`
	}

	AppConfig struct {
		App                App                `yaml:"app"`
		Bot                Bot                `yaml:"bot"`
		Redis              Redis              `yaml:"redis"`
		PG                 PG                 `yaml:"postgres"`
		HTTP               HTTP               `yaml:"http"`
		Log                Log                `yaml:"logger"`
		Scheduler          Scheduler          `yaml:"scheduler"`
		GlowReminderClient GlowReminderClient `yaml:"glow_reminder_client"`
	}
)

func getServiceFromConfig(path string, service interface{}) error {
	err := cleanenv.ReadConfig(path, service)
	if err != nil {
		return fmt.Errorf("read config error: %w", err)
	}

	err = cleanenv.ReadEnv(service)
	if err != nil {
		return fmt.Errorf("read env error: %w", err)
	}

	return nil
}

// Creates a new config entity after reading the configuration values
// from the YAML file and environment variables.
func New() (*AppConfig, error) {
	appConfig := &AppConfig{}
	if err := getServiceFromConfig("./config/config.yaml", appConfig); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return appConfig, nil
}
