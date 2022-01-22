package models

import (
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
		t := time.Now().UTC()
		f.LastOnline = &t
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