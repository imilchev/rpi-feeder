package repos

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/imilchev/rpi-feeder/pkg/mqtt/model"
	dbm "github.com/imilchev/rpi-feeder/pkg/service/db/models"
	"github.com/imilchev/rpi-feeder/pkg/service/models"
	"github.com/imilchev/rpi-feeder/tests/utils"
	modelUtils "github.com/imilchev/rpi-feeder/tests/utils/models"
	"github.com/stretchr/testify/suite"
)

type FeedersRepositorySuite struct {
	suite.Suite
	r *feedersRepository
}

func (suite *FeedersRepositorySuite) SetupTest() {
	suite.Require().NoError(utils.InitTestDb())
	db, err := utils.GetTestDb()
	suite.Require().NoError(err)
	suite.r = &feedersRepository{db: db}
}

func (suite *FeedersRepositorySuite) AfterTest(suiteName, testName string) {
	suite.Require().NoError(utils.CleanupDb(suite.r.db))
	db, err := suite.r.db.DB()
	suite.Require().NoError(err)
	db.Close()
}

func (suite *FeedersRepositorySuite) TestCreateFeeder() {
	f := modelUtils.RandomFeeder()
	ff, err := suite.r.CreateFeeder(f)
	suite.NoError(err)
	suite.Equal(f, ff)

	fDb := &dbm.Feeder{}
	suite.NoError(suite.r.db.First(fDb, "client_id = ?", f.ClientId).Error)
	fDb.ToApi(&ff)

	if f.LastOnline != nil {
		suite.Equal(f.LastOnline.Unix(), ff.LastOnline.Unix())
		f.LastOnline = nil
		ff.LastOnline = nil
	}

	suite.Equal(f, ff)
}

func (suite *FeedersRepositorySuite) TestCreateFeeder_AlreadyExists() {
	f := modelUtils.RandomDbFeeder()
	suite.NoError(suite.r.db.Create(&f).Error)

	f2 := models.Feeder{}
	f.ToApi(&f2)
	_, err := suite.r.CreateFeeder(f2)
	suite.Error(err)
	apiErr, ok := err.(*models.ApiError)
	suite.True(ok)
	suite.Equal(http.StatusConflict, apiErr.Code())
	suite.Equal(
		fmt.Sprintf("Feeder with ClientId %s already exists.", f.ClientId), apiErr.Error())
}

func (suite *FeedersRepositorySuite) TestCreateFeeder_ClientIdMissing() {
	f := modelUtils.RandomFeeder()
	f.ClientId = ""
	_, err := suite.r.CreateFeeder(f)
	suite.Error(err)
	apiErr, ok := err.(*models.ApiError)
	suite.True(ok)
	suite.Equal(http.StatusBadRequest, apiErr.Code())
}

func (suite *FeedersRepositorySuite) TestCreateFeeder_ClientIdTooLong() {
	f := modelUtils.RandomFeeder()
	f.ClientId = utils.RandString(61)
	_, err := suite.r.CreateFeeder(f)
	suite.Error(err)
	apiErr, ok := err.(*models.ApiError)
	suite.True(ok)
	suite.Equal(http.StatusBadRequest, apiErr.Code())
}

func (suite *FeedersRepositorySuite) TestCreateFeeder_SoftwareVersionMissing() {
	f := modelUtils.RandomFeeder()
	f.SoftwareVersion = ""
	_, err := suite.r.CreateFeeder(f)
	suite.Error(err)
	apiErr, ok := err.(*models.ApiError)
	suite.True(ok)
	suite.Equal(http.StatusBadRequest, apiErr.Code())
}

func (suite *FeedersRepositorySuite) TestCreateFeeder_SoftwareVersionTooLong() {
	f := modelUtils.RandomFeeder()
	f.SoftwareVersion = utils.RandString(61)
	_, err := suite.r.CreateFeeder(f)
	suite.Error(err)
	apiErr, ok := err.(*models.ApiError)
	suite.True(ok)
	suite.Equal(http.StatusBadRequest, apiErr.Code())
}

func (suite *FeedersRepositorySuite) TestCreateFeeder_StatusMissing() {
	f := modelUtils.RandomFeeder()
	f.Status = ""
	_, err := suite.r.CreateFeeder(f)
	suite.Error(err)
	apiErr, ok := err.(*models.ApiError)
	suite.True(ok)
	suite.Equal(http.StatusBadRequest, apiErr.Code())
}

func (suite *FeedersRepositorySuite) TestCreateFeeder_StatusTooLong() {
	f := modelUtils.RandomFeeder()
	f.Status = model.Status(utils.RandString(8))
	_, err := suite.r.CreateFeeder(f)
	suite.Error(err)
	apiErr, ok := err.(*models.ApiError)
	suite.True(ok)
	suite.Equal(http.StatusBadRequest, apiErr.Code())
}

