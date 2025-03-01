package bot

import "github.com/almostinf/glow-reminder/internal/domain"

type state int

const (
	menuState           state = 0
	timeChoosingState   state = 1
	textEnteringState   state = 2
	colourChoosingState state = 3
	effectChoosingState state = 4
	listRemindersState  state = 5
)

type userState struct {
	s        state
	reminder domain.Reminder
	offset   int64
}
