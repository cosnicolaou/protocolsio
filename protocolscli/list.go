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

type ProtocolsListFlags struct {
	ProtocolCommonFlags
	Query string `subcmd:"query,,'string may contain any characters, numbers and special symbols. System will seach around protocol name, description, authors. If the search keywords are enclosed into double quotes, then result contains only the exact match of the combined term'"`
	Order string `subcmd:"order,activity,'one of: 1. activity - index of protocol popularity; relevance - returns most relevant to the search key results at the top; date - date of publication; name - protocol name; id - id of protocol.'"`
	Sort  string `subcmd:"sort,asc,one of asc or desc"`
}

type protocolItemProcessor interface {
	Process(context.Context, api.ListProtocolsV3, checkpoint) error
}

type itemPrinter struct{}

func (ip *itemPrinter) Process(ctx context.Context, protocols api.ListProtocolsV3, cp checkpoint) error {
	errs := errors.M{}
	for _, item := range protocols.Items {
		var p api.Protocol
		if err := json.Unmarshal(item, &p); err != nil {
			errs.Append(err)
			continue
		}
		fmt.Printf("%v: URI: %v, Title: %v\n", p.ID, p.URI, p.Title)
	}
	return nil
}

func protocolsListCmd(ctx context.Context, values interface{}, args []string) error {
	ck, err := newCheckpointFromFlags(values.(*ProtocolsListFlags))
	if err != nil {
		return err
	}
	return getProtocols(ctx, ck, &itemPrinter{})
}
