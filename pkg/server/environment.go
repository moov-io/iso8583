// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package server

import (
	"context"
	"database/sql"
	"io/ioutil"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/moov-io/identity/pkg/config"
	"github.com/moov-io/identity/pkg/database"
	"github.com/moov-io/identity/pkg/logging"
	"github.com/moov-io/identity/pkg/stime"
	tmw "github.com/moov-io/tumbler/pkg/middleware"
	"github.com/moov-io/tumbler/pkg/webkeys"
)

// Environment - Contains everything thats been instantiated for this service.
type Environment struct {
	Logger       logging.Logger
	Config       *Config
	TimeService  *stime.TimeService
	GatewayKeys  webkeys.WebKeysService
	PublicRouter *mux.Router
	Shutdown     func()
}

// NewEnvironment - Generates a new default environment. Overrides can be specified via configs.
func NewEnvironment(env *Environment) (*Environment, error) {
	if env.Logger == nil {
		env.Logger = logging.NewDefaultLogger()
	}

	if env.Config == nil {
		ConfigService := config.NewConfigService(env.Logger)

		global := &GlobalConfig{}
		if err := ConfigService.Load(global); err != nil {
			return nil, err
		}

		env.Config = &global.ISO8583
	}

	//db setup
	db, close, err := initializeDatabase(env.Logger, env.Config.Database)
	if err != nil {
		close()
		return nil, err
	}
	_ = db // delete once used.

	if env.TimeService == nil {
		t := stime.NewSystemTimeService()
		env.TimeService = &t
	}

	// router
	if env.PublicRouter == nil {
		env.PublicRouter = mux.NewRouter()
	}

	// auth middleware for the tokens coming from the gateway
	GatewayMiddleware, err := tmw.NewTumblerMiddlewareFromConfig(env.Logger, *env.TimeService, env.Config.Gateway)
	if err != nil {
		return nil, env.Logger.Fatal().LogErrorF("Can't startup the Gateway middleware - %w", err)
	}

	env.PublicRouter.Use(GatewayMiddleware.Handler)

	// configure custom handlers
	ConfigureHandlers(env.PublicRouter)

	env.Shutdown = func() {
		close()
	}

	return env, nil
}

func initializeDatabase(logger logging.Logger, config database.DatabaseConfig) (*sql.DB, func(), error) {
	ctx, cancelFunc := context.WithCancel(context.Background())

	// migrate database
	db, err := database.New(ctx, logger, config)
	if err != nil {
		return nil, cancelFunc, logger.Fatal().LogError("Error creating database", err)
	}

	shutdown := func() {
		logger.Info().Log("Shutting down the db")
		cancelFunc()
		if err := db.Close(); err != nil {
			logger.Fatal().LogError("Error closing DB", err)
		}
	}

	backupFiles, _ := ioutil.ReadDir(filepath.Join("migrations"))
	if len(backupFiles) > 0 {
		if err := database.RunMigrations(logger, db, config); err != nil {
			return nil, shutdown, logger.Fatal().LogError("Error running migrations", err)
		}
	} else {
		logger.Info().Log("there is no backup files of database")
	}

	logger.Info().Log("finished initializing db")

	return db, shutdown, err
}
