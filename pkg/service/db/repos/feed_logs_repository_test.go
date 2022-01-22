package repos

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	dbm "github.com/imilchev/rpi-feeder/pkg/service/db/models"
	"github.com/imilchev/rpi-feeder/pkg/service/models"
	"github.com/imilchev/rpi-feeder/tests/utils"
	modelUtils "github.com/imilchev/rpi-feeder/tests/utils/models"
	"github.com/stretchr/testify/suite"
)

type FeedLogsRepositorySuite struct {
	suite.Suite
	r *feedLogsRepository
}

func (suite *FeedLogsRepositorySuite) SetupTest() {
	suite.Require().NoError(utils.InitTestDb())
	db, err := utils.GetTestDb()
	suite.Require().NoError(err)
	suite.r = &feedLogsRepository{db: db}
}

func (suite *FeedLogsRepositorySuite) AfterTest(suiteName, testName string) {
	suite.Require().NoError(utils.CleanupDb(suite.r.db))
	db, err := suite.r.db.DB()
	suite.Require().NoError(err)
	db.Close()
}

func (suite *FeedLogsRepositorySuite) TestCreateFeedLogs() {
	f := modelUtils.RandomDbFeeder()
	suite.NoError(suite.r.db.Create(&f).Error)

	logs := modelUtils.RandomFeedLogsForFeeder(f.ClientId)
	ll, err := suite.r.CreateFeedLogs(logs)
	suite.NoError(err)

	// Do not compare IDs since they were generated by the database.
	for i := range ll {
		ll[i].Id = 0
	}
	suite.Equal(logs, ll)
}

func (suite *FeedLogsRepositorySuite) TestCreateFeedLogs_ClientIdMissing() {
	f := modelUtils.RandomDbFeeder()
	suite.NoError(suite.r.db.Create(&f).Error)

	l := modelUtils.RandomFeedLogForFeeder(f.ClientId)
	l.ClientId = ""
	_, err := suite.r.CreateFeedLogs([]models.FeedLog{l})
	suite.Error(err)
	apiErr, ok := err.(*models.ApiError)
	suite.True(ok)
	suite.Equal(http.StatusBadRequest, apiErr.Code())
}

func (suite *FeedLogsRepositorySuite) TestCreateFeedLogs_ClientIdTooLong() {
	f := modelUtils.RandomDbFeeder()
	suite.NoError(suite.r.db.Create(&f).Error)

	l := modelUtils.RandomFeedLogForFeeder(f.ClientId)
	l.ClientId = utils.RandString(61)
	_, err := suite.r.CreateFeedLogs([]models.FeedLog{l})
	suite.Error(err)
	apiErr, ok := err.(*models.ApiError)
	suite.True(ok)
	suite.Equal(http.StatusBadRequest, apiErr.Code())
}

func (suite *FeedLogsRepositorySuite) TestCreateFeedLogs_PortionsMissing() {
	f := modelUtils.RandomDbFeeder()
	suite.NoError(suite.r.db.Create(&f).Error)

	l := modelUtils.RandomFeedLogForFeeder(f.ClientId)
	l.Portions = 0
	_, err := suite.r.CreateFeedLogs([]models.FeedLog{l})
	suite.Error(err)
	apiErr, ok := err.(*models.ApiError)
	suite.True(ok)
	suite.Equal(http.StatusBadRequest, apiErr.Code())
}

func (suite *FeedLogsRepositorySuite) TestCreateFeedLogs_TimestampMissing() {
	f := modelUtils.RandomDbFeeder()
	suite.NoError(suite.r.db.Create(&f).Error)

	l := modelUtils.RandomFeedLogForFeeder(f.ClientId)
	l.Timestamp = 0
	_, err := suite.r.CreateFeedLogs([]models.FeedLog{l})
	suite.Error(err)
	apiErr, ok := err.(*models.ApiError)
	suite.True(ok)
	suite.Equal(http.StatusBadRequest, apiErr.Code())
}

func (suite *FeedLogsRepositorySuite) TestGetFeedLogs() {
	feeders, logs := suite.seedFeedLogs()
	randFeeder := feeders[len(feeders)/2]

	var expected []models.FeedLog
	apiFeeder := &models.FeedLog{}
	for _, f := range logs {
		if f.ClientId == randFeeder.ClientId {
			f.ToApi(apiFeeder)
			expected = append(expected, *apiFeeder)
		}
	}

	feedLogs, err := suite.r.GetLogsForFeeder(randFeeder.ClientId)
	suite.NoError(err)
	suite.ElementsMatch(expected, feedLogs)
}

func (suite *FeedLogsRepositorySuite) TestGetFeedLogs_DoesNotExist() {
	clientId := utils.RandString(10)

	_, err := suite.r.GetLogsForFeeder(clientId)
	suite.Error(err)
	apiErr, ok := err.(*models.ApiError)
	suite.True(ok)
	suite.Equal(http.StatusNotFound, apiErr.Code())
	suite.Equal(
		fmt.Sprintf("Feeder with ClientId %s does not exist.", clientId), apiErr.Error())
}

func (suite *FeedLogsRepositorySuite) seedFeedLogs() (feeders []dbm.Feeder, feedLogs []dbm.FeedLog) {
	count := rand.Intn(10) + 1
	for i := 0; i < count; i++ {
		feeder := modelUtils.RandomDbFeeder()
		feeders = append(feeders, feeder)
		suite.NoError(suite.r.db.Create(&feeder).Error)

		seed := modelUtils.RandomDbFeedLogsForFeeder(feeder.ClientId)
		suite.NoError(suite.r.db.Create(&seed).Error)
		feedLogs = append(feedLogs, seed...)
	}
	return feeders, feedLogs
}

func TestFeedLogsRepositorySuite(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	suite.Run(t, new(FeedLogsRepositorySuite))
}