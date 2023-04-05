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
	samplePath   = "/sample/"    // folder containing the sample file for import
	arhiveName   = "cities.zip"  // name archive of cities
	jsonFilename = "cities.json" // file name containing cities
)

// CityRaw represents a struct for data from json.
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

// importCities performs transfer data of cities on database.
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
	_, err = a.Extract(rootDir + samplePath)
	if err != nil {
		return err
	}

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
	err = tx.Commit()
	if err != nil {
		return err
	}

	// remove json file
	os.Remove(rootDir + samplePath + jsonFilename)

	return nil
}

// checkDataCities performs a check exist data in table cities.
func checkDataCities(db *sqlx.DB) bool {
	var res int
	err := db.Get(&res, `SELECT COUNT(*) FROM cities`)

	if res == 0 || err != nil {
		return false
	}

	return true
}
