package config

type MqttConfig struct {
	Server            string `json:"server" validate:"required"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	ClientId          string `json:"clientId" validate:"required"`
	KeepAlive         uint16 `json:"keepAlive" validate:"required,gt=0"`
	ConnectRetryDelay uint16 `json:"connectRetryDelay" validate:"required,gt=0"`
}
