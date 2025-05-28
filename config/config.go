package config

import (
	"fmt"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func LoadEnv(file ...string) {
	var err error
	if len(file) > 0 && file[0] != "" {
		err = godotenv.Load(file[0])
	} else {
		err = godotenv.Load()
	}

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}

func InitDB() (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		GetEnv("DB_USER", ""),
		GetEnv("DB_PASS", ""),
		GetEnv("DB_HOST", ""),
		GetEnv("DB_PORT", ""),
		GetEnv("DB_NAME", ""),
	)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}

func InitInflux() influxdb2.Client {
	token := GetEnv("INFLUXDB_TOKEN", "")
	url := GetEnv("INFLUXDB_URL", "")
	db := influxdb2.NewClient(url, token)

	return db
}

func InitMQTTClient() mqtt.Client {
	o := mqtt.NewClientOptions().AddBroker(GetEnv("MQTT_BROKER", ""))
	o.SetClientID(GetEnv("MQTT_CLIENT", ""))
	client := mqtt.NewClient(o)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to broker: %v", token.Error())
	}

	return client
}
