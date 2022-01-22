package repos

import (
	"fmt"

	"github.com/imilchev/rpi-feeder/pkg/service/models"
)

// FakeFeedersRepository provides an easy way of mocking a FeedersRepository.
// The functions in this fake implementation do not perform any validation.
type FakeFeedersRepository struct {
	Feeders []models.Feeder

	// Error If this is set, any function will return it.
	Error error
}

func (r *FakeFeedersRepository) CreateFeeder(f models.Feeder) (models.Feeder, error) {
	if r.Error != nil {
		return models.Feeder{}, r.Error
	}

	r.Feeders = append(r.Feeders, f)
	return f, nil
}

func (r *FakeFeedersRepository) GetFeeders() (f []models.Feeder, err error) {
	if r.Error != nil {
		return f, r.Error
	}

	f = append(f, r.Feeders...)
	return f, nil
}

func (r *FakeFeedersRepository) GetFeederByClientId(cId string) (models.Feeder, error) {
	if r.Error != nil {
		return models.Feeder{}, r.Error
	}

	for _, f := range r.Feeders {
		if f.ClientId == cId {
			return f, nil
		}
	}
	return models.Feeder{}, fmt.Errorf("not found")
}

func (r *FakeFeedersRepository) UpdateFeeder(f models.Feeder) (models.Feeder, error) {
	if r.Error != nil {
		return models.Feeder{}, r.Error
	}

	for i, ff := range r.Feeders {
		if ff.ClientId == f.ClientId {
			r.Feeders[i] = f
			return f, nil
		}
	}
	return models.Feeder{}, fmt.Errorf("not found")
}
