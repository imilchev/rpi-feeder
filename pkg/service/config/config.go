package config

import "github.com/imilchev/rpi-feeder/pkg/mqtt/config"

type Config struct {
	Server   Server   `json:"server" validate:"required"`
	Database Database `json:"database" validate:"required"`
	//Jwt      Jwt      `json:"jwt" validate:"required"`
	Mqtt config.MqttConfig `json:"mqtt" validate:"required"`
}

type Server struct {
	Host        string `json:"host" validate:"required,url|ip"`
	Port        uint   `json:"port" validate:"gt=0"`
	ReadTimeout uint   `json:"readTimeout" validate:"gt=0"`
}

type Database struct {
	ConnectionString string `json:"connectionString" validate:"required"`
}

type Jwt struct {
	PublicKeyPath string `json:"publicKeyPath" validate:"required,file"`
	SigningMethod string `json:"signingMethod" validate:"required,oneof=RS256 RS384 RS512"`
}
