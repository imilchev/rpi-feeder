package db

import (
	"path/filepath"

	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
)

const dbName = "db"

type DbManager interface {
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
	return &dbManager{path: dbFullPath, db: db}, nil
}

func (dbm *dbManager) Close() {
	if err := dbm.db.Close(); err != nil {
		zap.S().Errorf("Failed to close db %s. %+v", dbm.path, err)
	}
	zap.S().Infof("Database %s closed.", dbm.path)
}
