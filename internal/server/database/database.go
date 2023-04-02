// Package database performs operations over the database.
package database

import (
	"errors"
	"math"
	"strings"
	"time"

	"github.com/alaleks/geospace/internal/server/config"
	"github.com/alaleks/geospace/internal/server/database/models"
	"github.com/alaleks/geospace/internal/server/database/schema"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	MaxIdleConns    = 100              // maximum number of concurrent connections to the database
	ConnMaxLifetime = 15 * time.Minute // the maximum length of time a connection can be reused
	// table names
	tableCities       = "cities"
	tableUsers        = "users"
	oneDegreesInKmLat = 110.574   // km in one degree latitude
	oneDegreesInKmLon = 111.320   // km in one degree longitude
	converFact        = 100000000 // number for convert floating point to uint
)

// typical errors
var (
	ErrUserAlreadyExists = errors.New("user with current email already exists")
)

// DB contains pointer to SQLX instance.
type DB struct {
	SQLX *sqlx.DB
}

// Connect performs creating a new connection to database.
func Connect(cfg config.Cfg) (*DB, error) {
	db, err := sqlx.Connect("mysql", cfg.CreateDSN())
	if err != nil {
		return nil, err
	}

	// set params of connection
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetConnMaxLifetime(ConnMaxLifetime)

	return &DB{
		SQLX: db,
	}, nil
}

// Migrate performs a migration schema of table to database.
func (db *DB) Migrate() {
	if !db.checkTableExist(tableCities) {
		db.SQLX.MustExec(schema.City)
	}

	if !db.checkTableExist(tableUsers) {
		db.SQLX.MustExec(schema.User)
	}
}

// Close perfoms closing the database connection.
func (db *DB) Close() error {
	return db.SQLX.Close()
}

// CreateUser performs a create user to database.
func (db *DB) CreateUser(name, email, password string) (int, error) {
	var res int
	err := db.SQLX.Get(&res, `SELECT COUNT(*) FROM users
	WHERE email = ?`, email)
	if err != nil {
		return 0, err
	}

	if res > 0 {
		return 0, ErrUserAlreadyExists
	}

	user := models.User{
		Name:      name,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now().Unix(),
	}

	_, err = db.SQLX.NamedExec(`INSERT INTO users (name, email, password, created_at) 
	VALUES (:name, :email, :password, :created_at)`,
		&user)
	if err != nil {
		return 0, err
	}

	return user.UID, err
}

// GetUser provides a get user from database by email.
func (db *DB) GetUser(email string) (models.User, error) {
	var user models.User
	err := db.SQLX.Get(&user, "SELECT * FROM users WHERE email=?", email)
	if err != nil {
		return user, err
	}

	return user, nil
}

// FindCityConc provides a get city by name from database (concurrently).
func (db *DB) FindCityConc(cityRaw string, chErr chan<- error, cityCh chan<- models.City) {
	var (
		cityName    string
		countryName string
		city        models.City
	)

	cityRawSplit := strings.Split(cityRaw, ",")
	if len(cityRawSplit) > 1 {
		cityName = strings.TrimSpace(cityRawSplit[0])
		countryName = strings.TrimSpace(cityRawSplit[1])
	} else {
		cityName = strings.TrimSpace(cityRaw)
	}

	err := db.SQLX.Get(&city, `SELECT cid, name, name_ascii, country_code, 
	country, timezone, latitude, longitude FROM cities 
	WHERE (name = ? OR alternative_names LIKE ?) 
	AND country LIKE ?;`, cityName, "%"+cityName+",%", countryName+"%")
	if err != nil {
		chErr <- err
	}

	cityCh <- city
}

// FindCity provides a get city by name from database.
func (db *DB) FindCity(cityRaw string) (models.City, error) {
	var (
		cityName    string
		countryName string
		city        models.City
	)

	cityRawSplit := strings.Split(cityRaw, ",")
	if len(cityRawSplit) > 1 {
		cityName = strings.TrimSpace(cityRawSplit[0])
		countryName = strings.TrimSpace(cityRawSplit[1])
	} else {
		cityName = strings.TrimSpace(cityRaw)
	}

	err := db.SQLX.Get(&city, `SELECT cid, name, name_ascii, country_code, 
	country, timezone, latitude, longitude FROM cities 
	WHERE (name = ? OR alternative_names LIKE ?) 
	AND country LIKE ?`, cityName, "%"+cityName+",%", countryName+"%")
	if err != nil {
		return city, err
	}

	return city, nil
}

// FindObjectsNearByName performs search for all objects at a distance
// until n km from the object by name.
// Returns city of departure, list of objects (cities) near the city and error.
func (db *DB) FindObjectsNearByName(departure string, distance int) (models.City, []models.City, error) {
	city, err := db.FindCity(departure)
	if err != nil {
		return city, nil, err
	}

	// convert float to uint
	latUint, lonUint := uint(city.Latitude*converFact), uint(city.Longitude*converFact)

	// convert km to degree
	degreeLat := float64(distance) / oneDegreesInKmLat
	degreeLon := float64(distance) / oneDegreesInKmLon * math.Cos(degreeLat)
	if degreeLon < 0 {
		degreeLon = -degreeLon
	}

	var cities []models.City

	err = db.SQLX.Select(&cities, `SELECT cid, name, country, 
		latitude, longitude FROM cities HAVING 
		ABS(CAST((latitude * ? - ?) AS INT)) <= ? AND 
		ABS(CAST((longitude * ? - ?) AS INT)) <= ?`,
		converFact, latUint, uint(degreeLat*converFact),
		converFact, lonUint, uint(degreeLon*converFact))
	if err != nil {
		return city, nil, err
	}

	return city, cities, nil
}

// FindObjectsNearByCoordperforms search for all objects at a distance
// until n km from the object by coordinates.
// Returns list of objects (cities) near these coordinates and error.
func (db *DB) FindObjectsNearByCoord(lat float64, lon float64, distance int) ([]models.City, error) {
	// convert float to uint
	latUint, lonUint := uint(lat*converFact), uint(lon*converFact)

	// convert km to degree
	degreeLat := float64(distance) / oneDegreesInKmLat
	degreeLon := float64(distance) / oneDegreesInKmLon * math.Cos(degreeLat)
	if degreeLon < 0 {
		degreeLon = -degreeLon
	}

	var cities []models.City

	err := db.SQLX.Select(&cities, `SELECT cid, name, country, 
		latitude, longitude FROM cities HAVING 
		ABS(CAST((latitude * ? - ?) AS INT)) <= ? AND 
		ABS(CAST((longitude * ? - ?) AS INT)) <= ?`,
		converFact, latUint, uint(degreeLat*converFact),
		converFact, lonUint, uint(degreeLon*converFact))
	if err != nil {
		return nil, err
	}

	return cities, nil
}

// checkTableExist checks if the table exists and returns
// false if it does not exist.
func (db *DB) checkTableExist(tableName string) bool {
	var res int
	err := db.SQLX.Get(&res, `SELECT COUNT(*) FROM 
	INFORMATION_SCHEMA.TABLES 
	WHERE TABLE_NAME = ?`, tableName)

	if res == 0 && err != nil {
		return false
	}

	return true
}

//
