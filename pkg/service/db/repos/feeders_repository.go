package repos

import (
	dbm "github.com/imilchev/rpi-feeder/pkg/service/db/models"
	"github.com/imilchev/rpi-feeder/pkg/service/models"
	"github.com/imilchev/rpi-feeder/pkg/utils"
	"gorm.io/gorm"
)

type FeedersRepository interface {
	CreateFeeder(u models.Feeder) (models.Feeder, error)
	GetFeeders() ([]models.Feeder, error)
	GetFeederByClientId(cId string) (models.Feeder, error)
	UpdateFeeder(f models.Feeder) (models.Feeder, error)
}

type feedersRepository struct {
	db *gorm.DB
}

func NewFeedersRepository(db *gorm.DB) FeedersRepository {
	return &feedersRepository{db: db}
}

func (r *feedersRepository) CreateFeeder(f models.Feeder) (models.Feeder, error) {
	if err := utils.Validate.Struct(f); err != nil {
		return models.Feeder{}, models.NewValidationError(err.Error())
	}

	if res := r.db.Where("client_id = ?", f.ClientId).Find(&dbm.Feeder{}); res.RowsAffected > 0 {
		return models.Feeder{},
			models.NewAlreadyExistsError("Feeder", "ClientId", f.ClientId)
	}

	dbModel := &dbm.Feeder{}
	dbModel.FromApi(f)
	if res := r.db.Create(dbModel); res.Error != nil {
		return models.Feeder{}, res.Error
	}
	createdFeeder := models.Feeder{}
	dbModel.ToApi(&createdFeeder)
	return createdFeeder, nil
}

func (r *feedersRepository) GetFeeders() (f []models.Feeder, err error) {
	var feeders []dbm.Feeder
	if res := r.db.Find(&feeders); res.Error != nil {
		return f, res.Error
	}

	apiFeeder := &models.Feeder{}
	for _, c := range feeders {
		c.ToApi(apiFeeder)
		f = append(f, *apiFeeder)
	}
	return f, nil
}

func (r *feedersRepository) GetFeederByClientId(cId string) (models.Feeder, error) {
	c := dbm.Feeder{}
	if res := r.db.Where("client_id = ?", cId).Find(&c); res.RowsAffected == 0 {
		return models.Feeder{}, models.NewDoesNotExistError("Feeder", "ClientId", cId)
	}

	cApi := models.Feeder{}
	c.ToApi(&cApi)

	return cApi, nil
}

func (r *feedersRepository) UpdateFeeder(f models.Feeder) (models.Feeder, error) {
	dbModel := &dbm.Feeder{}
	if res := r.db.Where("client_id = ?", f.ClientId).Find(dbModel); res.RowsAffected == 0 {
		return models.Feeder{}, models.NewDoesNotExistError("Feeder", "ClientId", f.ClientId)
	}

	dbModel.FromApi(f)
	if res := r.db.Model(dbModel).Where("client_id = ?", f.ClientId).
		Select("status", "last_online", "software_version").
		Updates(dbModel); res.Error != nil {
		return models.Feeder{}, res.Error
	}
	return f, nil
}
