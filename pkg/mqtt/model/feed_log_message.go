package model

import "time"

type FeedLogCollectionMessage struct {
	Value []FeedLogMessage `json:"value"`
}

type FeedLogMessage struct {
	Portions  uint      `json:"portions"`
	Timestamp time.Time `json:"timestamp"`
}
