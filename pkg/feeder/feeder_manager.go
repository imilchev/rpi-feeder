package feeder

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/imilchev/rpi-feeder/pkg/feeder/db"
	"github.com/imilchev/rpi-feeder/pkg/feeder/servo"
	"go.uber.org/zap"
)

type FeederManager struct {
}

func NewFeederManager() *FeederManager {
	return &FeederManager{}
}

func (fm *FeederManager) Start(dbPath string) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	zap.S().Info("Feeder started.")
	dbManager, err := db.NewDbManager(dbPath)
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
		select {
		case <-interrupt:
			zap.S().Info("Shutting down...")

			servoController.Close()
			dbManager.Close()

			zap.S().Info("Exit")
			return nil
		}
	}
}
