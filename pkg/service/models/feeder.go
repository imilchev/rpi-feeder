package models

import "time"

type Feeder struct {
	ClientId        string
	SoftwareVersion string
	Status          string

	// The timestamp of when the feeder was last observed to be online.
	// Only set if the feeder is offline.
	LastOnline *time.Time
}
