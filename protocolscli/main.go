// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"os"
	"strings"

	"cloudeng.io/cmdutil/signals"
	"cloudeng.io/cmdutil/subcmd"
	"github.com/cosnicolaou/protocolsio/protocolscli/glean"
)

type GlobalFlags struct {
	Config string `subcmd:"config,$HOME/.protocolsio.yaml,'protocolsio config file'"`
}

var (
	globalFlags  GlobalFlags
	globalConfig *Config
	cmdSet       *subcmd.CommandSetYAML
)

func indent(indent string, document string) string {
	out := strings.Builder{}
	for _, line := range strings.Split(document, "\n") {
		out.WriteString(indent)
		out.WriteString(line)
		out.WriteByte('\n')
	}
	return out.String()
}

var yamlSpec = `name: protocolsio
summary: utility for accessing protocols.io via its API
commands:
  - name: protocols
    summary: utilities for working with protocol objects
    commands:
      - name: list
        summary: list protocols
      - name: download
        summary: download protocols
      - name: get
        summary: get a specific protocol
        arguments:
          - id
          - ...
` + indent("  ", glean.SubcmdYAML)

func init() {
	cmdSet = subcmd.MustFromYAML(yamlSpec)
	globals := subcmd.GlobalFlagSet()
	globals.MustRegisterFlagStruct(&globalFlags, nil, nil)

	cmdSet.Set("protocols", "list").RunnerAndFlags(
		protocolsListCmd, subcmd.MustRegisteredFlagSet(&ProtocolsListFlags{}))

	cmdSet.Set("protocols", "download").RunnerAndFlags(
		protocolsDownloadCmd, subcmd.MustRegisteredFlagSet(&ProtocolsDownloadFlags{}))

	cmdSet.Set("protocols", "get").RunnerAndFlags(
		protocolsGetCmd, subcmd.MustRegisteredFlagSet(&ProtocolsGetFlags{}))

	glean.ConfigureCmdSet(cmdSet)
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
