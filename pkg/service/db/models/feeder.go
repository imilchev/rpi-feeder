package models

import (
	"time"

	"github.com/imilchev/rpi-feeder/pkg/mqtt/model"
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
	m.Status = model.Status(f.Status)
	m.LastOnline = nil
	if f.LastOnline != nil {
		t := f.LastOnline.UTC().Unix()
		m.LastOnline = &t
	}
}

func (f *Feeder) FromApi(m models.Feeder) {
	f.ClientId = m.ClientId
	f.SoftwareVersion = m.SoftwareVersion
	f.Status = string(m.Status)
	f.LastOnline = nil
	if m.LastOnline != nil {
		t := time.Unix(*m.LastOnline, 0)
		f.LastOnline = &t
	}
}
