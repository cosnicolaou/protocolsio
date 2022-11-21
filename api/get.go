// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

var ErrTooManyRequests = errors.New("too many requests")

func Get[T any](ctx context.Context, url string) (T, []byte, error) {
	var m T
	r, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return m, nil, err
	}
	if err := addAuthHeader(ctx, r); err != nil {
		return m, nil, err
	}
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return m, nil, err
	}
	if res.StatusCode == http.StatusTooManyRequests {
		return m, nil, ErrTooManyRequests
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return m, body, err
	}
	return parseJSON[T](body)
}

func parseJSON[T any](s []byte) (T, []byte, error) {
	var r T
	if err := json.Unmarshal(s, &r); err != nil {
		return r, s, err
	}
	return r, s, nil
}

func (p Pagination) Done() bool {
	return p.CurrentPage == p.TotalPages
}

func (p Pagination) PageInfo() (next, total int, done bool, err error) {
	if p.Done() {
		return 0, int(p.TotalPages), true, nil
	}
	u, err := url.Parse(p.NextPage)
	if err != nil {
		return 0, int(p.TotalPages), false, err
	}
	np := u.Query().Get("page_id")
	if len(np) == 0 {
		return 0, int(p.TotalPages), false, fmt.Errorf("%v: failed to find page_id parameter in %v: %#v", p.NextPage, u.String(), p)
	}
	npi, err := strconv.Atoi(np)
	if err != nil {
		return 0, int(p.TotalPages), false, fmt.Errorf("failed to parse %q: %v", np, err)
	}
	return npi, int(p.TotalPages), false, nil
}
