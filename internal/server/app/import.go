package app

import (
	"archive/zip"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/alaleks/geospace/internal/server/config"
	"github.com/alaleks/geospace/internal/server/database/models"
	"github.com/jmoiron/sqlx"
)

const (
	samplePath   = "/sample/"    // folder containing the sample file for import
	arhiveName   = "cities.zip"  // name archive of cities
	jsonFilename = "cities.json" // file name containing cities
	numCities    = 140_868
)

// CityRaw represents a struct for data from json.
type CityRaw struct {
	Name             string   `json:"name"`
	NameASCII        string   `json:"ascii_name"`
	CountryCode      string   `json:"country_code"`
	CountryName      string   `json:"label_en"`
	Timezone         string   `json:"timezone"`
	AlternativeNames []string `json:"alternate_names"`
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

	// open zip file
	archiveCities, err := zip.OpenReader(rootDir + samplePath + arhiveName)
	if err != nil {
		return err
	}

	defer archiveCities.Close()

	var fileCities []byte

	// read files in zip archive
	for _, f := range archiveCities.File {
		v, err := f.Open()
		if err != nil {
			return err
		}

		defer v.Close()

		b, err := io.ReadAll(v)
		if err != nil {
			return err
		}

		fileCities = b

		break // import file is one
	}

	citiesRaw := make([]CityRaw, 0, numCities)

	err = json.Unmarshal(fileCities, &citiesRaw)
	if err != nil {
		return err
	}

	cities := make([]models.City, 0, numCities)

	for _, city := range citiesRaw {
		cities = append(cities, models.City{
			Name:             city.Name,
			NameASCII:        city.NameASCII,
			AlternativeNames: strings.Join(city.AlternativeNames, ","),
			CountryCode:      city.CountryCode,
			Country:          city.CountryName,
			Timezone:         city.Timezone,
			Latitude:         city.Coordinates.Lat,
			Longitude:        city.Coordinates.Lon,
			CreatedAt:        time.Now().Unix(),
		})
	}

	// import data to cities table
	if _, err = db.NamedExec(`INSERT INTO cities (name, name_ascii, 
		alternative_names, country_code, country, 
		timezone, latitude, longitude, created_at) 
        VALUES (:name, :name_ascii, :alternative_names, :country_code, :country,
			 :country_code, :timezone, :latitude, :longitude, :created_at)`, cities); err != nil {
		return err
	}

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
