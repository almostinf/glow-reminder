package pg

import (
	"context"
	"fmt"

	"github.com/almostinf/glow-reminder/internal/domain"
	"github.com/almostinf/glow-reminder/pkg/logger"
	"github.com/almostinf/glow-reminder/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var _ ReminderRepo = (*reminderRepo)(nil)

type ReminderRepo interface {
	GetReminders(ctx context.Context, params domain.GetRemindersParams) ([]*domain.Reminder, error)
	GetReminder(ctx context.Context, id uuid.UUID) (*domain.Reminder, error)
	CreateReminder(ctx context.Context, reminder domain.Reminder) error
	UpdateReminder(ctx context.Context, reminder domain.Reminder) error
	DeleteReminder(ctx context.Context, id uuid.UUID) error
}

type reminderRepo struct {
	pg     *postgres.Postgres
	logger logger.Logger
}

func NewReminderRepo(pg *postgres.Postgres, logger logger.Logger) *reminderRepo {
	return &reminderRepo{
		pg:     pg,
		logger: logger,
	}
}

func (repo *reminderRepo) GetReminders(ctx context.Context, params domain.GetRemindersParams) ([]*domain.Reminder, error) {
	conn := repo.pg.GetTransactionConn(ctx)

	query := getRemindersQuery(params)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql query: %w", err)
	}

	rows, err := conn.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	reminders, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Reminder])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows: %w", err)
	}

	reminderPtrs := make([]*domain.Reminder, 0, len(reminders))
	for i := range reminders {
		reminderPtrs = append(reminderPtrs, &reminders[i])
	}

	return reminderPtrs, nil
}

func (repo *reminderRepo) GetReminder(ctx context.Context, id uuid.UUID) (*domain.Reminder, error) {
	conn := repo.pg.GetTransactionConn(ctx)

	query := getReminderQuery(id)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql query: %w", err)
	}

	rows, err := conn.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	reminder, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Reminder])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows: %w", err)
	}

	return &reminder, nil
}

func (repo *reminderRepo) CreateReminder(ctx context.Context, reminder domain.Reminder) error {
	conn := repo.pg.GetTransactionConn(ctx)

	query := createReminderQuery(reminder)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to get sql query: %w", err)
	}

	if _, err = conn.Exec(ctx, sqlQuery, args...); err != nil {
		return fmt.Errorf("failed to Exec: %w", err)
	}

	return nil
}

func (repo *reminderRepo) UpdateReminder(ctx context.Context, reminder domain.Reminder) error {
	conn := repo.pg.GetTransactionConn(ctx)

	query := updateReminderQuery(reminder)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to get sql query: %w", err)
	}

	if _, err = conn.Exec(ctx, sqlQuery, args...); err != nil {
		return fmt.Errorf("failed to Exec: %w", err)
	}

	return nil
}

func (repo *reminderRepo) DeleteReminder(ctx context.Context, id uuid.UUID) error {
	conn := repo.pg.GetTransactionConn(ctx)

	query := deleteReminderQuery(id)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to get sql query: %w", err)
	}

	if _, err = conn.Exec(ctx, sqlQuery, args...); err != nil {
		return fmt.Errorf("failed to Exec: %w", err)
	}

	return nil
}
