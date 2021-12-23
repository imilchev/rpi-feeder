package config

type Config struct {
	DbPath   string `json:"dbPath"`
	ServoPin uint8  `json:"servoPin" validate:"gt=0"`

	// The amount of ms to spin in a direction to drop 1
	// portion of food.
	PortionMs uint64     `json:"portionMs" validate:"gt=0"`
	Mqtt      MqttConfig `json:"mqtt" validate:"required"`
}

type MqttConfig struct {
	Server            string `json:"server" validate:"required"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	ClientId          string `json:"clientId" validate:"required"`
	KeepAlive         uint16 `json:"keepAlive" validate:"required,gt=0"`
	ConnectRetryDelay uint16 `json:"connectRetryDelay" validate:"required,gt=0"`
}
