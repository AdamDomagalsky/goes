package entities

import "time"

// The Event struct contains the information about special events that
// the library is hosting, such as game nights and book signings.
type Event struct {
	ID          int
	Start       *time.Time
	Duration    *time.Duration
	Name        string
	Description string
}
