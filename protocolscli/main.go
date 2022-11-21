// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"os"

	"cloudeng.io/cmdutil/signals"
	"cloudeng.io/cmdutil/subcmd"
)

type GlobalFlags struct {
	Config string `subcmd:"config,$HOME/.protocolsio.yaml,'config file'"`
}

var (
	globalFlags  GlobalFlags
	globalConfig *Config
	cmdSet       *subcmd.CommandSetYAML
)

func init() {
	cmdSet = subcmd.MustFromYAML(`name: protocolsio
summary: utility for accessing protocols.io via its API
commands:
  - name: protocols
    summary: utilities for working with protocol objects
    commands:
      - name: list
        summary: list protocols
      - name: download
        summary: download protocols
`)

	globals := subcmd.GlobalFlagSet()
	globals.MustRegisterFlagStruct(&globalFlags, nil, nil)

	cmdSet.Set("protocols", "list").RunnerAndFlags(
		protocolsListCmd, subcmd.MustRegisteredFlagSet(&ProtocolsListFlags{}))

	cmdSet.Set("protocols", "download").RunnerAndFlags(
		protocolsDownloadCmd, subcmd.MustRegisteredFlagSet(&ProtocolsDownloadFlags{}))

	cmdSet.WithGlobalFlags(globals)
	cmdSet.WithMain(mainWrapper)
}

func mainWrapper(ctx context.Context, cmdRunner func(ctx context.Context) error) error {
	cfg, err := ParseConfig(globalFlags.Config)
	if err != nil {
		return err
	}
	globalConfig = cfg
	return cmdRunner(ctx)
}

func main() {
	ctx, _ := signals.NotifyWithCancel(context.Background(), os.Interrupt)
	cmdSet.MustDispatch(ctx)
}
