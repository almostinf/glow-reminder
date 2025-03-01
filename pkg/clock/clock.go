package clock

import "time"

//go:generate mockgen -package mocks -destination mocks/clock_mocks.go github.com/almostinf/glow-reminder/pkg/clock Clock

// Clock defines the interface for accessing time-related functionality.
type Clock interface {
	NowUTC() time.Time
	NowUnix() int64
}

type clock struct{}

// New creates a new instance of Clock.
func New() Clock {
	return &clock{}
}

// NowUTC returns the current time in UTC timezone.
func (clock) NowUTC() time.Time {
	return time.Now().UTC()
}

// NowUnix returns the current time as Unix timestamp.
func (clock) NowUnix() int64 {
	return time.Now().Unix()
}
