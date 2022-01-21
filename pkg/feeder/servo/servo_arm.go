package servo

import (
	"time"

	"github.com/stianeikeland/go-rpio/v4"
	"go.uber.org/zap"
)

type servoController struct {
	pin         rpio.Pin
	stopChan    chan struct{}
	stoppedChan chan struct{}
	isRotating  bool
}

func NewServoController(pinNumber uint8) (ServoController, error) {
	err := rpio.Open()
	if err != nil {
		zap.S().Errorf("Could not initialize GPIO library. %+v", err)
		return nil, err
	}

	pin := rpio.Pin(pinNumber)
	pin.Mode(rpio.Output)
	zap.S().Infof(
		"Initialized GPIO library. Using pin %d to control the servo.", uint8(pin))
	return &servoController{
		pin:         pin,
		stopChan:    make(chan struct{}),
		stoppedChan: make(chan struct{}),
	}, nil
}

func (sc *servoController) RotateClockwise() {
	zap.S().Debug("Rotating servo clockwise...")
	sc.isRotating = true
	go func() {
		for {
			select {
			default:
				sc.pin.High()
				time.Sleep(2 * time.Millisecond)
				sc.pin.Low()
				time.Sleep(17 * time.Millisecond)
			case <-sc.stopChan:
				zap.S().Debug("Servo rotation stopped.")
				sc.stoppedChan <- struct{}{}
				return
			}
		}
	}()
}

func (sc *servoController) RotateCounterClockwise() {
	zap.S().Debug("Rotating servo counter-clockwise...")
	sc.isRotating = true
	go func() {
		for {
			select {
			default:
				sc.pin.High()
				time.Sleep(1 * time.Millisecond)
				sc.pin.Low()
				time.Sleep(19 * time.Millisecond)
			case <-sc.stopChan:
				zap.S().Debug("Servo rotation stopped.")
				sc.stoppedChan <- struct{}{}
				return
			}
		}
	}()
}

func (sc *servoController) Stop() {
	if sc.isRotating {
		sc.stopChan <- struct{}{}
		<-sc.stoppedChan
		sc.isRotating = false
	}
}

func (sc *servoController) Close() {
	if sc.isRotating {
		sc.Stop()
	}
	close(sc.stopChan)
	close(sc.stoppedChan)
	if err := rpio.Close(); err != nil {
		zap.S().Errorf("Failed to close GPIO library. %+v", err)
	}
	zap.S().Info("GPIO library closed.")
}

func Move() {
	err := rpio.Open()
	if err != nil {
		zap.S().Errorf("Bla %+v", err)
	}
	defer rpio.Close()

	zap.S().Debug("Initialized GPIO library.")

	pin := rpio.Pin(17)
	pin.Mode(rpio.Output)

	// ccw
	for i := 0; i < 250; i++ {
		pin.High()
		time.Sleep(1 * time.Millisecond)
		pin.Low()
		time.Sleep(19 * time.Millisecond)
	}
	zap.S().Debugf("Changing direction")

	// cw
	zap.S().Debugf("Rotating now")
	for i := 0; i < 250; i++ {
		pin.High()
		time.Sleep(2 * time.Millisecond)
		pin.Low()
		time.Sleep(18 * time.Millisecond)
	}
}
