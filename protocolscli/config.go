// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cosnicolaou/protocolsio/api"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Auth struct {
		PublicToken  string `yaml:"public_token"`
		ClientID     string `yaml:"public_clientid"`
		ClientSecret string `yaml:"public_secret"`
	}
	Cache struct {
		Path string `yaml:"path"`
	}
	Endpoints struct {
		ListProtocolsV3 string `yaml:"list_protocols_v3"`
	}
}

func (c *Config) String() string {
	var out strings.Builder
	out.WriteString("auth:\n")
	if len(c.Auth.PublicToken) > 0 {
		fmt.Fprintf(&out, "  token: **redacted**\n")
	}
	return out.String()
}

func (c *Config) WithAuth(ctx context.Context) context.Context {
	return api.WithPublicToken(ctx, c.Auth.PublicToken)
}

func ParseConfig(file string) (*Config, error) {
	cfg := &Config{}
	data, err := os.ReadFile(file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("warning: %q not found\n", file)
			return cfg, nil
		}
		return nil, err
	}
	if err := yaml.Unmarshal([]byte(data), cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
