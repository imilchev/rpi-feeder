package v1

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/imilchev/rpi-feeder/pkg/mqtt/model"
	"github.com/imilchev/rpi-feeder/pkg/service/db/repos"
	"github.com/imilchev/rpi-feeder/pkg/service/models"
	"github.com/imilchev/rpi-feeder/pkg/service/mqtt"
	"github.com/imilchev/rpi-feeder/pkg/utils"
	"gorm.io/gorm"
)

type FeederController struct {
	repo repos.FeedersRepository
	mqtt mqtt.MqttManager
}

func NewFeederController(db *gorm.DB, mqtt mqtt.MqttManager) *FeederController {
	return &FeederController{mqtt: mqtt, repo: repos.NewFeedersRepository(db)}
}

func (c *FeederController) RegisterHandlers(a *fiber.App) {
	route := a.Group(apiGroup)
	route.Get("/feeder", c.GetFeeders)
	route.Post("/feeder/:clientId/feed", c.FeedPortions)
}

func (c *FeederController) GetFeeders(ctx *fiber.Ctx) error {
	feeders, err := c.repo.GetFeeders()
	if err != nil {
		return err
	}
	return ctx.Status(http.StatusOK).JSON(feeders)
}

func (c *FeederController) FeedPortions(ctx *fiber.Ctx) error {
	clientId := ctx.Params("clientId")
	if clientId == "" {
		return models.NewValidationError("Missing clientId.")
	}

	feeder, err := c.repo.GetFeederByClientId(clientId)
	if err != nil {
		return err
	}

	if feeder.Status != model.OnlineStatus {
		return models.NewValidationError(
			fmt.Sprintf("Feeder %s is not online.", feeder.ClientId))
	}

	request := models.FeedRequest{}
	if err := ctx.BodyParser(&request); err != nil {
		return models.NewValidationError(fmt.Sprintf("Cannot parse request body. %v", err))
	}

	if err := utils.Validate.Struct(request); err != nil {
		return models.NewValidationError(err.Error())
	}

	msg := model.FeedMessage{
		Portions: request.Portions,
	}
	if err := c.mqtt.SendFeedCommand(clientId, msg); err != nil {
		return err
	}
	return ctx.Status(http.StatusNoContent).JSON(fiber.Map{})
}
