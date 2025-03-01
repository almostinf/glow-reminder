package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/almostinf/glow-reminder/internal/domain"
	"github.com/almostinf/glow-reminder/internal/repository/pg"
	"github.com/almostinf/glow-reminder/internal/repository/redis"
	"github.com/almostinf/glow-reminder/pkg/clock"
	"github.com/almostinf/glow-reminder/pkg/glow_reminder/client/operations"
	"github.com/almostinf/glow-reminder/pkg/glow_reminder/models"
	"github.com/almostinf/glow-reminder/pkg/logger"
	"gopkg.in/tomb.v2"
)

type ReminderScheduler interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type reminderScheduler struct {
	cfg                Config
	reminderTaskRepo   redis.ReminderTaskRepo
	reminderRepo       pg.ReminderRepo
	logger             logger.Logger
	tomb               tomb.Tomb
	clock              clock.Clock
	glowReminderClient operations.ClientService
}

func New(
	cfg Config,
	reminderTaskRepo redis.ReminderTaskRepo,
	reminderRepo pg.ReminderRepo,
	logger logger.Logger,
	clock clock.Clock,
	glowReminderClient operations.ClientService,
) *reminderScheduler {
	return &reminderScheduler{
		cfg:                cfg,
		reminderTaskRepo:   reminderTaskRepo,
		reminderRepo:       reminderRepo,
		logger:             logger,
		clock:              clock,
		tomb:               tomb.Tomb{},
		glowReminderClient: glowReminderClient,
	}
}

func (scheduler *reminderScheduler) Start(ctx context.Context) error {
	scheduler.logger.Debug("Start reminder scheduler", map[string]interface{}{})

	scheduler.tomb.Go(func() error {
		ctx = context.WithoutCancel(ctx)

		ticker := time.NewTicker(scheduler.cfg.CycleDuration)
		defer ticker.Stop()

		for {
			select {
			case <-scheduler.tomb.Dying():
				if err := scheduler.tomb.Err(); err != nil {
					scheduler.logger.Error("failed to tomb dying", map[string]interface{}{
						"error": err.Error(),
					})
				}
			case <-ticker.C:
				scheduler.tomb.Go(func() error {
					if err := scheduler.processReminderTasks(ctx); err != nil {
						scheduler.logger.Error("failed to process reminder tasks", map[string]interface{}{
							"error": err.Error(),
						})
					}
					return nil
				})
			}
		}
	})

	return nil
}

func (scheduler *reminderScheduler) processReminderTasks(ctx context.Context) error {
	scheduler.logger.Debug("Start process reminder tasks", map[string]interface{}{})

	now := scheduler.clock.NowUnix()

	reminderTasks, err := scheduler.reminderTaskRepo.GetReminderTasks(ctx, now)
	if err != nil {
		return fmt.Errorf("failed to GetReminderTasks: %w", err)
	}

	scheduler.logger.Info("Process reminder tasks", map[string]interface{}{
		"reminder_tasks": reminderTasks,
	})

	reminders := make([]*domain.Reminder, 0, len(reminderTasks))
	for _, reminderTask := range reminderTasks {
		reminder, err := scheduler.reminderRepo.GetReminder(ctx, reminderTask.ID)
		if err != nil {
			return fmt.Errorf("failed to GetReminder %s: %w", reminderTask.ID, err)
		}
		reminders = append(reminders, reminder)
	}

	for _, reminder := range reminders {
		_, err := scheduler.glowReminderClient.GlowReminder(&operations.GlowReminderParams{
			Body: &models.GlowReminder{
				Colour: int64(reminder.Colour),
				Mode:   int64(reminder.Mode),
			},
		})
		if err != nil {
			return fmt.Errorf("failed to GlowReminder %s: %w", reminder.ID, err)
		}

		if err = scheduler.reminderRepo.DeleteReminder(ctx, reminder.ID); err != nil {
			return fmt.Errorf("failed to DeleteReminder %v: %w", reminder.ID, err)
		}
	}

	return nil
}

func (scheduler *reminderScheduler) Stop(_ context.Context) error {
	scheduler.logger.Debug("Stop reminder scheduler", map[string]interface{}{})

	return scheduler.tomb.Wait()
}
