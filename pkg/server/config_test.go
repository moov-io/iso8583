// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package server_test

import (
	"testing"

	"github.com/moov-io/identity/pkg/config"
	"github.com/moov-io/identity/pkg/logging"

	"github.com/moov-io/iso8583/pkg/server"
	"github.com/stretchr/testify/require"
)

func Test_ConfigLoading(t *testing.T) {
	logger := logging.NewNopLogger()

	ConfigService := config.NewConfigService(logger)

	gc := &server.GlobalConfig{}
	err := ConfigService.Load(gc)
	require.Nil(t, err)
}
