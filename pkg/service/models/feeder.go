package models

import (
	"time"

	"github.com/imilchev/rpi-feeder/pkg/mqtt/model"
)

type Feeder struct {
	ClientId        string       `validate:"required,max=60"`
	SoftwareVersion string       `validate:"required,max=60"`
	Status          model.Status `validate:"required,max=7"`

	// The timestamp of when the feeder was last observed to be online.
	// Only set if the feeder is offline.
	LastOnline *time.Time
}

type FeedRequest struct {
	Portions uint `validate:"numeric,gt=0"`
}
