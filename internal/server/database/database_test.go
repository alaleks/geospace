package database_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/alaleks/geospace/internal/server/config"
	"github.com/alaleks/geospace/internal/server/database"
	"github.com/alaleks/geospace/internal/server/database/models"
)

func BenchmarkFindCity(b *testing.B) {
	db, err := connectDB()
	if err != nil {
		b.Errorf(err.Error())
	}

	b.ResetTimer()

	b.Run("Find City by Name", func(b *testing.B) {
		_, _ = db.FindCity("Rome, It")
		_, _ = db.FindCity("Venice, It")
	})

	b.ResetTimer()

	b.Run("Find City by Alternative Name", func(b *testing.B) {
		_, _ = db.FindCity("Рим, It")
		_, _ = db.FindCity("Венеция, It")
	})

	db.Close()
}

func BenchmarkFindCityConc(b *testing.B) {
	db, err := connectDB()
	if err != nil {
		b.Errorf(err.Error())
	}

	chErr := make(chan error, 1)
	cityDepartureCh := make(chan models.City, 1)
	cityDestinationCh := make(chan models.City, 1)

	b.ResetTimer()

	b.Run("Find City by Name (Concurrency)", func(b *testing.B) {
		go db.FindCityConc("Rome, It", chErr, cityDepartureCh)
		go db.FindCityConc("Venice, It", chErr, cityDestinationCh)

		var n atomic.Int64

		for {
			if n.Load() == 2 {
				break
			}

			select {
			case <-chErr:
				break
			case <-cityDestinationCh:
				n.Add(1)
			case <-cityDepartureCh:
				n.Add(1)
			case <-time.After(500 * time.Millisecond):
				break
			}
		}
	})

	b.ResetTimer()

	b.Run("Find City by Alternative Name (Concurrency)", func(b *testing.B) {
		go db.FindCityConc("Рим, It", chErr, cityDepartureCh)
		go db.FindCityConc("Венеция, It", chErr, cityDestinationCh)

		var n atomic.Int64

		for {
			if n.Load() == 2 {
				break
			}

			select {
			case <-chErr:
				break
			case <-cityDestinationCh:
				n.Add(1)
			case <-cityDepartureCh:
				n.Add(1)
			case <-time.After(500 * time.Millisecond):
				break
			}
		}
	})

	db.Close()
}

func BenchmarkFindObjectsNearByName(b *testing.B) {
	db, err := connectDB()
	if err != nil {
		b.Errorf(err.Error())
	}

	b.ResetTimer()

	b.Run("Find Cities Nearby by Name", func(b *testing.B) {
		_, _, _ = db.FindObjectsNearByName("Rome, It", 100)
	})

	b.ResetTimer()

	b.Run("Find Cities Nearby by Alternative Name", func(b *testing.B) {
		_, _, _ = db.FindObjectsNearByName("Рим, It", 100)
	})

	db.Close()
}

func connectDB() (*database.DB, error) {
	cfg, err := config.ReadCfgFile()
	if err != nil {
		return nil, err
	}

	db, err := database.Connect(cfg)
	if err != nil {
		return nil, err
	}

	return db, nil
}
