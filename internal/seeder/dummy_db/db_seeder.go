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

	subjectRepo := repository.NewSubjectRepository(db)
	deviceRepo := repository.NewDeviceRepository(db, nil)

	subjectNames := []string{"John Doe", "M. John", "A. John"}
	deviceClientNames := []string{"aa:aa:aa:aa:aa:aa", "bb:bb:bb:bb:bb:bb", "cc:cc:cc:cc:cc:cc"}
	count := len(subjectNames) + len(deviceClientNames)

	errChan := make(chan error, count)
	var (
		wg   sync.WaitGroup
		errs []error
	)

	wg.Add(count)

	for _, n := range subjectNames {
		go func(s string) {
			defer wg.Done()

			_, err := subjectRepo.AddSubject(context.Background(), repository.Subject{
				Name: s,
			})
			if err != nil {
				errChan <- err
			}
		}(n)
	}

	for _, n := range deviceClientNames {
		go func(s string) {
			defer wg.Done()

			_, err := deviceRepo.AddDevice(context.Background(), repository.Device{
				ClientID: s,
			})
			if err != nil {
				errChan <- err
			}
		}(n)
	}

	wg.Wait()
	close(errChan)

	for e := range errChan {
		errs = append(errs, e)
	}
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Printf("Failed to seed: %v\n", e)
		}
		log.Fatal("Seeding is stopped")
	}

	fmt.Println("Seeder successfully done")
}
