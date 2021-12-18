package db

import (
	"encoding/json"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/imilchev/rpi-feeder/pkg/db/model"
	"github.com/stretchr/testify/suite"
	bolt "go.etcd.io/bbolt"
)

const dbPath = "./"

type DbManagerSuite struct {
	suite.Suite
	db *dbManager
}

func (suite *DbManagerSuite) SetupTest() {
	m, err := NewDbManager(dbPath)
	suite.NoError(err)
	suite.db = m.(*dbManager)
}

func (suite *DbManagerSuite) AfterTest(suiteName, testName string) {
	suite.db.Close()
	suite.NoError(os.Remove(filepath.Join(dbPath, dbName)))
}

func (suite *DbManagerSuite) TestAddFeedLog() {
	testLog := model.FeedLog{
		Portions:  uint(rand.Intn(100)),
		Timestamp: time.Now(),
	}
	suite.NoError(suite.db.AddFeedLog(testLog))

	var logs []model.FeedLog
	suite.NoError(suite.db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(logBucketName))
		return b.ForEach(func(k, v []byte) error {
			l := &model.FeedLog{}
			suite.NoError(json.Unmarshal(v, l))
			logs = append(logs, *l)
			return nil
		})
	}))

	suite.Equal(1, len(logs))
	suite.Equal(testLog.Portions, logs[0].Portions)
	suite.Equal(testLog.Timestamp.UTC(), logs[0].Timestamp)
}

func (suite *DbManagerSuite) TestListFeedLog() {
	logsCount := rand.Intn(100)
	testLogs := make([]model.FeedLog, logsCount)

	for i := 0; i < logsCount; i++ {
		testLogs[i] = model.FeedLog{
			Portions:  uint(rand.Intn(100)),
			Timestamp: time.Now().UTC(),
		}
	}

	suite.NoError(suite.db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(logBucketName))

		for i := 0; i < logsCount; i++ {
			id, err := b.NextSequence()
			suite.NoError(err)
			testLogs[i].Id = int(id)
			data, err := json.Marshal(testLogs[i])
			suite.NoError(err)
			suite.NoError(b.Put(itob(testLogs[i].Id), data))
		}
		return nil
	}))

	logs, err := suite.db.ListFeedLog()
	suite.NoError(err)

	suite.Equal(logsCount, len(logs))
	suite.EqualValues(testLogs, logs)
}

func TestDbManagerSuite(t *testing.T) {
	suite.Run(t, new(DbManagerSuite))
}
