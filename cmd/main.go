package main

import (
	"caps_influx/config"
	"caps_influx/internal/server"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	db := config.InitDB()
	defer db.Close()

	influxClient := config.InitInflux()
	defer influxClient.Close()

	mqttClient := config.InitMQTTClient()
	defer mqttClient.Disconnect(250)

	redisClient := config.InitRedis()
	defer redisClient.Close()

	eng := gin.Default()

	server.StartEngine(eng, db, influxClient, mqttClient, redisClient)

	port := config.GetEnv("APP_PORT", "")
	eng.Run(":" + port)
}
