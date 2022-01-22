package db

import (
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/imilchev/rpi-feeder/pkg/service/config"
	"github.com/imilchev/rpi-feeder/pkg/service/db/migrations"
	"github.com/imilchev/rpi-feeder/tests/utils"
	"github.com/stretchr/testify/suite"
)

type DatabaseSuite struct {
	suite.Suite
	cfg config.Database
}

func (suite *DatabaseSuite) SetupTest() {
	suite.cfg = config.Database{
		ConnectionString: utils.GetDbConnString(),
	}
}

func (suite *DatabaseSuite) TestInit_Migrate() {
	// Drop all tables first.
	suite.Require().NoError(utils.DropAllTables())

	mFs, err := iofs.New(migrations.FS, ".")
	suite.NoError(err)
	m, err := migrate.NewWithSourceInstance("iofs", mFs, suite.cfg.ConnectionString)
	suite.NoError(err)

	// Make sure the database is empty.
	_, _, err = m.Version()
	suite.Error(err) // migrate returns an error if no migrations have been executed

	migrations, err := utils.GetMigrationsCount()
	suite.NoError(err)

	db, err := NewDatabaseConnection(suite.cfg)
	suite.NoError(err)

	// Verify the migrations have been executed.
	version, dirty, err := m.Version()
	suite.NoError(err)
	suite.False(dirty)
	suite.Equal(migrations, version)

	suite.Require().NoError(db.Close())
}

func (suite *DatabaseSuite) TestInit_AlreadyMigrated() {
	// Make sure the DB is migrated.
	suite.Require().NoError(utils.InitTestDb())

	mFs, err := iofs.New(migrations.FS, ".")
	suite.NoError(err)
	m, err := migrate.NewWithSourceInstance("iofs", mFs, suite.cfg.ConnectionString)
	suite.NoError(err)

	migrations, err := utils.GetMigrationsCount()
	suite.NoError(err)

	// Make sure the database is up-to-date.
	version, dirty, err := m.Version()
	suite.NoError(err)
	suite.False(dirty)
	suite.Equal(migrations, version)

	db, err := NewDatabaseConnection(suite.cfg)
	suite.NoError(err)

	// Verify the migrations have been executed.
	version, dirty, err = m.Version()
	suite.NoError(err)
	suite.False(dirty)
	suite.Equal(migrations, version)

	suite.Require().NoError(db.Close())
}

func TestDatabaseSuiteSuite(t *testing.T) {
	suite.Run(t, new(DatabaseSuite))
}
