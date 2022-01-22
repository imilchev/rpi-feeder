package service

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/imilchev/rpi-feeder/pkg/mqtt/model"
	"github.com/imilchev/rpi-feeder/pkg/service/config"
	"github.com/imilchev/rpi-feeder/pkg/service/controllers"
	v1 "github.com/imilchev/rpi-feeder/pkg/service/controllers/v1"
	"github.com/imilchev/rpi-feeder/pkg/service/db"
	"github.com/imilchev/rpi-feeder/pkg/service/db/repos"
	"github.com/imilchev/rpi-feeder/pkg/service/middleware"
	"github.com/imilchev/rpi-feeder/pkg/service/models"
	"github.com/imilchev/rpi-feeder/pkg/service/mqtt"
	"github.com/imilchev/rpi-feeder/pkg/utils"
	"go.uber.org/zap"
)

type Service struct {
	config       config.Config
	app          *fiber.App
	db           *db.Database
	feedersRepo  repos.FeedersRepository
	feedLogsRepo repos.FeedLogsRepository
	mqtt         mqtt.MqttManager
	shutdownChan chan os.Signal
	controllers  []controllers.Controller
}

func NewService(configPath string) (*Service, error) {
	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		return nil, err
	}

	if err := utils.Validate.Struct(cfg); err != nil {
		return nil, err
	}

	// See: https://docs.gofiber.io/api/fiber#config
	fCfg := fiber.Config{
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		ErrorHandler: middleware.ErrorHandler,
	}
	db, err := db.NewDatabaseConnection(cfg.Database)
	if err != nil {
		return nil, err
	}

	app := &Service{
		config:       *cfg,
		app:          fiber.New(fCfg),
		db:           db,
		feedersRepo:  repos.NewFeedersRepository(db.DB),
		feedLogsRepo: repos.NewFeedLogsRepository(db.DB),
		shutdownChan: make(chan os.Signal, 1),
	}

	mqtt, err := mqtt.NewMqttManager(cfg.Mqtt, app.updateFeederStatus, app.storeFeedLogs)
	if err != nil {
		return nil, err
	}
	app.mqtt = mqtt
	app.controllers = []controllers.Controller{
		v1.NewFeederController(db.DB, mqtt),
	}

	signal.Notify(app.shutdownChan, os.Interrupt) // Catch OS signals.
	app.registerHandlers()
	return app, nil
}

// StartServerWithGracefulShutdown function for starting server with a graceful shutdown.
func (s *Service) StartServerWithGracefulShutdown() {
	defer zap.S().Sync() //nolint
	zap.S().Info("Starting RPi feeder web service...")
	// Create channel for idle connections.
	idleConnsClosed := make(chan struct{})

	go func() {
		<-s.shutdownChan
		zap.S().Info("Shutting down RPi feeder web service...")
		// Received an interrupt signal, shutdown.
		if err := s.app.Shutdown(); err != nil {
			// Error from closing listeners, or context timeout:
			zap.S().Errorf("Oops... Server is not shutting down! Reason: %v", err)
		}
		close(idleConnsClosed)
	}()

	// Run server.
	if err := s.app.Listen(s.connUrl()); err != nil {
		zap.S().Errorf("Oops... Server is not running! Reason: %v", err)
	}

	<-idleConnsClosed
	zap.S().Info("Sucessfully closed all API connections.")

	if err := s.db.Close(); err != nil {
		return
	}

	if err := s.mqtt.Stop(); err != nil {
		zap.S().Errorf("Failed to gracefully shutdown MQTT client. %+v", err)
	}
	zap.S().Info("RPi feeder web service gracefully shut down!")
}

func (a *Service) connUrl() string {
	return fmt.Sprintf("%s:%d", a.config.Server.Host, a.config.Server.Port)
}

func (a *Service) registerHandlers() {
	// jwtHandler, err := middleware.JwtHandler(a.config.Jwt)
	// if err != nil {
	// 	panic(err)
	// }

	// // Register the JWT handler. Every handler registered below this point
	// // will require JWT auth.
	// a.app.Use(jwtHandler)
	for _, c := range a.controllers {
		c.RegisterHandlers(a.app)
	}
}

func (s *Service) updateFeederStatus(clientId string, msg model.StatusMessage) error {
	m := models.Feeder{
		ClientId:        clientId,
		SoftwareVersion: msg.SoftwareVersion,
		Status:          msg.Status,
	}
	if msg.Status == model.OfflineStatus {
		t := time.Now().UTC().Unix()
		m.LastOnline = &t
	}

	_, err := s.feedersRepo.GetFeederByClientId(clientId)
	if err != nil {
		if _, err := s.feedersRepo.CreateFeeder(m); err != nil {
			return err
		}
		return nil
	}

	_, err = s.feedersRepo.UpdateFeeder(m)
	return err
}

func (s *Service) storeFeedLogs(clientId string, msg model.FeedLogCollectionMessage) error {
	_, err := s.feedersRepo.GetFeederByClientId(clientId)
	if err != nil {
		return err
	}

	var f []models.FeedLog
	for _, m := range msg.Value {
		f = append(f, models.FeedLog{
			ClientId:  clientId,
			Portions:  m.Portions,
			Timestamp: m.Timestamp.UTC().Unix(),
		})
	}

	_, err = s.feedLogsRepo.CreateFeedLogs(f)
	return err
}
