package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
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

func InitDB() *sqlx.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		GetEnv("DB_USER", ""),
		GetEnv("DB_PASS", ""),
		GetEnv("DB_HOST", "localhost"),
		GetEnv("DB_PORT", "3306"),
		GetEnv("DB_NAME", ""),
	)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	fmt.Println("MySQL database successfully connected")
	return db
}

func InitInflux() influxdb2.Client {
	token := GetEnv("INFLUXDB_TOKEN", "")
	url := GetEnv("INFLUXDB_URL", "http://localhost:8086")
	db := influxdb2.NewClient(url, token)

	fmt.Println("InfluxDB successfully connected")
	return db
}

func InitMQTTClient() mqtt.Client {
	o := mqtt.NewClientOptions().AddBroker(GetEnv("MQTT_BROKER", "localhost:1883"))
	o.SetClientID(GetEnv("MQTT_CLIENT", ""))
	client := mqtt.NewClient(o)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to broker: %v", token.Error())
	}

	fmt.Println("Successfully connected to MQTT broker")
	return client
}

func InitRedis() *redis.Client {
	dbIndex, err := strconv.Atoi(GetEnv("REDIS_DB", "0"))
	if err != nil {
		log.Fatalf("Failed to get and convert Redis DB index from .env: %v", err)
	}

	rd := redis.NewClient(&redis.Options{
		Addr:     GetEnv("REDIS_HOST", "localhost:6379"),
		Password: GetEnv("REDIS_PASS", ""),
		DB:       dbIndex,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rd.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	fmt.Println("Redis successfully connected")
	return rd
}
