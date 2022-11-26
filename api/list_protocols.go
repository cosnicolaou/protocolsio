// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package api

import "encoding/json"

type ListProtocolsV3 struct {
	Extras       json.RawMessage
	Items        []json.RawMessage `json:"items"`
	Pagination   Pagination        `json:"pagination"`
	Total        int64             `json:"total"`
	TotalPages   int64             `json:"total_pages"`
	TotalResults int64             `json:"total_results"`
}

type Payload struct {
	Payload    json.RawMessage `json:"payload"`
	StatusCode int             `json:"status_code"`
}

type Protocol struct {
	ID          int64  `json:"id"`
	URI         string `json:"uri"`
	Title       string `json:"title"`
	Description string `json:"description"`
	VersionID   int    `json:"version_id"`
}
