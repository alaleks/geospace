package database_test

import (
	"fmt"
	"testing"

	"github.com/alaleks/geospace/internal/server/config"
	"github.com/alaleks/geospace/internal/server/database"
)

func BenchmarkFindCity(b *testing.B) {
	cfg, err := config.ReadCfgFile()
	if err != nil {
		b.Errorf(err.Error())
	}

	db, err := database.Connect(cfg)
	if err != nil {
		b.Errorf(err.Error())
	}

	b.ResetTimer()

	b.Run(fmt.Sprintf("Find City by Name"), func(b *testing.B) {
		_, _ = db.FindCity("Rome, It")
	})

	b.ResetTimer()

	b.Run(fmt.Sprintf("Find City by Alternative Name"), func(b *testing.B) {
		_, _ = db.FindCity("Рим, It")
	})

	db.Close()
}

func BenchmarkFindObjectsNearByName(b *testing.B) {
	cfg, err := config.ReadCfgFile()
	if err != nil {
		b.Errorf(err.Error())
	}

	db, err := database.Connect(cfg)
	if err != nil {
		b.Errorf(err.Error())
	}

	b.ResetTimer()

	b.Run(fmt.Sprintf("Find Cities Nearby by Name"), func(b *testing.B) {
		_, _, _ = db.FindObjectsNearByName("Rome, It", 100)
	})

	b.ResetTimer()

	b.Run(fmt.Sprintf("Find Cities Nearby by Alternative Name"), func(b *testing.B) {
		_, _, _ = db.FindObjectsNearByName("Рим, It", 100)
	})

	db.Close()
}
