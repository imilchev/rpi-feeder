package db

import (
	"encoding/json"
	"path/filepath"

	"github.com/imilchev/rpi-feeder/pkg/db/model"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
)

const (
	dbName = "feeder.db"
)

var (
	logBucketName = []byte("feeder-log")
)

func initBuckets(db *bolt.DB) error {
	zap.S().Debug("Initializing buckets...")
	return db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(logBucketName); err != nil {
			return err
		}
		zap.S().Debugf("Initialized bucket %s.", logBucketName)
		return nil
	})
}

type DbManager interface {
	AddFeedLog(model.FeedLog) error
	ListFeedLog() ([]model.FeedLog, error)
	CleanFeedLog() error
	Close()
}

type dbManager struct {
	path string
	db   *bolt.DB
}

func NewDbManager(path string) (DbManager, error) {
	dbFullPath := filepath.Join(path, dbName)
	db, err := bolt.Open(dbFullPath, 0666, nil)
	if err != nil {
		zap.S().Errorf("Failed to initialize database. %+v", err)
		return nil, err
	}
	if err := initBuckets(db); err != nil {
		return nil, err
	}
	return &dbManager{path: dbFullPath, db: db}, nil
}

func (m *dbManager) AddFeedLog(log model.FeedLog) error {
	return m.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(logBucketName)
		id, err := bucket.NextSequence()
		log.Id = int(id)
		if err != nil {
			return err
		}
		data, err := json.Marshal(log)
		if err != nil {
			return err
		}
		if err := bucket.Put(itob(log.Id), data); err != nil {
			return err
		}
		zap.S().Debugf("Written feed log %+v.", log)
		return nil
	})
}

func (m *dbManager) ListFeedLog() ([]model.FeedLog, error) {
	var logs []model.FeedLog
	err := m.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(logBucketName)

		return bucket.ForEach(func(k, v []byte) error {
			l := &model.FeedLog{}
			if err := json.Unmarshal(v, l); err != nil {
				return err
			}
			logs = append(logs, *l)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (m *dbManager) CleanFeedLog() error {
	err := m.db.Update(func(tx *bolt.Tx) error {
		if err := tx.DeleteBucket(logBucketName); err != nil {
			return err
		}

		_, err := tx.CreateBucket(logBucketName)
		return err
	})
	return err
}

func (m *dbManager) Close() {
	if err := m.db.Close(); err != nil {
		zap.S().Errorf("Failed to close db %s. %+v", m.path, err)
	}
	zap.S().Infof("Database %s closed.", m.path)
}
