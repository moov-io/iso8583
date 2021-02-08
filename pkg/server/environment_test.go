// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package server_test

import (
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	baseLog "github.com/moov-io/base/log"
	"github.com/moov-io/iso8583/pkg/server"
	"github.com/stretchr/testify/assert"
)

func Test_Environment_Startup(t *testing.T) {
	a := assert.New(t)

	env := &server.Environment{
		Logger: baseLog.NewLogger(log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))),
	}

	env, err := server.NewEnvironment(env)
	a.Nil(err)

	shutdown := env.RunServers(false)

	env1, err := server.NewEnvironment(&server.Environment{})
	a.Nil(err)
	env1.Shutdown()

	t.Cleanup(shutdown)
}
