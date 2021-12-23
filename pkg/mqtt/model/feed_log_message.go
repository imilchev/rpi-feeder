package model

import "time"

type FeedLogMessage struct {
	Portions  uint      `json:"portions"`
	Timestamp time.Time `json:"timestamp"`
}
