package main

import (
	"caps_influx/config"
	"caps_influx/internal/repository"
	"context"
	"fmt"
	"log"
	"sync"
)

func main() {
	config.LoadEnv()

	db := config.InitDB()
	defer db.Close()

	deviceRepo := repository.NewDeviceRepository(db, nil)

	deviceClients := []string{"c8:2e:18:26:65:90"}

	count := len(deviceClients)
	errChan := make(chan error, count)
	var (
		wg   sync.WaitGroup
		errs []error
	)
	wg.Add(count)

	for _, d := range deviceClients {
		go func(s string) {
			defer wg.Done()

			_, err := deviceRepo.AddDevice(context.Background(), repository.Device{
				ClientID: d,
			})
			if err != nil {
				errChan <- err
			}
		}(d)
	}

	wg.Wait()
	close(errChan)

	for e := range errChan {
		errs = append(errs, e)
	}
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Printf("Failed to seed: %v", e)
		}
		log.Fatal("Seeding is stopped")
	}

	fmt.Println("Database successfully seeded")
}
