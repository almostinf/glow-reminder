package pg

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/almostinf/glow-reminder/internal/domain"
	"github.com/google/uuid"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

var nilTime = time.Time{}

func getRemindersQuery(params domain.GetRemindersParams) sq.SelectBuilder {
	query := psql.Select(
		"id",
		"user_id",
		"msg",
		"colour",
		"mode",
		"scheduled_at",
		"created_at",
		"updated_at",
	).
		From("reminders")

	if params.UserID != 0 {
		query = query.
			Where(sq.Eq{
				"user_id": params.UserID,
			})
	}

	if params.Offset != 0 {
		query = query.Offset(params.Offset)
	}

	if params.Limit != 0 {
		query = query.Limit(params.Limit)
	}

	return query.
		OrderBy("scheduled_at")
}

func getReminderQuery(id uuid.UUID) sq.SelectBuilder {
	return psql.Select(
		"id",
		"user_id",
		"msg",
		"colour",
		"mode",
		"scheduled_at",
		"created_at",
		"updated_at",
	).
		From("reminders").
		Where(sq.Eq{
			"id": id,
		})
}

func createReminderQuery(reminder domain.Reminder) sq.InsertBuilder {
	return psql.Insert("reminders").
		Columns(
			"id",
			"user_id",
			"msg",
			"colour",
			"mode",
			"scheduled_at",
			"created_at",
			"updated_at",
		).
		Values(
			reminder.ID,
			reminder.UserID,
			reminder.Msg,
			reminder.Colour,
			reminder.Mode,
			reminder.ScheduledAt,
			reminder.CreatedAt,
			reminder.UpdatedAt,
		)
}

func updateReminderQuery(reminder domain.Reminder) sq.UpdateBuilder {
	query := psql.Update("reminders")

	if reminder.Colour != domain.UnknownColour {
		query = query.Set("colour", reminder.Colour)
	}

	if domain.Colour(reminder.Mode) != domain.Colour(domain.UnknownMode) {
		query = query.Set("mode", reminder.Mode)
	}

	if reminder.ScheduledAt != nilTime {
		query = query.Set("scheduled_at", reminder.ScheduledAt)
	}

	return query.
		Set("updated_at", reminder.UpdatedAt).
		Where(sq.Eq{
			"id": reminder.ID,
		})
}

func deleteReminderQuery(id uuid.UUID) sq.DeleteBuilder {
	return psql.Delete("reminders").
		Where(sq.Eq{
			"id": id,
		})
}
