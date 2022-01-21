package db

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/imilchev/rpi-feeder/pkg/service/config"
	"github.com/imilchev/rpi-feeder/pkg/service/db/migrations"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

type Database struct {
	config config.Database
	DB     *gorm.DB
}

func NewDatabaseConnection(cfg config.Database) (*Database, error) {
	logger := zapgorm2.New(zap.L())
	logger.SetAsDefault() // optional: configure gorm to use this zapgorm.Logger for callbacks

	db, err := gorm.Open(postgres.Open(cfg.ConnectionString), &gorm.Config{Logger: logger})
	if err != nil {
		return nil, err
	}
	d := &Database{config: cfg, DB: db}
	return d, d.migrateDatabase()
}

func (d *Database) Close() error {
	dbConn, err := d.DB.DB()
	if err != nil {
		return err
	}

	if err := dbConn.Close(); err != nil {
		zap.S().Errorf("Failed to close database. %v", err)
		return err
	}
	zap.S().Info("Successfully closed database connections.")
	return nil
}

func (d *Database) migrateDatabase() error {
	zap.S().Debugf("Migrating the database...")
	mFs, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("iofs", mFs, d.config.ConnectionString)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		// The framework return "no change" if the db is up-to-date. In such
		// case we do not want to return an error as we run Up on every start
		// of GoClapy.
		if err.Error() == "no change" {
			zap.S().Info("Database is already up-to-date.")
			return nil
		}

		return err
	}

	zap.S().Info("Successfully migrated database.")
	return nil
}
