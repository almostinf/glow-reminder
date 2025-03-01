package domain

import (
	"time"

	"github.com/google/uuid"
)

type Colour int8

const (
	UnknownColour Colour = 0
	Red           Colour = 1
	Green         Colour = 2
	Blue          Colour = 3
)

type Mode int8

const (
	UnknownMode Mode = 0
	Static      Mode = 1
	Blinking    Mode = 2
)

type Reminder struct {
	ID          uuid.UUID `db:"id"`
	UserID      int64     `db:"user_id"`
	Msg         string    `db:"msg"`
	Colour      Colour    `db:"colour"`
	Mode        Mode      `db:"mode"`
	ScheduledAt time.Time `db:"scheduled_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type ReminderTask struct {
	ID          uuid.UUID `json:"id"`
	ScheduledAt time.Time `json:"-"`
}

type GetRemindersParams struct {
	UserID int64
	Limit  uint64
	Offset uint64
}
