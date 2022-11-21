// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"cloudeng.io/cmdutil/flags"
	"cloudeng.io/errors"
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

type ProtocolsListFlags struct {
	ProtocolCommonFlags
	Query string `subcmd:"query,,'string may contain any characters, numbers and special symbols. System will seach around protocol name, description, authors. If the search keywords are enclosed into double quotes, then result contains only the exact match of the combined term'"`
	Order string `subcmd:"order,activity,'one of: 1. activity - index of protocol popularity; relevance - returns most relevant to the search key results at the top; date - date of publication; name - protocol name; id - id of protocol.'"`
	Sort  string `subcmd:"sort,asc,one of asc or desc"`
}

type ProtocolsDownloadFlags struct {
	ProtocolsListFlags
	CacheDir       string `subcmd:"cachepath,,'location of cache of download protocol objects that overides that specified in the global yaml config'"`
	CheckpointFile string `subcmd:"resume,,checkpoint file to resume download from"`
}

func stringPtr(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

type protocolItemProcessor interface {
	Process([]api.Item, *checkpoint) error
}

type itemPrinter struct{}

func (ip *itemPrinter) Process(items []api.Item, cp *checkpoint) error {
	for _, item := range items {
		fmt.Printf("%v: URI: %v, Title: %v, DOI: %v\n", item.ID, item.URI, item.Title, stringPtr(item.Doi))
	}
	return nil
}

func protocolsListCmd(ctx context.Context, values interface{}, args []string) error {
	ck, err := newCheckpointFromFlags(values.(*ProtocolsListFlags))
	if err != nil {
		return err
	}
	return getProtocols(ctx, ck, args, &itemPrinter{})
}

func protocolsDownloadCmd(ctx context.Context, values interface{}, args []string) error {
	fv := values.(*ProtocolsDownloadFlags)
	dir := fv.CacheDir
	if len(dir) == 0 {
		dir = globalConfig.Cache.Path
	}
	if len(dir) == 0 {
		return fmt.Errorf("no cache path specified either via --cachepath or via the global yaml config file")
	}
	saver, err := newItemSaver(dir)
	if err != nil {
		return err
	}
	var cp *checkpoint
	if len(fv.CheckpointFile) != 0 {
		data, err := os.ReadFile(fv.CheckpointFile)
		if err != nil {
			return fmt.Errorf("failed to read checkpoint: %v", err)
		}
		cp = &checkpoint{}
		if err := json.Unmarshal(data, cp); err != nil {
			return fmt.Errorf("failed to decode checkpoint file: %v: %v", fv.CheckpointFile, err)
		}
	} else {
		cp, err = newCheckpointFromFlags(&fv.ProtocolsListFlags)
		if err != nil {
			return err
		}
	}
	return getProtocols(ctx, cp, args, saver)
}

type downloadedItems struct {
	items      *[]api.Item
	checkpoint *checkpoint
}

func getProtocols(ctx context.Context, checkpoint *checkpoint, args []string, proccessor protocolItemProcessor) error {
	ch := make(chan downloadedItems, 1000)
	errCh := make(chan error)

	go func() {
		errCh <- getProtocolsCall(ctx, checkpoint, args, ch)
		close(ch)
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return err
		case downloaded, ok := <-ch:
			if !ok {
				return nil
			}
			if err := proccessor.Process(*downloaded.items, downloaded.checkpoint); err != nil {
				return err
			}
		}
	}
}

func getProtocolsCall(ctx context.Context, checkpoint *checkpoint, args []string, ch chan<- downloadedItems) error {
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

	initialDelay := time.Minute
	maxDelay := time.Minute * 16
	delay := initialDelay
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		resp, _, err := api.Get[api.ProtocolsV3](ctx, u+v.Encode())
		if err != nil {
			if errors.Is(err, api.ErrTooManyRequests) {
				if delay >= maxDelay {
					return fmt.Errorf("failed after retrying for %v", delay)
				}
				fmt.Printf("too many requests: sleeping for %v\n", delay)
				time.Sleep(delay)
				delay *= 2
				continue
			}
			return err
		}
		if delay != initialDelay {
			fmt.Printf("succeeded after retry with delay of %v\n", delay)
		}
		delay = initialDelay
		var result downloadedItems

		done, nextPage, err := checkpoint.update(&resp.Pagination)
		if err != nil {
			return err
		}
		result.checkpoint = checkpoint

		if checkpoint.Total > 0 {
			nItems += len(resp.Items)
			if nItems <= checkpoint.Total {
				result.items = &resp.Items
				ch <- result
			} else {
				rem := resp.Items[:checkpoint.Total-(nItems-len(resp.Items))]
				result.items = &rem
				ch <- result
			}
			if nItems >= checkpoint.Total {
				break
			}
		} else {
			result.items = &resp.Items
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
	return nil
}
