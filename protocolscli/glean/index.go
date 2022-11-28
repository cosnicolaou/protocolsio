// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package glean

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cosnicolaou/glean/gleancli/config"
	"github.com/cosnicolaou/gleansdk"
	"github.com/cosnicolaou/protocolsio/api"
)

type BulkIndexFlags struct {
	config.ConfigFlags
	UploadID      string `subcmd:"upload-id,upload,id to use for this bulk upload"`
	ForceRestart  bool   `subcmd:"force-restart,false,restart the bulk upload"`
	ForceDeletion bool   `subcmd:"force-sync-deletion,false,synchronously delete stale documents on upload of last bulk indexing batch"`
}

func bulkIndexCmd(ctx context.Context, values interface{}, args []string) error {
	fv := values.(*BulkIndexFlags)
	cfg, err := config.ParseConfig(fv.Config)
	if err != nil {
		return err
	}
	ctx, client := cfg.NewAPIClient(ctx)

	dl, err := newDirLister(args[0])
	if err != nil {
		return err
	}

	ch := make(chan dirListResult, 100)
	go func() {
		dl.stream(ctx, ch)
		close(ch)
	}()

	var (
		firstPage = true
		indexed   = 0
		duration  time.Duration
	)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case lr, ok := <-ch:
			if !ok {
				return nil
			}
			if err := lr.err; err != nil {
				return err
			}
			gd := bulkIndexReq(lr.protocols)
			gd.SetIsFirstPage(firstPage)
			gd.SetIsLastPage(lr.lastPage)
			if firstPage {
				gd.SetForceRestartUpload(fv.ForceRestart)
				firstPage = false
			}
			gd.SetDatasource(DatasourceName)
			gd.SetUploadId(fv.UploadID)
			reqStart := time.Now()
			if err := executeBulkIndexRequest(ctx, client, gd); err != nil {
				return err
			}
			took := time.Since(reqStart)
			duration += took
			indexed += len(gd.Documents)
			avg := time.Duration(int64(duration) / int64(indexed))
			fmt.Printf("indexed: total # docs: % 5v, per req # docs: % 3v in % 8v (avg: %8v)\n", indexed, len(gd.Documents), took, avg)
			if lr.lastPage {
				fmt.Printf("indexed: all # docs: % 5v docs in % 8v, (avg: %8v)\n", indexed, duration, avg)
			}
		}
	}
}

func executeBulkIndexRequest(ctx context.Context, client *gleansdk.APIClient, gd gleansdk.BulkIndexDocumentsRequest) error {
	resp, err := client.DocumentsApi.BulkindexdocumentsPost(ctx).BulkIndexDocumentsRequest(gd).Execute()
	if err != nil {
		fmt.Printf("response: %v\n", resp)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected http status code: %v", resp.Status)
	}
	return nil
}

func bulkIndexReq(protocols []*api.Protocol) gleansdk.BulkIndexDocumentsRequest {
	var req gleansdk.BulkIndexDocumentsRequest
	req.Datasource = DatasourceName
	for _, p := range protocols {
		gd := gleanDocument(p)
		req.Documents = append(req.Documents, gd)
	}
	return req
}

func gleanDocument(p *api.Protocol) gleansdk.DocumentDefinition {
	gd := gleansdk.DocumentDefinition{}
	gd.Datasource = DatasourceName
	gd.SetId(p.URI)
	gd.SetViewURL(p.URL)
	gd.SetTitle(p.Title)
	gd.Summary = &gleansdk.ContentDefinition{}
	gd.Summary.SetMimeType("text/plain")
	var tmpDesc struct {
		Blocks []struct {
			Text string
			Key  string
		}
	}
	if err := json.Unmarshal([]byte(p.Description), &tmpDesc); err == nil && len(tmpDesc.Blocks) > 0 {
		gd.Summary.SetTextContent(tmpDesc.Blocks[0].Text)
	} else {
		gd.Summary.SetTextContent(p.Description)
	}

	gd.Author = &gleansdk.UserReferenceDefinition{}
	gd.Author.SetName(p.Creator.Name)
	gd.Author.SetEmail(p.Creator.Username + "@protocols.io")
	gd.Author.SetDatasourceUserId(p.Creator.Username)
	gd.Permissions = &gleansdk.DocumentPermissionsDefinition{}
	gd.Permissions.SetAllowAnonymousAccess(true)
	gd.CreatedAt = new(int64)
	*gd.CreatedAt = int64(p.CreatedOn)
	return gd
}

func newDirLister(dir string) (*dirLister, error) {
	d, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	return &dirLister{dir: dir, df: d}, nil
}

type dirLister struct {
	dir string
	df  *os.File
}

type dirListResult struct {
	protocols []*api.Protocol
	err       error
	lastPage  bool
}

func (dl dirLister) stream(ctx context.Context, ch chan<- dirListResult) {
	numEntries := 50
	for {
		select {
		case <-ctx.Done():
			ch <- dirListResult{err: ctx.Err()}
		default:
		}
		de, err := dl.df.ReadDir(numEntries)
		if err == io.EOF {
			return
		}
		if err != nil {
			ch <- dirListResult{err: err}
			return
		}
		dl := dl.readFiles(de)
		dl.lastPage = len(de) != numEntries
		ch <- dl
	}
}

func (dl dirLister) readFiles(entries []fs.DirEntry) dirListResult {
	var lr dirListResult
	for _, f := range entries {
		if !strings.HasSuffix(f.Name(), ".detail") {
			continue
		}
		buf, err := os.ReadFile(filepath.Join(dl.dir, f.Name()))
		if err != nil {
			return dirListResult{err: err}
		}
		p, err := api.ParsePayload[api.Protocol](buf)
		if err != nil {
			return dirListResult{err: err}
		}
		lr.protocols = append(lr.protocols, &p)
	}
	return lr
}
