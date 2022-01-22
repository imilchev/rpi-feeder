package mqtt

import (
	"github.com/imilchev/rpi-feeder/pkg/mqtt/model"
)

type FeedRequests struct {
	ClientId string
	Msg      model.FeedMessage
}

// FakeServiceMqttManager provides an easy way of mocking a MqttManager.
// The functions in this fake implementation do not perform any validation.
type FakeServiceMqttManager struct {
	Feeds []FeedRequests

	// Error If this is set, any function will return it.
	Error error
}

func (m *FakeServiceMqttManager) SendFeedCommand(clientId string, msg model.FeedMessage) error {
	if m.Error != nil {
		return m.Error
	}

	m.Feeds = append(m.Feeds, FeedRequests{ClientId: clientId, Msg: msg})
	return nil
}

func (m *FakeServiceMqttManager) Stop() error {
	if m.Error != nil {
		return m.Error
	}
	return nil
}