func (suite *FeedersRepositorySuite) TestGetFeeders() {
	seed := suite.seedFeeders()
	var expected []models.Feeder
	apiFeeder := &models.Feeder{}
	for _, f := range seed {
		f.ToApi(apiFeeder)
		expected = append(expected, *apiFeeder)
	}

	feeders, err := suite.r.GetFeeders()
	suite.NoError(err)

	// Compare last online in unix as the nanosecs precision is lost when stored
	// in the db.
	for i, f := range expected {
		if f.LastOnline != nil {
			suite.Equal(f.LastOnline.Unix(), feeders[i].LastOnline.Unix())
			expected[i].LastOnline = nil
			feeders[i].LastOnline = nil
		}
	}
	// Compare the rest of the properties.
	suite.ElementsMatch(expected, feeders)
}

func (suite *FeedersRepositorySuite) TestGetFeeders_NoFeeders() {
	feeders, err := suite.r.GetFeeders()
	suite.NoError(err)
	suite.Equal(0, len(feeders))
}

func (suite *FeedersRepositorySuite) TestGetFeederByClientId() {
	seed := suite.seedFeeders()
	randFeeder := models.Feeder{}
	seed[len(seed)/2].ToApi(&randFeeder)

	feeder, err := suite.r.GetFeederByClientId(randFeeder.ClientId)
	suite.NoError(err)
	suite.Equal(randFeeder, feeder)
}

func (suite *FeedersRepositorySuite) TestGetFeederByClientId_DoesNotExist() {
	suite.seedFeeders()
	cId := "does-not-exist"
	_, err := suite.r.GetFeederByClientId(cId)
	suite.Error(err)
	apiErr, ok := err.(*models.ApiError)
	suite.True(ok)
	suite.Equal(http.StatusNotFound, apiErr.Code())
	suite.Equal(
		fmt.Sprintf("Feeder with ClientId %s does not exist.", cId), apiErr.Error())
}

func (suite *FeedersRepositorySuite) TestUpdateFeeder_OnlineToOffline() {
	f := modelUtils.RandomFeeder()
	f.Status = model.OnlineStatus
	f.LastOnline = nil
	suite.NoError(suite.r.db.Create(f).Error)

	t := time.Now().UTC()
	f.LastOnline = &t
	f.Status = model.OfflineStatus

	ff, err := suite.r.UpdateFeeder(f)
	suite.NoError(err)
	suite.Equal(f, ff)

	fDb := &dbm.Feeder{}
	suite.NoError(suite.r.db.First(fDb, "client_id = ?", f.ClientId).Error)
	fDb.ToApi(&ff)

	suite.Equal(f.LastOnline.Unix(), ff.LastOnline.Unix())
	f.LastOnline = nil
	ff.LastOnline = nil

	suite.Equal(f, ff)
}

func (suite *FeedersRepositorySuite) TestUpdateFeeder_OfflineToOnline() {
	t := time.Now().UTC()
	f := modelUtils.RandomFeeder()
	f.Status = model.OfflineStatus
	f.LastOnline = &t
	suite.NoError(suite.r.db.Create(f).Error)

	f.LastOnline = nil
	f.Status = model.OnlineStatus

	ff, err := suite.r.UpdateFeeder(f)
	suite.NoError(err)
	suite.Equal(f, ff)

	fDb := &dbm.Feeder{}
	suite.NoError(suite.r.db.First(fDb, "client_id = ?", f.ClientId).Error)
	fDb.ToApi(&ff)
	suite.Equal(f, ff)
}

func (suite *FeedersRepositorySuite) TestUpdateFeeder_DoesNotExist() {
	f := modelUtils.RandomFeeder()

	_, err := suite.r.UpdateFeeder(f)
	suite.Error(err)
	apiErr, ok := err.(*models.ApiError)
	suite.True(ok)
	suite.Equal(http.StatusNotFound, apiErr.Code())
	suite.Equal(
		fmt.Sprintf("Feeder with ClientId %s does not exist.", f.ClientId), apiErr.Error())
}

func (suite *FeedersRepositorySuite) seedFeeders() (feeders []dbm.Feeder) {
	count := rand.Intn(20) + 1
	for i := 0; i < count; i++ {
		f := modelUtils.RandomDbFeeder()
		feeders = append(feeders, f)
	}
	suite.Require().NoError(suite.r.db.Create(&feeders).Error)
	return feeders
}

func TestFeedersRepositorySuite(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	suite.Run(t, new(FeedersRepositorySuite))
}
