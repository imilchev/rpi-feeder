package models

import (
	"time"

	"github.com/imilchev/rpi-feeder/pkg/service/models"
)

type FeedLog struct {
	Id        int `gorm:"primaryKey"`
	ClientId  string
	Portions  uint
	Timestamp time.Time
}

func (f FeedLog) ToApi(m *models.FeedLog) {
	m.Id = f.Id
	m.ClientId = f.ClientId
	m.Portions = f.Portions
	m.Timestamp = f.Timestamp.UTC().Unix()
}

func (f *FeedLog) FromApi(m models.FeedLog) {
	f.Id = m.Id
	f.ClientId = m.ClientId
	f.Portions = m.Portions
	f.Timestamp = time.Unix(m.Timestamp, 0)
}
