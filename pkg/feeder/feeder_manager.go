package feeder

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/imilchev/rpi-feeder/pkg/config"
	"github.com/imilchev/rpi-feeder/pkg/db"
	dbm "github.com/imilchev/rpi-feeder/pkg/db/model"
	"github.com/imilchev/rpi-feeder/pkg/mqtt"
	"github.com/imilchev/rpi-feeder/pkg/mqtt/model"
	"github.com/imilchev/rpi-feeder/pkg/servo"
	"github.com/imilchev/rpi-feeder/pkg/utils"
	"go.uber.org/zap"
)

type FeederManager struct {
	config          *config.Config
	dbManager       db.DbManager
	servoController servo.ServoController
	mqttManager     mqtt.MqttManager
}

func NewFeederManager(configPath string) (*FeederManager, error) {
	config, err := config.ReadConfig(configPath)
	if err != nil {
		return nil, err
	}

	if err := utils.Validate.Struct(config); err != nil {
		return nil, err
	}

	dbManager, err := db.NewDbManager(config.DbPath)
	if err != nil {
		return nil, err
	}

	servoController, err := servo.NewServoController(config.ServoPin)
	if err != nil {
		return nil, err
	}

	fm := &FeederManager{
		config:          config,
		dbManager:       dbManager,
		servoController: servoController,
	}

	fm.mqttManager, err = mqtt.NewMqttManager(
		config.Mqtt,
		func(msg model.FeedMessage) error { return fm.feed(msg.Portions) })
	if err != nil {
		return nil, err
	}

	return fm, nil
}

func (fm *FeederManager) Start() error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	zap.S().Info("Feeder started.")

	if err := fm.flushFeedLog(); err != nil {
		zap.S().Error("Failed to flush feed log. %v", err)
	}

	// fm.servoController.RotateClockwise()
	// time.Sleep(3 * time.Second)
	// zap.S().Debugf("Sending stop signal...")
	// fm.servoController.Stop()

	// fm.servoController.RotateCounterClockwise()

	//for {
	//select {
	//case <-interrupt:
	<-interrupt
	zap.S().Info("Shutting down...")

	fm.servoController.Stop()
	fm.servoController.Close()
	fm.dbManager.Close()
	if err := fm.mqttManager.Stop(); err != nil {
		zap.S().Errorf("Failed to stop MQTT manager. %+v", err)
	}

	zap.S().Info("Exit")
	return nil
}

func (fm *FeederManager) flushFeedLog() error {
	feedLog, err := fm.dbManager.ListFeedLog()
	if err != nil {
		zap.S().Error("Failed to list feed log.")
		return err
	}

	if len(feedLog) == 0 {
		zap.S().Info("No feed logs have to be flushed.")
		return nil
	}

	zap.S().Infof("Found %d feed logs to be flushed.", len(feedLog))
	msg := model.FeedLogCollectionMessage{}
	for _, f := range feedLog {
		fmsg := model.FeedLogMessage{Portions: f.Portions, Timestamp: f.Timestamp}
		msg.Value = append(msg.Value, fmsg)
	}
	if err := fm.mqttManager.SendFeedLog(msg); err != nil {
		zap.S().Error("Failed to send feed log.")
		return err
	}
	zap.S().Info("Feed log flushed.")
	return fm.dbManager.CleanFeedLog()
}

func (fm *FeederManager) feed(portions uint) error {
	zap.S().Debugf("Serving %d portions...", portions)
	fm.servoController.RotateClockwise()

	for i := uint(0); i < portions; i++ {
		time.Sleep(time.Duration(fm.config.PortionMs) * time.Millisecond)
	}
	fm.servoController.Stop()
	zap.S().Infof("Served %d portions.", portions)

	// send the log via mqtt and if that fails store it locally
	msg := model.FeedLogCollectionMessage{
		Value: []model.FeedLogMessage{
			{Portions: portions, Timestamp: time.Now().UTC()},
		},
	}
	if err := fm.mqttManager.SendFeedLog(msg); err != nil {
		zap.S().Warnf("Failed to send feed log to server. %v", err)

		feedLog := dbm.FeedLog{
			Portions:  msg.Value[0].Portions,
			Timestamp: msg.Value[0].Timestamp,
		}
		return fm.dbManager.AddFeedLog(feedLog)
	}
	return nil
}
