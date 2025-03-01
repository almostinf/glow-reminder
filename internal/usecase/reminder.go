package usecase

import (
	"context"
	"fmt"

	"github.com/almostinf/glow-reminder/internal/domain"
	"github.com/almostinf/glow-reminder/internal/repository/pg"
	"github.com/almostinf/glow-reminder/internal/repository/redis"
	"github.com/almostinf/glow-reminder/pkg/logger"
	"github.com/google/uuid"
)

type ReminderUsecase interface {
	GetReminders(ctx context.Context, params domain.GetRemindersParams) ([]*domain.Reminder, error)
	CreateReminder(ctx context.Context, reminder domain.Reminder) error
	DeleteReminder(ctx context.Context, id uuid.UUID) error
}

type reminderUsecase struct {
	reminderRepo     pg.ReminderRepo
	reminderTaskRepo redis.ReminderTaskRepo
	logger           logger.Logger
}

func NewReminder(reminderRepo pg.ReminderRepo, reminderTaskRepo redis.ReminderTaskRepo, logger logger.Logger) *reminderUsecase {
	return &reminderUsecase{
		reminderRepo:     reminderRepo,
		reminderTaskRepo: reminderTaskRepo,
		logger:           logger,
	}
}

func (usecase *reminderUsecase) GetReminders(ctx context.Context, params domain.GetRemindersParams) ([]*domain.Reminder, error) {
	return usecase.reminderRepo.GetReminders(ctx, params)
}

func (usecase *reminderUsecase) CreateReminder(ctx context.Context, reminder domain.Reminder) error {
	if err := usecase.reminderTaskRepo.AddReminderTask(ctx, &domain.ReminderTask{
		ID:          reminder.ID,
		ScheduledAt: reminder.ScheduledAt,
	}); err != nil {
		return fmt.Errorf("failed to AddReminderTask: %w", err)
	}

	return usecase.reminderRepo.CreateReminder(ctx, reminder)
}

func (usecase *reminderUsecase) DeleteReminder(ctx context.Context, id uuid.UUID) error {
	return usecase.reminderRepo.DeleteReminder(ctx, id)
}
