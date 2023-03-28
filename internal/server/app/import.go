package app

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/alaleks/geospace/internal/server/config"
	"github.com/gen2brain/go-unarr"
	"github.com/jmoiron/sqlx"
)

const (
	samplePath   = "/sample/"
	arhiveName   = "cities.zip"
	jsonFilename = "cities.json"
)

type CityRaw struct {
	Name             string   `json:"name"`
	NameASCII        string   `json:"ascii_name"`
	AlternativeNames []string `json:"alternate_names"`
	CountryCode      string   `json:"country_code"`
	CountryName      string   `json:"label_en"`
	Timezone         string   `json:"timezone"`
	Coordinates      struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coordinates"`
}

// import performs transfer data of cities on database.
func importCities(db *sqlx.DB) error {
	// if data exists skip import of cities
	if checkDataCities(db) {
		return nil
	}

	rootDir, err := config.GetRootDir()
	if err != nil {
		return err
	}

	a, err := unarr.NewArchive(rootDir + samplePath + arhiveName)
	if err != nil {
		return err
	}

	defer a.Close()

	// extract data from archive
	a.Extract(rootDir + samplePath)

	file, err := os.Open(rootDir + samplePath + jsonFilename)
	if err != nil {
		return err
	}

	var cities []CityRaw

	err = json.NewDecoder(file).Decode(&cities)
	if err != nil {
		return err
	}

	// import data to cities table
	tx := db.MustBegin()
	for _, v := range cities {
		tx.MustExec(tx.Rebind(`INSERT INTO cities (name, name_ascii, 
			alternative_names, country_code, country, 
			timezone, latitude, longitude, created_at) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`), v.Name, v.NameASCII, strings.Join(v.AlternativeNames, ","),
			v.CountryCode, v.CountryName, v.Timezone, v.Coordinates.Lat, v.Coordinates.Lon, time.Now().Unix())
	}
	tx.Commit()

	// remove json file
	os.Remove(rootDir + samplePath + jsonFilename)

	return nil
}

// checkDataCities performs a check exist data in table cities.
func checkDataCities(db *sqlx.DB) bool {
	var res int
	db.Get(&res, `SELECT COUNT(*) FROM cities`)

	if res == 0 {
		return false
	}

	return true
}
