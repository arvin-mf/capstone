package main

import (
	"caps_influx/config"
	"caps_influx/internal/server"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	db, err := config.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	influxClient := config.InitInflux()
	defer influxClient.Close()

	mqttClient := config.InitMQTTClient()
	defer mqttClient.Disconnect(250)

	eng := gin.Default()

	server.StartEngine(eng, db, influxClient, mqttClient)

	port := config.GetEnv("APP_PORT", "")
	eng.Run(":" + port)
}
