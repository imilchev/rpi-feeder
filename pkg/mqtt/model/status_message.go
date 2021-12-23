package model

type Status string

const (
	OnlineStatus  Status = "online"
	OfflineStatus Status = "offline"
)

type StatusMessage struct {
	SoftwareVersion string `json:"softwareVersion"`
	Status          Status `json:"status"`
}
