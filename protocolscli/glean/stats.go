// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package glean

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cosnicolaou/glean/gleancli/config"
	"github.com/cosnicolaou/gleansdk"
)

type StatsFlags struct {
	config.ConfigFlags
}

func statsCmd(ctx context.Context, values interface{}, args []string) error {
	fv := values.(*StatsFlags)
	cfg, err := config.ParseConfig(fv.Config)
	if err != nil {
		return err
	}
	ctx, client := cfg.NewAPIClient(ctx)

	req := gleansdk.NewGetDocumentCountRequest(DatasourceName)
	fmt.Printf("%#v\n", req)

	resp, r, err := client.DocumentsApi.GetdocumentcountPost(ctx).GetDocumentCountRequest(*req).Execute()
	out, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return parseError(r, err)

}

func parseError(r *http.Response, err error) error {
	if err == nil {
		return nil
	}
	oapiErr, ok := err.(*gleansdk.GenericOpenAPIError)
	if !ok {
		return fmt.Errorf("%v: %v: %v\n", r.Request.URL, r.StatusCode, err)
	}
	var tmp any
	body := oapiErr.Body()
	if json.Unmarshal(body, &tmp) == nil {
		return fmt.Errorf("%s: %v", body, err)
	}
	if body, nerr := io.ReadAll(r.Body); nerr == nil {
		return fmt.Errorf("%s: %v", body, err)
	}
	return err
}
