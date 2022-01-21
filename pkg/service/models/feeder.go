package models

import (
	"time"

	"github.com/imilchev/rpi-feeder/pkg/mqtt/model"
)

type Feeder struct {
	ClientId        string
	SoftwareVersion string
	Status          model.Status

	// The timestamp of when the feeder was last observed to be online.
	// Only set if the feeder is offline.
	LastOnline *time.Time
}

type FeedRequest struct {
	Portions uint `validate:"numeric,gt=0"`
}
