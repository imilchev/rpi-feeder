package models

import (
	"time"

	"github.com/imilchev/rpi-feeder/pkg/service/models"
)

type Feeder struct {
	ClientId        string `gorm:"primaryKey"`
	SoftwareVersion string
	Status          string

	// The timestamp of when the feeder was last observed to be online.
	// Only set if the feeder is offline.
	LastOnline *time.Time
}

func (f Feeder) ToApi(m *models.Feeder) {
	m.ClientId = f.ClientId
	m.SoftwareVersion = f.SoftwareVersion
	m.Status = f.Status
	m.LastOnline = f.LastOnline
}

func (f *Feeder) FromApi(m models.Feeder) {
	f.ClientId = m.ClientId
	f.SoftwareVersion = m.SoftwareVersion
	f.Status = m.Status
	f.LastOnline = m.LastOnline
}
