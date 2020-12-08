// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/moov-io/base/admin"
	_ "github.com/moov-io/identity" // need to import the embedded files

	log "github.com/moov-io/identity/pkg/logging"
)

// RunServers - Boots up all the servers and awaits till they are stopped.
func (env *Environment) RunServers(await bool) func() {

	// Listen for application termination.
	terminationListener := newTerminationListener()

	adminServer := bootAdminServer(terminationListener, env.Logger, env.Config.Servers.Admin)

	_, shutdownPublicServer := bootHTTPServer("public", env.PublicRouter, terminationListener, env.Logger, env.Config.Servers.Public)

	if await {
		awaitTermination(env.Logger, terminationListener)
	}

	return func() {
		adminServer.Shutdown()
		shutdownPublicServer()
	}
}

func newTerminationListener() chan error {
	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	return errs
}

func awaitTermination(logger log.Logger, terminationListener chan error) {
	if err := <-terminationListener; err != nil {
		logger.Fatal().LogError("Terminated", err)
	}
}

func bootHTTPServer(name string, routes *mux.Router, errs chan<- error, logger log.Logger, config HTTPConfig) (*http.Server, func()) {

	// Create main HTTP server
	serve := &http.Server{
		Addr:    config.Bind.Address,
		Handler: routes,
		TLSConfig: &tls.Config{
			InsecureSkipVerify:       false,
			PreferServerCipherSuites: true,
			MinVersion:               tls.VersionTLS12,
		},
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start main HTTP server
	go func() {
		logger.Info().Log(fmt.Sprintf("%s listening on %s", name, config.Bind.Address))
		if err := serve.ListenAndServe(); err != nil {
			errs <- logger.Fatal().LogErrorF("problem starting http: %w", err)
		}
	}()

	shutdownServer := func() {
		if err := serve.Shutdown(context.TODO()); err != nil {
			logger.Fatal().LogError(name, err)
		}
	}

	return serve, shutdownServer
}

func bootAdminServer(errs chan<- error, logger log.Logger, config HTTPConfig) *admin.Server {
	adminServer := admin.NewServer(config.Bind.Address)

	go func() {
		logger.Info().Log(fmt.Sprintf("listening on %s", adminServer.BindAddr()))
		if err := adminServer.Listen(); err != nil {
			errs <- logger.Fatal().LogErrorF("problem starting admin http: %w", err)
		}
	}()

	return adminServer
}
