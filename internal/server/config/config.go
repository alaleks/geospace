// Package config performs configuration of the application.
package config

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/alaleks/geospace/pkg/genkey"
	"github.com/golang-module/dongle"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

const (
	cfgDirName    = "/cfg"
	cfgFile       = "/config.yaml"
	sizeIVCipher  = 8
	sizeKeyCipher = 48
	sizeKeySecret = 64
)

type (
	// Cfg contains the configuration of app.
	Cfg struct {
		CfgDatabase CfgDatabase `yaml:"database"`
		App         App         `yaml:"app"`
		Secure      Secure      `yaml:"secure"`
	}

	// CfgDatabase contains the configuration for a database connection.
	CfgDatabase struct {
		Name       string `yaml:"name"`        // Name of the database
		User       string `yaml:"user"`        // User name of the database
		Password   string `yaml:"password"`    // Password of the database
		UnixSocket string `yaml:"unix_socket"` // Socket for connections (faster than TCP connection)
		Port       int    `yaml:"port"`        // Port of the database for TCP connections
	}

	// App contains the params of settings.
	App struct {
		Name       string `yaml:"name"`        // Name of the application
		Port       string `yaml:"port"`        // Port for running the application
		MaxRequest int    `yaml:"max_request"` // Max request quantity in seconds
		Expiration int    `yaml:"expiration"`  // Expiration period in seconds
	}

	// Secure contains the params for encryption
	// and decryption private data.
	Secure struct {
		SecretJWT string `yaml:"secret_jwt"` // SecretJWT needed create JWT token.
		Key       string `yaml:"key"`        // Key needed for create new cipher.
		IV        string `yaml:"iv"`         // IV  needed for create new cipher.
	}
)

// New create config and returns a pointer to the Cfg struct and error.
func New(logger *zap.SugaredLogger) (*Cfg, error) {
	cfg := new(Cfg)
	// read config file
	if err := cfg.readCfgFile(); err == nil {
		return cfg, nil
	}

	// default parameters
	cfg.App = App{
		Name:       "geospace",
		Port:       ":3000",
		MaxRequest: 100,
		Expiration: 1,
	}

	// parce flags
	cfg.readParamFlags()
	err := cfg.validateConfig()
	if err != nil {
		return cfg, err
	}

	// check connection with database
	// this check only then config file not existing.
	db, err := sqlx.Connect("mysql", cfg.CreateDSN())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.Close()

	// generate keys for encryption/decryption
	cfg.Secure = Secure{
		Key: dongle.Encode.FromString(genkey.Create(sizeKeyCipher)).
			ByBase64().ToString(),
		IV: dongle.Encode.FromString(genkey.Create(sizeIVCipher)).
			ByBase64().ToString(),
		SecretJWT: dongle.Encode.FromString(genkey.Create(sizeKeySecret)).
			ByBase64().ToString(),
	}

	// create config file
	err = cfg.createCfgFile()
	if err != nil {
		// if can't create the config file but config is valid
		// write about this error in log and using parameters of configuration
		// in only current session
		logger.Error(err.Error())
	}

	return cfg, nil
}

// CreateDSN returns a string for connecting to the database.
func (cfg *Cfg) CreateDSN() string {
	// if indicated unix socket then we make a connection through it.
	if cfg.CfgDatabase.UnixSocket != "" {
		return fmt.Sprintf("%s:%s@unix(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.CfgDatabase.User, cfg.CfgDatabase.Password,
			cfg.CfgDatabase.UnixSocket, cfg.CfgDatabase.Name)
	}

	return fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=%d",
		cfg.CfgDatabase.User, cfg.CfgDatabase.Password,
		cfg.CfgDatabase.Name, cfg.CfgDatabase.Port)
}

// GetKeyCipher returns key after decrypt.
func (s *Secure) GetKeyCipher() string {
	return dongle.Decode.FromString(s.Key).ByBase64().ToString()
}

// GetIVCipher returns IV value after decrypt.
func (s *Secure) GetIVCipher() string {
	return dongle.Decode.FromString(s.IV).ByBase64().ToString()
}

// GetSecretJWT returns SecretJWT  after decrypt.
func (s *Secure) GetSecretJWT() string {
	return dongle.Decode.FromString(s.SecretJWT).ByBase64().ToString()
}

// GetRootDir returns the root directory of the project.
func GetRootDir() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	rootDir, err := os.Open(path.Join(currentDir, "../../"))
	if err != nil {
		return "", err
	}

	return rootDir.Name(), nil
}

// validateConfig performs validation of the configuration and
// returns error if required config parameters are not valid.
func (cfg *Cfg) validateConfig() error {
	switch {
	case len(cfg.CfgDatabase.Name) == 0:
		return fmt.Errorf("database name cannot be empty")
	case len(cfg.CfgDatabase.User) == 0:
		return fmt.Errorf("database username cannot be empty")
	case len(cfg.CfgDatabase.Password) == 0:
		return fmt.Errorf("database password cannot be empty")
	}

	if len(cfg.CfgDatabase.UnixSocket) == 0 && cfg.CfgDatabase.Port == 0 {
		return fmt.Errorf("unix socket or port of database cannot be empty")
	}

	return nil
}

// readCfgFile performs read configuration from yaml file.
func (cfg *Cfg) readCfgFile() error {
	// get root directory
	rootDir, err := GetRootDir()
	if err != nil {
		return err
	}

	// create path to cfg dir
	cfgPath := rootDir + cfgDirName

	// check dir of cfg exists
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return err
	}

	// open cfg file
	f, err := os.Open(cfgPath + cfgFile)
	if err != nil {
		return err
	}

	defer f.Close()

	// read file to buffer
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(f)
	if err != nil {
		return err
	}

	// unmarshal config from yaml
	return yaml.Unmarshal(buf.Bytes(), cfg)
}

// readCfgFile performs read configuration from yaml file.
func (cfg *Cfg) createCfgFile() error {
	// get root directory
	rootDir, err := GetRootDir()
	if err != nil {
		return err
	}

	// create cfg dir
	cfgPath := rootDir + cfgDirName
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		err = os.Mkdir(cfgPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// marshal config to yaml
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	// write config to file
	return os.WriteFile(cfgPath+cfgFile, b, os.ModePerm)
}

// readParamFlags performs reading parameters from flags the command line.
func (cfg *Cfg) readParamFlags() {
	var (
		// required parameters
		dbName   = flag.String("d", "", "Name of the database")
		dbUser   = flag.String("u", "", "User name of the database")
		dbPass   = flag.String("p", "", "Password of the database")
		dbSocket = flag.String("s", "", "Socket of connection to database")
		dbPort   = flag.Int("t", 0, "Port of the database")

		// optional parameters
		appName    = flag.String("n", "", "Name of the database")
		port       = flag.Int("a", 0, "Port for running the application")
		maxRequest = flag.Int("r", 0, "Max request quantity in seconds")
		expiration = flag.Int("e", 0, "Expiration period in seconds")
	)

	flag.Parse()

	cfg.CfgDatabase = CfgDatabase{
		Name:       *dbName,
		User:       *dbUser,
		Password:   *dbPass,
		UnixSocket: *dbSocket,
		Port:       *dbPort,
	}

	if *port != 0 {
		cfg.App.Port = fmt.Sprintf(":%d", *port)
	}

	if *maxRequest != 0 {
		cfg.App.MaxRequest = *maxRequest
	}

	if *expiration != 0 {
		cfg.App.Expiration = *expiration
	}

	if *appName != "" {
		cfg.App.Name = *appName
	}
}
