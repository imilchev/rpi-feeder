package config

import "github.com/imilchev/rpi-feeder/pkg/mqtt/config"

type Config struct {
	DbPath   string `json:"dbPath"`
	ServoPin uint8  `json:"servoPin" validate:"gt=0"`

	// The amount of ms to spin in a direction to drop 1
	// portion of food.
	PortionMs uint64            `json:"portionMs" validate:"gt=0"`
	Mqtt      config.MqttConfig `json:"mqtt" validate:"required"`
}
