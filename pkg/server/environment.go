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
	"github.com/moov-io/base/config"
	"github.com/moov-io/base/database"
	"github.com/moov-io/base/log"
	"github.com/moov-io/base/stime"
)

// Environment - Contains everything thats been instantiated for this service.
type Environment struct {
	Logger       log.Logger
	Config       *Config
	TimeService  *stime.TimeService
	PublicRouter *mux.Router
	Shutdown     func()
}

// NewEnvironment - Generates a new default environment. Overrides can be specified via configs.
func NewEnvironment(env *Environment) (*Environment, error) {
	if env.Logger == nil {
		env.Logger = log.NewDefaultLogger()
	}

	if env.Config == nil {
		ConfigService := config.NewService(env.Logger)

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

	// configure custom handlers
	ConfigureHandlers(env.PublicRouter)

	env.Shutdown = func() {
		close()
	}

	return env, nil
}

func initializeDatabase(logger log.Logger, config database.DatabaseConfig) (*sql.DB, func(), error) {
	ctx, cancelFunc := context.WithCancel(context.Background())

	// migrate database
	db, err := database.New(ctx, logger, config)
	if err != nil {
		return nil, cancelFunc, logger.Fatal().LogErrorf("Error creating database", err).Err()
	}

	shutdown := func() {
		logger.Info().Log("Shutting down the db")
		cancelFunc()
		if err := db.Close(); err != nil {
			logger.Fatal().LogErrorf("Error closing DB", err).Err()
		}
	}

	backupFiles, _ := ioutil.ReadDir(filepath.Join("migrations"))
	if len(backupFiles) > 0 {
		if err := database.RunMigrations(logger, config); err != nil {
			return nil, shutdown, logger.Fatal().LogErrorf("Error running migrations", err).Err()
		}
	} else {
		logger.Info().Log("there is no backup files of database")
	}

	logger.Info().Log("finished initializing db")

	return db, shutdown, err
}
