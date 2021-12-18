package servo

type ServoController interface {
	RotateClockwise()
	RotateCounterClockwise()
	Stop()
	Close()
}
