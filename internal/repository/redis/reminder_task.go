package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/almostinf/glow-reminder/internal/domain"
	"github.com/almostinf/glow-reminder/pkg/logger"
	rediswrapper "github.com/almostinf/glow-reminder/pkg/redis"
	"github.com/redis/go-redis/v9"
)

const reminderTasksKey = "reminder-tasks"

var _ ReminderTaskRepo = (*reminderTaskRepo)(nil)

type ReminderTaskRepo interface {
	AddReminderTask(ctx context.Context, reminderTask *domain.ReminderTask) error
	GetReminderTasks(ctx context.Context, to int64) ([]*domain.ReminderTask, error)
}

type reminderTaskRepo struct {
	redis  *rediswrapper.Redis
	logger logger.Logger
}

func NewReminderTaskRepo(redis *rediswrapper.Redis, logger logger.Logger) *reminderTaskRepo {
	return &reminderTaskRepo{
		redis:  redis,
		logger: logger,
	}
}

func (repo *reminderTaskRepo) AddReminderTask(ctx context.Context, reminderTask *domain.ReminderTask) error {
	reminderTaskBytes, err := json.Marshal(reminderTask)
	if err != nil {
		return fmt.Errorf("failed to marshal reminder task: %w", err)
	}

	z := redis.Z{
		Member: reminderTaskBytes,
		Score:  float64(reminderTask.ScheduledAt.Unix()),
	}

	if err = repo.redis.ZAdd(ctx, reminderTasksKey, z).Err(); err != nil {
		return fmt.Errorf("failed to ZAdd reminder task %v: %w", reminderTask.ID, err)
	}

	repo.logger.Info("Add new reminder task", map[string]interface{}{
		"id":           reminderTask.ID,
		"scheduled_at": reminderTask.ScheduledAt,
	})

	return nil
}

func (repo *reminderTaskRepo) GetReminderTasks(ctx context.Context, to int64) ([]*domain.ReminderTask, error) {
	pipe := repo.redis.TxPipeline()

	min := "-inf"
	max := strconv.FormatInt(to, 10)

	pipe.ZRangeByScore(ctx, reminderTasksKey, &redis.ZRangeBy{
		Min: min,
		Max: max,
	})

	pipe.ZRemRangeByScore(ctx, reminderTasksKey, min, max)

	resp, err := pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to Exec GetReminderTasks: %w", err)
	}

	reminderTasksBytes, err := resp[0].(*redis.StringSliceCmd).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to read reminder tasks: %w", err)
	}

	reminderTasks := make([]*domain.ReminderTask, 0, len(reminderTasksBytes))
	for _, reminderTaskBytes := range reminderTasksBytes {
		reminderTask := &domain.ReminderTask{}
		if err = json.Unmarshal([]byte(reminderTaskBytes), &reminderTask); err != nil {
			return nil, fmt.Errorf("failed to unmarshal reminder task bytes: %w", err)
		}
		reminderTasks = append(reminderTasks, reminderTask)
	}

	return reminderTasks, nil
}

func (repo *reminderTaskRepo) DeleteReminderTask(ctx context.Context, reminderTask *domain.ReminderTask) error {
	pipe := repo.redis.TxPipeline()

	reminderTaskBytes, err := json.Marshal(reminderTask)
	if err != nil {
		return fmt.Errorf("failed to marshal reminder task: %w", err)
	}

	if err = pipe.ZRem(ctx, reminderTasksKey, reminderTaskBytes).Err(); err != nil {
		return fmt.Errorf("failed to ZRem reminder task %v: %w", reminderTask.ID, err)
	}

	return nil
}
