package models

import (
	"math/rand"
	"time"

	"github.com/imilchev/rpi-feeder/pkg/mqtt/model"
	dbm "github.com/imilchev/rpi-feeder/pkg/service/db/models"
	"github.com/imilchev/rpi-feeder/pkg/service/models"
	"github.com/imilchev/rpi-feeder/tests/utils"
)

func RandomFeeder() models.Feeder {
	isOffline := utils.RandBool()
	f := models.Feeder{
		ClientId:        utils.RandString(10),
		SoftwareVersion: utils.RandString(10),
		Status:          model.OnlineStatus,
	}
	if isOffline {
		f.Status = model.OfflineStatus
		t := time.Now().UTC().Unix()
		f.LastOnline = &t
	}
	return f
}

func RandomFeeders() []models.Feeder {
	var f []models.Feeder
	count := rand.Intn(10) + 1
	for i := 0; i < count; i++ {
		f = append(f, RandomFeeder())
	}
	return f
}

func RandomDbFeeder() dbm.Feeder {
	isOffline := utils.RandBool()
	f := dbm.Feeder{
		ClientId:        utils.RandString(10),
		SoftwareVersion: utils.RandString(10),
		Status:          string(model.OnlineStatus),
	}
	if isOffline {
		f.Status = string(model.OfflineStatus)
		t := time.Now().UTC()
		f.LastOnline = &t
	}
	return f
}

// func CompareFeedersLastOnline(suite suite.Suite, expected, actual *models.Feeder) {
// 	if expected.LastOnline != nil {
// 		//suite.Equal(expected.LastOnline.Unix(), actual.LastOnline.Unix())
// 		expected.LastOnline = nil
// 		actual.LastOnline = nil
// 	}
// }
