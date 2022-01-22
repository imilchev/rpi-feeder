package models

type FeedLog struct {
	Id        int
	ClientId  string `validate:"required,max=60"`
	Portions  uint   `validate:"required,gt=0"`
	Timestamp int64  `validate:"required"`
}
