package feeder

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/imilchev/rpi-feeder/pkg/config"
	"github.com/imilchev/rpi-feeder/pkg/db"
	"github.com/imilchev/rpi-feeder/pkg/servo"
	"go.uber.org/zap"
)

type FeederManager struct {
	config *config.Config
}

func NewFeederManager(configPath string) (*FeederManager, error) {
	config, err := config.ReadConfig(configPath)
	if err != nil {
		return nil, err
	}

	return &FeederManager{config: config}, nil
}

func (fm *FeederManager) Start() error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	zap.S().Info("Feeder started.")
	dbManager, err := db.NewDbManager(fm.config.DbPath)
	if err != nil {
		return err
	}

	servoController, err := servo.NewServoController(17)
	if err != nil {
		return err
	}

	servoController.RotateClockwise()
	time.Sleep(3 * time.Second)
	zap.S().Debugf("Sending stop signal...")
	servoController.Stop()

	servoController.RotateCounterClockwise()

	for {
		//select {
		//case <-interrupt:
		<-interrupt
		zap.S().Info("Shutting down...")

		servoController.Stop()
		servoController.Close()
		dbManager.Close()

		zap.S().Info("Exit")
		return nil
		//}
	}
}
