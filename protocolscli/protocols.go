// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"cloudeng.io/cmdutil/flags"
	"github.com/cosnicolaou/protocolsio/api"
)

func ProtocolListFilters() []string {
	return []string{"public", "user_public", "user_private", "shared_with_user"}
}

func ProtocolListOrderField() []string {
	return []string{"activity", "relevance", "date", "name", "id"}
}

func ProtocolListOrderDirection() []string {
	return []string{"asc", "desc"}
}

type ProtocolCommonFlags struct {
	Pages    flags.IntRangeSpec `subcmd:"pages,1,page range to return"`
	PageSize int                `subcmd:"size,20,number of items in each page"`
	Total    int                `subcmd:"total,,total number of items to return"`
	Filter   string             `subcmd:"filter,public,'one of: 1. public - list of all public protocols;	2. user_public - list of public protocols that was publiches by concrete user; 3. user_private - list of private protocols that was created by concrete user; 4. shared_with_user - list of public protocols that was shared with concrete user.'"`
}

type downloadedItems struct {
	err        error
	protocols  api.ListProtocolsV3
	checkpoint checkpoint
}

func getProtocols(ctx context.Context, checkpoint checkpoint, proccessor protocolItemProcessor) error {
	ch := make(chan downloadedItems, 1000)

	go func() {
		getProtocolsCall(ctx, checkpoint, ch)
		close(ch)
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case downloaded, ok := <-ch:
			if !ok {
				return nil
			}
			if err := downloaded.err; err != nil {
				return err
			}
			if err := proccessor.Process(ctx, downloaded.protocols, downloaded.checkpoint); err != nil {
				return err
			}
		}
	}

}

func getProtocolsCall(ctx context.Context, checkpoint checkpoint, ch chan<- downloadedItems) {
	lastPage := strconv.Itoa(checkpoint.Pages.To)
	if checkpoint.Pages.To == 0 && !checkpoint.Pages.ExtendsToEnd {
		lastPage = strconv.Itoa(checkpoint.Pages.From)
	}
	ctx = globalConfig.WithAuth(ctx)
	u := globalConfig.Endpoints.ListProtocolsV3
	u = strings.TrimSuffix(u, "/?")
	u += "?"
	v := url.Values{}
	checkpoint.initHeaders(&v)
	nItems := 0

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		var result downloadedItems

		resp, _, err := api.Get[api.ListProtocolsV3](ctx, u+v.Encode())
		if err != nil {
			result.err = err
			ch <- result
			return
		}

		done, nextPage, err := checkpoint.update(resp.Pagination)
		if err != nil {
			ch <- result
			return
		}
		result.checkpoint = checkpoint
		result.protocols.Extras = resp.Extras

		if checkpoint.Total > 0 {
			nItems += len(resp.Items)
			if nItems <= checkpoint.Total {
				result.protocols.Items = resp.Items
				ch <- result
			} else {
				rem := resp.Items[:checkpoint.Total-(nItems-len(resp.Items))]
				result.protocols.Items = rem
				ch <- result
			}
			if nItems >= checkpoint.Total {
				break
			}
		} else {
			result.protocols.Items = resp.Items
			ch <- result
		}
		if done {
			break
		}
		if checkpoint.Total == 0 {
			if !checkpoint.Pages.ExtendsToEnd && v.Get("page_id") == lastPage {
				break
			}
		}
		v.Set("page_id", strconv.Itoa(nextPage))
	}
	return
}
