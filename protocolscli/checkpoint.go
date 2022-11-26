// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/url"
	"strconv"

	"cloudeng.io/cmdutil/flags"
	"cloudeng.io/errors"
	"github.com/cosnicolaou/protocolsio/api"
)

type checkpoint struct {
	CurrentPage int // Used for naming the checkpoint file
	TotalPages  int // Used for naming the checkpoint file
	// The following are set originally from command line flags.
	Pages      flags.IntRangeSpec // Range of pages equested, eg. 1-
	PageSize   int                // Number of itmes per page.
	Filter     string             // Which protocols to request
	FieldOrder string             //
	Order      string             //
	Total      int                // Total number of protocols to download.

	// Pagination returned by API server, used to resume the download.
	Pagination api.Pagination
	// Names of files this checkpoint covers.
	Files []string
}

func newCheckpointFromFlags(fv *ProtocolsListFlags) (checkpoint, error) {
	errs := errors.M{}
	errs.Append(flags.OneOf(fv.Filter).Validate("public", ProtocolListFilters()...))
	errs.Append(flags.OneOf(fv.Order).Validate("activity", ProtocolListOrderField()...))
	errs.Append(flags.OneOf(fv.Sort).Validate("asc", "desc"))
	if err := errs.Err(); err != nil {
		return checkpoint{}, err
	}
	return checkpoint{
		CurrentPage: 0,
		TotalPages:  0,
		Pages:       fv.Pages,
		PageSize:    fv.PageSize,
		Filter:      fv.Filter,
		FieldOrder:  fv.Order,
		Order:       fv.Sort,
		Total:       fv.Total,
	}, nil
}

func (cp *checkpoint) initHeaders(v *url.Values) {
	v.Add("page_size", strconv.Itoa(cp.PageSize))
	v.Add("filter", string(cp.Filter))
	v.Add("field_order", string(cp.FieldOrder))
	v.Add("order", string(cp.Order))
	v.Set("page_id", strconv.Itoa(cp.Pages.From))
}

func (cp *checkpoint) resetFiles() {
	cp.Files = make([]string, 0, 20)
}

func (cp *checkpoint) appendFile(file string) {
	cp.Files = append(cp.Files, file)
}

func (cp *checkpoint) update(p api.Pagination) (bool, int, error) {
	cp.CurrentPage = int(p.CurrentPage)
	cp.TotalPages = int(p.TotalPages)
	done := p.Done()
	if done {
		return done, 0, nil
	}
	u, err := url.Parse(p.NextPage)
	if err != nil {
		return done, 0, err
	}
	np := u.Query().Get("page_id")
	if len(np) == 0 {
		return done, 0, fmt.Errorf("%v: failed to find page_id parameter in %v: %#v", p.NextPage, u.String(), p)
	}
	npi, err := strconv.Atoi(np)
	if err != nil {
		return done, 0, fmt.Errorf("failed to parse %q: %v", np, err)
	}
	cp.Pages.From = npi
	cp.Pagination = p
	return done, npi, err
}

func (cp *checkpoint) filename() string {
	return fmt.Sprintf("checkpoint_%05v_%05v", cp.CurrentPage, cp.TotalPages)
}
