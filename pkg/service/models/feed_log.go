package models

import "time"

type FeedLog struct {
	Id        int
	ClientId  string    `validate:"required"`
	Portions  uint      `validate:"required,gt=0"`
	Timestamp time.Time `validate:"required"`
}
