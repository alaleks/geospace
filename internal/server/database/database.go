// Package database performs operations over the database.
package database

import (
	"errors"
	"time"

	"github.com/alaleks/geospace/internal/server/config"
	"github.com/alaleks/geospace/internal/server/database/models"
	"github.com/alaleks/geospace/internal/server/database/schema"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	MaxIdleConns    = 100
	ConnMaxLifetime = 15 * time.Minute
	tableNameCities = "cities"
	tableNameUsers  = "users"
)

var (
	ErrUserAlreadyExists = errors.New("user with current email already exists")
)

type DB struct {
	SQLX *sqlx.DB
}

// Connect performs creating a new connection to database.
func Connect(cfg config.Cfg) (*DB, error) {
	db, err := sqlx.Connect("mysql", cfg.CreateDSN())
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(MaxIdleConns)
	db.SetConnMaxLifetime(ConnMaxLifetime)

	return &DB{
		SQLX: db,
	}, nil
}

// Migrate performs a migration schema of table to database.
func (db *DB) Migrate() {
	if !db.checkTableExist(tableNameCities) {
		db.SQLX.MustExec(schema.City)
	}

	if !db.checkTableExist(tableNameUsers) {
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
	db.SQLX.Get(&res, `SELECT COUNT(*) FROM users
	WHERE email = ?`, email)
	if res > 0 {
		return 0, ErrUserAlreadyExists
	}

	user := models.User{
		Name:      name,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now().Unix(),
	}

	_, err := db.SQLX.NamedExec(`INSERT INTO users (name, email, password, created_at) 
	VALUES (:name, :email, :password, :created_at)`,
		&user)

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

// FindCityByName provides a get city by name from database.
func (db *DB) FindCityByName(cityName, countryCode string) (models.City, error) {
	var city models.City
	err := db.SQLX.Get(&city, `SELECT * FROM cities 
		WHERE name = ?`, cityName)
	if err == nil {
		return city, nil
	}

	err = db.SQLX.Get(&city, `SELECT * FROM cities 
		WHERE alternative_names LIKE ?`, "%"+cityName+",%")
	if err != nil {
		return city, err
	}

	return city, nil
}

// FindCityByNameAndCountryCode returns a city by name and country code.
func (db *DB) FindCityByNameAndCountryCode(cityName, countryCode string) (models.City, error) {
	var city models.City
	err := db.SQLX.Get(&city, `SELECT * FROM cities 
		WHERE name = ? AND country_code = ?`, cityName, countryCode)
	if err == nil {
		return city, nil
	}

	err = db.SQLX.Get(&city, `SELECT * FROM cities 
		WHERE alternative_names LIKE ? AND country_code = ?`, "%"+cityName+",%", countryCode)
	if err != nil {
		return city, err
	}

	return city, nil
}

// checkTableExist checks if the table exists and returns
// false if it does not exist.
func (db *DB) checkTableExist(tableName string) bool {
	var res int
	db.SQLX.Get(&res, `SELECT COUNT(*) FROM 
	INFORMATION_SCHEMA.TABLES 
	WHERE TABLE_NAME = ?`, tableName)

	if res == 0 {
		return false
	}

	return true
}
