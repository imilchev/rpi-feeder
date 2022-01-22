package v1

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/imilchev/rpi-feeder/pkg/mqtt/model"
	"github.com/imilchev/rpi-feeder/pkg/service/middleware"
	"github.com/imilchev/rpi-feeder/pkg/service/models"
	"github.com/imilchev/rpi-feeder/tests/fake/mqtt"
	fake "github.com/imilchev/rpi-feeder/tests/fake/repos"
	"github.com/imilchev/rpi-feeder/tests/utils"
	modelUtils "github.com/imilchev/rpi-feeder/tests/utils/models"
	"github.com/stretchr/testify/suite"
)

type FeederControllerSuite struct {
	suite.Suite
	app      *fiber.App
	feeders  *fake.FakeFeedersRepository
	feedLogs *fake.FakeFeedLogsRepository
	mqtt     *mqtt.FakeServiceMqttManager
}

func (suite *FeederControllerSuite) SetupTest() {
	suite.app = fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})
	suite.feeders = &fake.FakeFeedersRepository{}
	suite.feedLogs = &fake.FakeFeedLogsRepository{}
	suite.mqtt = &mqtt.FakeServiceMqttManager{}
	c := FeederController{
		feedersRepo:  suite.feeders,
		feedLogsRepo: suite.feedLogs,
		mqtt:         suite.mqtt}
	c.RegisterHandlers(suite.app)
}

func (suite *FeederControllerSuite) TestGetFeeders() {
	fs := modelUtils.RandomFeeders()
	suite.feeders.Feeders = fs

	req := httptest.NewRequest(http.MethodGet, "/v1/feeders", nil)
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var rFs []models.Feeder
	suite.NoError(utils.ParseResponse(&rFs, resp))
	suite.ElementsMatch(fs, rFs)
}

func (suite *FeederControllerSuite) TestGetFeeders_Error() {
	suite.feeders.Error = fmt.Errorf("error")
	req := httptest.NewRequest(http.MethodGet, "/v1/feeders", nil)
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, resp.StatusCode)
}

func (suite *FeederControllerSuite) TestGetFeedLogsForFeeder() {
	fs := modelUtils.RandomFeeders()
	suite.feeders.Feeders = fs

	f := fs[len(fs)/2]
	ls := modelUtils.RandomFeedLogsForFeeder(f.ClientId)
	suite.feedLogs.FeedLogs = ls

	req := httptest.NewRequest(
		http.MethodGet, fmt.Sprintf("/v1/feeders/%s/logs", f.ClientId), nil)
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var rLs []models.FeedLog
	suite.NoError(utils.ParseResponse(&rLs, resp))
	suite.ElementsMatch(ls, rLs)
}

func (suite *FeederControllerSuite) TestGetFeedLogsForFeeder_FeederDoesNotExist() {
	clientId := utils.RandString(10)
	ls := modelUtils.RandomFeedLogsForFeeder(clientId)
	suite.feedLogs.FeedLogs = ls

	req := httptest.NewRequest(
		http.MethodGet, fmt.Sprintf("/v1/feeders/%s/logs", clientId), nil)
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, resp.StatusCode)
}

func (suite *FeederControllerSuite) TestFeedPortions() {
	f := modelUtils.RandomFeeder()
	f.Status = model.OnlineStatus
	suite.feeders.Feeders = append(suite.feeders.Feeders, f)

	m := models.FeedRequest{Portions: uint(rand.Intn(10) + 1)}
	req := utils.PostJsonRequest(fmt.Sprintf("/v1/feeders/%s/feed", f.ClientId), m)
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusNoContent, resp.StatusCode)

	suite.Equal(1, len(suite.mqtt.Feeds))
	suite.Equal(f.ClientId, suite.mqtt.Feeds[0].ClientId)
	suite.Equal(m.Portions, suite.mqtt.Feeds[0].Msg.Portions)
}

func (suite *FeederControllerSuite) TestFeedPortions_PortionsMissing() {
	f := modelUtils.RandomFeeder()
	f.Status = model.OfflineStatus
	suite.feeders.Feeders = append(suite.feeders.Feeders, f)

	m := models.FeedRequest{}
	req := utils.PostJsonRequest(fmt.Sprintf("/v1/feeders/%s/feed", f.ClientId), m)
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
	suite.Empty(suite.mqtt.Feeds)
}

func (suite *FeederControllerSuite) TestFeedPortions_FeederOffline() {
	f := modelUtils.RandomFeeder()
	f.Status = model.OfflineStatus
	suite.feeders.Feeders = append(suite.feeders.Feeders, f)

	m := models.FeedRequest{Portions: uint(rand.Intn(10) + 1)}
	req := utils.PostJsonRequest(fmt.Sprintf("/v1/feeders/%s/feed", f.ClientId), m)
	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	e := models.ApiError{}
	suite.NoError(utils.ParseResponse(&e, resp))
	suite.Equal(fmt.Sprintf("Feeder %s is not online.", f.ClientId), e.Message)
	suite.Empty(suite.mqtt.Feeds)
}

func TestFeederControllerSuite(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	suite.Run(t, new(FeederControllerSuite))
}
