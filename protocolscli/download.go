// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"cloudeng.io/errors"
	"github.com/cosnicolaou/protocolsio/api"
)

type ProtocolsDownloadFlags struct {
	ProtocolsListFlags
	CacheDir       string `subcmd:"cachepath,,'location of cache of download protocol objects that overides that specified in the global yaml config'"`
	CheckpointFile string `subcmd:"resume,,checkpoint file to resume download from"`
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
	var cp checkpoint
	if len(fv.CheckpointFile) != 0 {
		data, err := os.ReadFile(fv.CheckpointFile)
		if err != nil {
			return fmt.Errorf("failed to read checkpoint: %v", err)
		}
		if err := json.Unmarshal(data, &cp); err != nil {
			return fmt.Errorf("failed to decode checkpoint file: %v: %v", fv.CheckpointFile, err)
		}
	} else {
		cp, err = newCheckpointFromFlags(&fv.ProtocolsListFlags)
		if err != nil {
			return err
		}
	}
	return getProtocols(ctx, cp, saver)
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
	return is.write(buf.Bytes(), filename)
}

func (is *itemSaver) write(buf []byte, filename string) error {
	file := filepath.Join(is.root, filename)
	err := os.WriteFile(file, buf, 0600)
	if err != nil {
		fmt.Printf("%s: write error: %v\n", file, err)
		return err
	}
	fmt.Printf("%s (%v)\n", file, is.totalItems)
	return nil
}

func (is *itemSaver) fileVersion(filename string) (int, bool, error) {
	file := filepath.Join(is.root, filename)
	buf, err := os.ReadFile(file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0, false, nil
		}
		fmt.Printf("%s: read error: %v\n", file, err)
		return 0, false, err
	}
	protocol, err := api.ParsePayload[api.Protocol](buf)
	if err != nil {
		fmt.Printf("%s: decode error: %v\n", file, err)
		return 0, true, err
	}
	return protocol.VersionID, true, nil
}

func (is *itemSaver) Process(ctx context.Context, protocols api.ListProtocolsV3, cp checkpoint) error {
	cp.resetFiles()
	errs := errors.M{}
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	for _, item := range protocols.Items {
		is.totalItems++
		var p api.Protocol
		if err := json.Unmarshal(item, &p); err != nil {
			errs.Append(err)
			continue
		}
		filebase := fmt.Sprintf("%06d", p.ID)
		file := filebase + ".list"
		cp.appendFile(file)
		tmp := struct {
			Extras json.RawMessage
			Item   json.RawMessage
		}{
			Extras: protocols.Extras,
			Item:   item,
		}
		if err := is.encodeAndWrite(enc, buf, tmp, file); err != nil {
			errs.Append(err)
			continue
		}

		// Fetch the protocol if it has not already been downloaded
		// or there's a newer version.
		file = filebase + ".detail"
		version, exists, err := is.fileVersion(file)
		if err != nil {
			errs.Append(err)
			continue
		}
		if exists && version >= p.VersionID {
			fmt.Printf("%v: [current] (%v >= %v)\n", file, version, p.VersionID)
			continue
		}
		fmt.Printf("%v: [new] (%v, %v < %v)\n", file, exists, version, p.VersionID)
		// Issue a get for this individual protocol since the
		// protocol struct returned by List is incomplete, in particular
		// it does not contain the description field.
		_, body, err := getProtocol(ctx, strconv.Itoa(int(p.ID)))
		if err != nil {
			errs.Append(err)
			continue
		}
		if err := is.write(body, file); err != nil {
			errs.Append(err)
			continue
		}
	}
	if err := errs.Err(); err != nil {
		return err
	}
	// only write the checkpoint if every download operation completed successfully.
	return is.encodeAndWrite(enc, buf, cp, cp.filename())
}
