package repos

import (
	dbm "github.com/imilchev/rpi-feeder/pkg/service/db/models"
	"github.com/imilchev/rpi-feeder/pkg/service/models"
	"github.com/imilchev/rpi-feeder/pkg/utils"
	"gorm.io/gorm"
)

type FeedLogsRepository interface {
	CreateFeedLogs(f []models.FeedLog) ([]models.FeedLog, error)
	GetLogsForFeeder(clientId string) ([]models.FeedLog, error)
}

type feedLogsRepository struct {
	db *gorm.DB
}

func NewFeedLogsRepository(db *gorm.DB) FeedLogsRepository {
	return &feedLogsRepository{db: db}
}

func (r *feedLogsRepository) CreateFeedLogs(f []models.FeedLog) ([]models.FeedLog, error) {
	for _, v := range f {
		if err := utils.Validate.Struct(v); err != nil {
			return []models.FeedLog{}, models.NewValidationError(err.Error())
		}
	}

	var dbModels []dbm.FeedLog
	for _, m := range f {
		dbm := &dbm.FeedLog{}
		dbm.FromApi(m)
		dbModels = append(dbModels, *dbm)
	}

	if res := r.db.CreateInBatches(dbModels, 100); res.Error != nil {
		return []models.FeedLog{}, res.Error
	}

	f = make([]models.FeedLog, 0)
	for _, m := range dbModels {
		createdFeedLog := models.FeedLog{}
		m.ToApi(&createdFeedLog)
		f = append(f, createdFeedLog)
	}
	return f, nil
}

func (r *feedLogsRepository) GetLogsForFeeder(clientId string) (f []models.FeedLog, err error) {
	var feedLogs []dbm.FeedLog
	if res := r.db.Where("client_id", clientId).Find(&feedLogs); res.Error != nil {
		return f, res.Error
	}

	apiFeedLog := &models.FeedLog{}
	for _, c := range feedLogs {
		c.ToApi(apiFeedLog)
		f = append(f, *apiFeedLog)
	}
	return f, nil
}
