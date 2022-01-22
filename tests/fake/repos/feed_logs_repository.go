package repos

import (
	"github.com/imilchev/rpi-feeder/pkg/service/models"
)

// FakeFeedLogsRepository provides an easy way of mocking a FeedLogsRepository.
// The functions in this fake implementation do not perform any validation.
type FakeFeedLogsRepository struct {
	FeedLogs []models.FeedLog

	// Error If this is set, any function will return it.
	Error error
}

func (r *FakeFeedLogsRepository) CreateFeedLogs(f []models.FeedLog) (fs []models.FeedLog, err error) {
	if r.Error != nil {
		return fs, r.Error
	}

	r.FeedLogs = append(r.FeedLogs, f...)
	return f, nil
}

func (r *FakeFeedLogsRepository) GetLogsForFeeder(clientId string) (f []models.FeedLog, err error) {
	if r.Error != nil {
		return f, r.Error
	}

	f = append(f, r.FeedLogs...)
	return f, nil
}
