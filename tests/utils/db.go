package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/imilchev/rpi-feeder/pkg/service/db/migrations"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	dbPortEnvVar  = "FEEDER_TEST_DB_PORT"
	dbDefaultPort = 5432
	dbHostEnvVar  = "FEEDER_TEST_DB_HOST"
	dbDefaultHost = "localhost"
)

var isInitialized = false

func GetDbConnString() string {
	port := dbDefaultPort
	p := os.Getenv(dbPortEnvVar)
	if p != "" {
		pInt, err := strconv.ParseInt(p, 10, 32)
		if err == nil {
			port = int(pInt)
		}
	}

	host := dbDefaultHost
	h := os.Getenv(dbHostEnvVar)
	if h != "" {
		host = h
	}
	return fmt.Sprintf(
		"postgres://postgres:SuperSecret@%s:%d/postgres?sslmode=disable", host, port)
}

func InitTestDb() error {
	if isInitialized {
		return nil
	}

	d, err := iofs.New(migrations.FS, ".")
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, GetDbConnString())
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		// The framework return "no change" if the db is up-to-date. In such
		// case we do not want to return an error as we run Up on every start
		// of GoClapy.
		if err.Error() == "no change" {
			isInitialized = true
			return nil
		}
		return err
	}
	isInitialized = true
	return nil
}

func DropAllTables() error {
	d, err := iofs.New(migrations.FS, ".")
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, GetDbConnString())
	if err != nil {
		return err
	}
	if err := m.Down(); err != nil {
		return err
	}
	isInitialized = false
	return nil
}

func GetTestDb() (*gorm.DB, error) {
	// logger := zapgorm2.New(zap.L())
	// logger.IgnoreRecordNotFoundError = false
	// logger.SetAsDefault()

	return gorm.Open(postgres.Open(GetDbConnString())) //, &gorm.Config{Logger: logger})
}

func CleanupDb(db *gorm.DB) error {
	return db.Exec(`TRUNCATE TABLE "feed_logs" CASCADE;
					TRUNCATE TABLE "feeders" CASCADE;`).Error
}

func GetMigrationsCount() (uint, error) {
	// For each migration there are 2 SQL files - one for up and one for down. The
	// amount of migrations we have is equal to half of the SQL files.
	entries, err := migrations.FS.ReadDir(".")
	return uint(len(entries) / 2), err
}
