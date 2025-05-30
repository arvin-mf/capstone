package repository

import (
	"caps_influx/config"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

var (
	testDB      *sqlx.DB
	redisClient *redis.Client
	deviceRepo  DeviceRepository
	subjectRepo SubjectRepository
)

func TestMain(m *testing.M) {
	config.LoadEnv("../../.env.test")

	testDB = config.InitDB()
	defer testDB.Close()

	redisClient = config.InitRedis()
	defer redisClient.Close()

	deviceRepo = NewDeviceRepository(testDB, redisClient)
	subjectRepo = NewSubjectRepository(testDB)

	os.Exit(m.Run())
}
