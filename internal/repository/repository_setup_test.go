package repository

import (
	"caps_influx/config"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
)

var testDB *sqlx.DB
var deviceRepo DeviceRepository

func TestMain(m *testing.M) {
	config.LoadEnv("../../.env.test")

	var err error
	testDB, err = config.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}
	defer testDB.Close()

	testDB.Exec("TRUNCATE TABLE devices")

	deviceRepo = NewDeviceRepository(testDB)

	os.Exit(m.Run())
}
