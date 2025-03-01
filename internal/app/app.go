package app

import (
	"context"

	"github.com/almostinf/glow-reminder/config"
	"github.com/almostinf/glow-reminder/internal/bot"
	"github.com/almostinf/glow-reminder/internal/repository/pg"
	"github.com/almostinf/glow-reminder/internal/repository/redis"
	"github.com/almostinf/glow-reminder/internal/scheduler"
	"github.com/almostinf/glow-reminder/internal/usecase"
	"github.com/almostinf/glow-reminder/pkg/clock"
	"github.com/almostinf/glow-reminder/pkg/glow_reminder/client"
	"github.com/almostinf/glow-reminder/pkg/glow_reminder/client/operations"
	"github.com/almostinf/glow-reminder/pkg/logger"
	"github.com/almostinf/glow-reminder/pkg/postgres"
	rediswrapper "github.com/almostinf/glow-reminder/pkg/redis"
	"github.com/go-openapi/strfmt"
	"go.uber.org/fx"
)

func CreateApp() fx.Option {
	return fx.Options(
		fx.Provide(
			appCtx,
			config.New,
			logger.FromAppConfig,
			logger.NewLogrusLogger,
			bot.FromAppConfig,
			clock.New,
			bot.New,
			usecase.NewReminder,
			fx.Annotate(usecase.NewReminder, fx.As(new(usecase.ReminderUsecase))),
			pg.NewReminderRepo,
			fx.Annotate(pg.NewReminderRepo, fx.As(new(pg.ReminderRepo))),
			postgres.FromAppConfig,
			postgres.New,
			fx.Annotate(bot.New, fx.As(new(bot.Bot))),
			rediswrapper.FromAppConfig,
			rediswrapper.New,
			redis.NewReminderTaskRepo,
			fx.Annotate(redis.NewReminderTaskRepo, fx.As(new(redis.ReminderTaskRepo))),
			defaultStrfmtRegistry,
			glowReminderTransportCfg,
			client.NewHTTPClientWithConfig,
			glowReminderOperations,
			scheduler.FromAppConfig,
			scheduler.New,
			fx.Annotate(scheduler.New, fx.As(new(scheduler.ReminderScheduler))),
		),
		fx.Invoke(
			startBot,
			startScheduler,
		),
	)
}

func startBot(b bot.Bot, lc fx.Lifecycle) error {
	lc.Append(
		fx.Hook{
			OnStart: b.Start,
			OnStop:  b.Stop,
		},
	)

	return nil
}

func startScheduler(scheduler scheduler.ReminderScheduler, lc fx.Lifecycle) error {
	lc.Append(
		fx.Hook{
			OnStart: scheduler.Start,
			OnStop:  scheduler.Stop,
		},
	)

	return nil
}

func appCtx() context.Context {
	return context.Background()
}

func defaultStrfmtRegistry() strfmt.Registry {
	return strfmt.Default
}

func glowReminderTransportCfg(cfg *config.AppConfig) *client.TransportConfig {
	return &client.TransportConfig{
		Host:    cfg.GlowReminderClient.Host,
		Schemes: []string{"http"},
	}
}

func glowReminderOperations(client *client.GlowReminder) operations.ClientService {
	return client.Operations
}
