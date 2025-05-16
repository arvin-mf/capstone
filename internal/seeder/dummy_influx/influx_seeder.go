package main

import (
	"caps_influx/config"
	"caps_influx/internal/repository"
	"context"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	config.LoadEnv()

	client := config.InitInflux()
	defer client.Close()

	influxRepo := repository.NewInfluxRepository(
		client,
		config.GetEnv("INFLUXDB_ORG", ""),
		config.GetEnv("INFLUXDB_BUCKET", ""),
	)

	count := 1000

	for range count {
		err := influxRepo.WritePerpetual(context.Background(), repository.InfluxPerpetualPointParam{
			DeviceID:  "5",
			SubjectID: "1",
			RawEcg:    1000 * rand.Float32(),
			Timestamp: time.Now().UnixNano(),
		})
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
		time.Sleep(50 * time.Millisecond)
	}

	fmt.Println("Seeding successfully done")
}
