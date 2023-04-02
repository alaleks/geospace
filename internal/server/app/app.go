// Package app perfoms run application
package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alaleks/geospace/internal/server/app/authentication"
	"github.com/alaleks/geospace/internal/server/app/handlers"
	"github.com/alaleks/geospace/internal/server/config"
	"github.com/alaleks/geospace/internal/server/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type App struct {
	cfg    *config.Cfg        // configuration
	srv    *fiber.App         // server
	hdls   *handlers.Hdls     // handlers
	logger *zap.SugaredLogger // zap logger
}

// New returns a pointer to a new App instance.
func New() *App {
	logger, err := createLogger()
	if err != nil {
		log.Fatal(err)
	}

	app := &App{
		logger: logger,
	}

	cfg, err := config.New(logger)
	if err != nil {
		logger.Fatal(err)
	}

	db, err := database.Connect(*cfg)
	if err != nil {
		logger.Fatal(err)
	}

	// migrate schemes of tables
	db.Migrate()

	// import data to table
	err = importCities(db.SQLX)
	if err != nil {
		logger.Fatal(err)
	}

	// create server and handlers
	app.cfg = cfg
	app.createServer()
	app.hdls = handlers.New(db, authentication.Init(db, cfg.Secure))

	return app
}

// Run performs start application.
func (app *App) Run() {
	// run goroutine for catch os signals for shutdown server.
	go app.catchSign()

	// register routes
	app.RegRouters()

	// use recovery from panic
	app.srv.Use(recover.New())

	// listen port
	err := app.srv.Listen(app.cfg.App.Port)
	if err != nil {
		app.logger.Fatal(err)
	}
}

// RegRouters install routes for the given application.
func (app *App) RegRouters() {
	// ping server
	app.srv.Get("/ping", app.hdls.Ping)

	// v1
	v1 := app.srv.Group("/v1")
	// registration for the using application.
	v1.Post("/register", app.hdls.SignUp)
	// login for the using application.
	v1.Post("/login", app.hdls.Login)
	// logout user
	v1.Get("/logout", app.hdls.Logout)
	// get list countries
	v1.Get("/country", app.hdls.GetCountry)

	// these routes available only auth user
	user := v1.Group("/user", app.hdls.CheckAuthentication)
	user.Get("/distance", app.hdls.CalculateDistance)
	//user.Get("/cities-area")

	// api, these routes available only auth user
	api := v1.Group("/api", app.hdls.CheckAuthentication)
	api.Get("/distance", app.hdls.CalculateDistanceAPI)
	api.Get("/find-by-name", app.hdls.FindObjectsNearByNameAPI)
	api.Get("/find-by-coord", app.hdls.FindObjectsNearByCoordAPI)
}

// catchSign will catch SIGINT, SIGHUP, SIGQUIT and SIGTERM and shutdown the server.
func (app *App) catchSign() {
	termSignals := make(chan os.Signal, 1)
	reloadSignals := make(chan os.Signal, 1)

	signal.Notify(termSignals,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	signal.Notify(reloadSignals, syscall.SIGUSR1)

	for {
		select {
		case <-termSignals:
			fmt.Printf("%s shutdown\n", app.cfg.App.Name)
			err := app.srv.Shutdown()
			if err != nil {
				app.logger.Fatal(err)
			}
		case <-reloadSignals:
			fmt.Printf("%s shutdown\n", app.cfg.App.Name)
			err := app.srv.Shutdown()
			if err != nil {
				app.logger.Fatal(err)
			}
		}
	}
}

// createServer performs initialization a new server.
func (app *App) createServer() {
	app.srv = fiber.New(fiber.Config{
		AppName: app.cfg.App.Name,
	})
}

// createLogger performs initialization a new logger.
func createLogger() (*zap.SugaredLogger, error) {
	cfgZap := zap.NewProductionConfig()
	cfgZap.EncoderConfig.TimeKey = "timestamp"
	cfgZap.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("02.01.2006 15:04:05")
	cfgZap.EncoderConfig.StacktraceKey = ""

	logger, err := cfgZap.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
