package repository

import (
	"caps_influx/config"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
)

var (
	testDB      *sqlx.DB
	deviceRepo  DeviceRepository
	subjectRepo SubjectRepository
)

func TestMain(m *testing.M) {
	config.LoadEnv("../../.env.test")

	var err error
	testDB, err = config.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}
	defer testDB.Close()

	deviceRepo = NewDeviceRepository(testDB)
	subjectRepo = NewSubjectRepository(testDB)

	os.Exit(m.Run())
}
