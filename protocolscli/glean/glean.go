// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package glean

import (
	"cloudeng.io/cmdutil/subcmd"
)

// DatasourceName is the name of the protocols.io datasource
// registered with the Glean SDK.
const DatasourceName = "protocolsio"

var SubcmdYAML = `
- name: glean
  summary: Glean related commands
  commands:
    - name: bulk-index
      summary: index protocols.io protocol objects using Glean.
      arguments:
        - documents-directory - containing previously downloaded documents to be indexed.
    - name: stats
      summary: retrieve statistics for the protocolsio datasource.
`

func ConfigureCmdSet(cmdSet *subcmd.CommandSetYAML) {
	cmdSet.Set("glean", "bulk-index").RunnerAndFlags(
		bulkIndexCmd, subcmd.MustRegisteredFlagSet(&BulkIndexFlags{}))
	cmdSet.Set("glean", "stats").RunnerAndFlags(
		statsCmd, subcmd.MustRegisteredFlagSet(&StatsFlags{}))
}
