// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"encoding/json"
	"fmt"

	"cloudeng.io/errors"
	"github.com/cosnicolaou/protocolsio/api"
)

type ProtocolsGetFlags struct{}

func protocolsGetCmd(ctx context.Context, values interface{}, args []string) error {
	errs := errors.M{}
	for _, id := range args {
		_, body, err := getProtocol(ctx, id)
		if err != nil {
			errs.Append(err)
			continue
		}
		fmt.Printf("%s\n", body)
	}
	return errs.Err()
}

func getProtocol(ctx context.Context, id string) (json.RawMessage, []byte, error) {
	ctx = globalConfig.WithAuth(ctx)
	u := globalConfig.Endpoints.GetProtocolV4 + "/" + id
	resp, body, err := api.Get[api.Payload](ctx, u)
	if resp.StatusCode != 0 {
		return nil, body, fmt.Errorf("unexpected status_code: %v", resp.StatusCode)
	}
	return resp.Payload, body, err
}
