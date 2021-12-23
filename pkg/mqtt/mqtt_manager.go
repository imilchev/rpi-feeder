package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/imilchev/rpi-feeder/pkg/config"
	"github.com/imilchev/rpi-feeder/pkg/mqtt/model"
	"go.uber.org/zap"
)

type FeedHandler func(model.FeedMessage) error

type MqttManager interface {
	SendFeedLog(msg model.FeedLogMessage) error
	Stop() error
}

type mqttManager struct {
	clientId string
	c        *autopaho.ConnectionManager
}

func NewMqttManager(cfg config.MqttConfig, fh FeedHandler) (MqttManager, error) {
	serverUrl, err := url.Parse(cfg.Server)
	if err != nil {
		return nil, err
	}

	router := paho.NewStandardRouter()
	router.RegisterHandler(
		fmt.Sprintf("feeder/%s/feed", cfg.ClientId),
		func(p *paho.Publish) { internalFeedHandler(p, fh) })

	pahoCfg := autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{serverUrl},
		KeepAlive:         cfg.KeepAlive,
		ConnectRetryDelay: time.Duration(cfg.ConnectRetryDelay) * time.Second,
		OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
			zap.S().Info("MQTT connection is up.")
			msg := model.StatusMessage{SoftwareVersion: "dev", Status: model.OnlineStatus}
			if err := sendStatusMessage(msg, cm, cfg.ClientId); err != nil {
				zap.S().Errorf("Failed to send status message. %v", err)
				return
			}

			if _, err := cm.Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: map[string]paho.SubscribeOptions{
					feedTopic(cfg.ClientId): {QoS: byte(2)},
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

	willMsg := model.StatusMessage{SoftwareVersion: "dev", Status: model.OfflineStatus}
	willData, err := json.Marshal(willMsg)
	if err != nil {
		return nil, err
	}
	pahoCfg.SetWillMessage(statusTopic(cfg.ClientId), willData, byte(1), true)

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
	msg := model.StatusMessage{SoftwareVersion: "dev", Status: model.OfflineStatus}
	if err := sendStatusMessage(msg, m.c, m.clientId); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err := m.c.Disconnect(ctx)
	if err == nil {
		zap.S().Info("MQTT connection closed.")
	}
	return err
}

func (m *mqttManager) SendFeedLog(msg model.FeedLogMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = m.c.Publish(context.Background(), &paho.Publish{
		Topic:   feedLogTopic(m.clientId),
		QoS:     byte(2),
		Payload: data,
	})
	return err
}

func sendStatusMessage(msg model.StatusMessage, cm *autopaho.ConnectionManager, clientId string) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = cm.Publish(context.Background(), &paho.Publish{
		Topic:   statusTopic(clientId),
		QoS:     byte(1),
		Retain:  true,
		Payload: data,
	})
	return err
}

func internalFeedHandler(p *paho.Publish, fh FeedHandler) {
	msg := model.FeedMessage{}
	if err := json.Unmarshal(p.Payload, &msg); err != nil {
		zap.S().Errorf("Failed to deserialize message %s. %v", string(p.Payload), err)
		return
	}
	if err := fh(msg); err != nil {
		zap.S().Errorf("Failed to feed %d portions. %v", msg.Portions, err)
		return
	}
	zap.S().Infof("Feed %d portions", msg.Portions)
}
