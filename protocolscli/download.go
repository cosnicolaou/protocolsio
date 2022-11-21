// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"cloudeng.io/cmdutil/flags"
	"cloudeng.io/errors"
	"github.com/cosnicolaou/protocolsio/api"
)

type checkpoint struct {
	CurrentPage int
	TotalPages  int
	Pages       flags.IntRangeSpec
	PageSize    int
	Filter      string
	FieldOrder  string
	Order       string
	Total       int
	Pagination  *api.Pagination
	Files       []string
}

func newCheckpointFromFlags(fv *ProtocolsListFlags) (*checkpoint, error) {
	errs := errors.M{}
	errs.Append(flags.OneOf(fv.Filter).Validate("public", ProtocolListFilters()...))
	errs.Append(flags.OneOf(fv.Order).Validate("activity", ProtocolListOrderField()...))
	errs.Append(flags.OneOf(fv.Sort).Validate("asc", "desc"))
	if err := errs.Err(); err != nil {
		return nil, err
	}
	return &checkpoint{
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

func (cp *checkpoint) update(p *api.Pagination) (bool, int, error) {
	done := p.Done()
	cp.CurrentPage = int(p.CurrentPage)
	cp.TotalPages = int(p.TotalPages)
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

type itemSaver struct {
	root       string
	totalItems int
}

func newItemSaver(dir string) (protocolItemProcessor, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	return &itemSaver{dir, 0}, nil
}

func (is *itemSaver) encodeAndWrite(enc *json.Encoder, buf *bytes.Buffer, item any, filename string) error {
	buf.Reset()
	file := filepath.Join(is.root, filename)
	if err := enc.Encode(item); err != nil {
		fmt.Printf("%s: encode error: %v\n", file, err)
		return err
	}
	err := os.WriteFile(file, buf.Bytes(), 0600)
	if err != nil {
		fmt.Printf("%s: write error: %v\n", file, err)
		return err
	}
	fmt.Printf("%s (%v)\n", file, is.totalItems)
	return nil
}

func (is *itemSaver) Process(items []api.Item, cp *checkpoint) error {
	cp.resetFiles()
	errs := errors.M{}
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	for _, item := range items {
		is.totalItems++
		file := fmt.Sprintf("%06d", item.ID)
		cp.appendFile(file)
		if err := is.encodeAndWrite(enc, buf, item, file); err != nil {
			errs.Append(err)
			continue
		}
	}
	if err := errs.Err(); err != nil {
		return err
	}
	// only write the checkpoint if every completed successfully.
	return is.encodeAndWrite(enc, buf, cp, cp.filename())
}
