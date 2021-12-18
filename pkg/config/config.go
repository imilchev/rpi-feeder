package config

type Config struct {
	DbPath   string `json:"dbPath"`
	ServoPin uint8  `json:"servoPin"`

	// The amount of ms to spin in a direction to drop 1
	// portion of food.
	PortionMs int64 `json:"portionMs"`
}
