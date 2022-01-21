package model

import "time"

type FeedLog struct {
	Id        int       `json:"id"`
	Portions  uint      `json:"portions"`
	Timestamp time.Time `json:"timestamp"`
}
