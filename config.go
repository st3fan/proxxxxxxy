// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Hostname       string
	Interface      string
	Prefix         string
	Address        string
	APIPort        int
	ProxyPortStart int
	Allow          []string
	Verbose        bool
}

func newConfigFromFile(path string) (Config, error) {
	encodedConfig, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(encodedConfig, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}
