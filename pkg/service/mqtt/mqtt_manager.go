package mqtt

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/imilchev/rpi-feeder/pkg/mqtt"
	"github.com/imilchev/rpi-feeder/pkg/mqtt/config"
	"github.com/imilchev/rpi-feeder/pkg/mqtt/model"
	"go.uber.org/zap"
)

type FeederStatusHandler func(clientId string, msg model.StatusMessage) error

type MqttManager interface {
	SendFeedCommand(clientId string, msg model.FeedMessage) error
	Stop() error
}

type mqttManager struct {
	clientId string
	c        *autopaho.ConnectionManager
}

func NewMqttManager(cfg config.MqttConfig, fsh FeederStatusHandler) (MqttManager, error) {
	serverUrl, err := url.Parse(cfg.Server)
	if err != nil {
		return nil, err
	}

	router := paho.NewStandardRouter()
	router.RegisterHandler(
		mqtt.StatusTopic(nil),
		func(p *paho.Publish) { internalStatusHandler(p, fsh) })

	pahoCfg := autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{serverUrl},
		KeepAlive:         cfg.KeepAlive,
		ConnectRetryDelay: time.Duration(cfg.ConnectRetryDelay) * time.Second,
		OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
			zap.S().Info("MQTT connection is up.")
			// msg := model.StatusMessage{SoftwareVersion: "dev", Status: model.OnlineStatus}
			// if err := sendStatusMessage(msg, cm, cfg.ClientId); err != nil {
			// 	zap.S().Errorf("Failed to send status message. %v", err)
			// 	return
			// }

			if _, err := cm.Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: map[string]paho.SubscribeOptions{
					mqtt.StatusTopic(nil): {QoS: byte(1)},
				},
			}); err != nil {
				zap.S().Errorf("Failed to subscribe (%v). This is likely to mean no messages will be received.", err)
				return
			}
			zap.S().Info("MQTT subscriptions are made.")
		},
		OnConnectError: func(err error) { zap.S().Warnf("Error whilst attempting connection: %v", err) },
		ClientConfig: paho.ClientConfig{
			ClientID:      cfg.ClientId,
			Router:        router,
			OnClientError: func(err error) { zap.S().Errorf("Client error: %s", err) },
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					zap.S().Warnf("Server requested disconnect: %s", d.Properties.ReasonString)
				} else {
					zap.S().Warnf("Server requested disconnect; reason code: %d", d.ReasonCode)
				}
			},
		},
	}
	pahoCfg.SetUsernamePassword(cfg.Username, []byte(cfg.Password))

	// willMsg := model.StatusMessage{SoftwareVersion: "dev", Status: model.OfflineStatus}
	// willData, err := json.Marshal(willMsg)
	// if err != nil {
	// 	return nil, err
	// }
	// pahoCfg.SetWillMessage(mqtt.StatusTopic(&cfg.ClientId), willData, byte(1), true)

	// Connect to the broker
	cm, err := autopaho.NewConnection(context.Background(), pahoCfg)
	if err != nil {
		return nil, err
	}

	m := &mqttManager{clientId: cfg.ClientId, c: cm}
	if err := cm.AwaitConnection(context.Background()); err != nil {
		return nil, err
	}

	zap.S().Infof("Connected to %s.", cfg.Server)
	return m, nil
}

func (m *mqttManager) Stop() error {
	// msg := model.StatusMessage{SoftwareVersion: "dev", Status: model.OfflineStatus}
	// if err := sendStatusMessage(msg, m.c, m.clientId); err != nil {
	// 	return err
	// }

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err := m.c.Disconnect(ctx)
	if err == nil {
		zap.S().Info("MQTT connection closed.")
	}
	return err
}

func (m *mqttManager) SendFeedCommand(clientId string, msg model.FeedMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = m.c.Publish(context.Background(), &paho.Publish{
		Topic:   mqtt.FeedTopic(&clientId),
		QoS:     byte(2),
		Payload: data,
	})
	return err
}

func internalStatusHandler(p *paho.Publish, fsh FeederStatusHandler) {
	msg := model.StatusMessage{}
	if err := json.Unmarshal(p.Payload, &msg); err != nil {
		zap.S().Errorf("Failed to deserialize message %s. %v", string(p.Payload), err)
		return
	}
	clientId := mqtt.ClientIdFromTopic(p.Topic)
	if err := fsh(clientId, msg); err != nil {
		zap.S().Errorf("Failed to set status for feeder %s. %v", clientId, err)
		return
	}
	zap.S().Infof("Status of feeder %s set to %s", clientId, msg.Status)
}
