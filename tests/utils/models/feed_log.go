package models

import (
	"math/rand"
	"time"

	dbm "github.com/imilchev/rpi-feeder/pkg/service/db/models"
	"github.com/imilchev/rpi-feeder/pkg/service/models"
)

func RandomFeedLogForFeeder(clientId string) models.FeedLog {
	return models.FeedLog{
		ClientId:  clientId,
		Portions:  uint(rand.Intn(10) + 1),
		Timestamp: time.Now().UTC().Unix(),
	}
}

func RandomFeedLogsForFeeder(clientId string) []models.FeedLog {
	var f []models.FeedLog
	count := rand.Intn(15) + 1

	for i := 0; i < count; i++ {
		f = append(f, RandomFeedLogForFeeder(clientId))
	}
	return f
}

func RandomDbFeedLogsForFeeder(clientId string) []dbm.FeedLog {
	var f []dbm.FeedLog
	count := rand.Intn(15) + 1

	for i := 0; i < count; i++ {
		f = append(f, dbm.FeedLog{
			ClientId:  clientId,
			Portions:  uint(rand.Intn(10) + 1),
			Timestamp: time.Unix(time.Now().UTC().Unix(), 0),
		})
	}
	return f
}
