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
	/*
		buf, err := os.ReadFile("f")
		if err != nil {
			panic(err)
		}
		var p api.Protocol
		if err := json.Unmarshal(buf, &p); err != nil {
			panic(err)
		}

		fmt.Printf("%#v\n", p)
		var a map[string]interface{}
		if err := json.Unmarshal(buf, &a); err != nil {
			panic(err)
		}
		for k := range a {
			fmt.Printf("%v\n", k)
		}
	return nil*/

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
