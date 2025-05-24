/*
Copyright 2024 4rcadia

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mcdr

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"lucy/fscache"
	"lucy/logger"
	"lucy/tools"
	"net/http"
)

// Everything in this context is an extremely long json file of plugins info
// used by the official mcdr catalogue website.

var getEverything = tools.MemoizeE(fetchEverything)

const gzFilename = "everything.json.gz"

const (
	everythingAPIEndpoint = "https://raw.githubusercontent.com/MCDReforged/PluginCatalogue/meta/everything.json.gz"
)

func fetchEverything() (everything *everything, err error) {
	if hit, file, err := fscache.Network.Get(everythingAPIEndpoint); hit && err == nil {
		data, err := io.ReadAll(file)
		if err == nil {
			everything, err = readEverythingGz(data)
			if err == nil {
				return everything, nil
			}
			logger.Warn(fmt.Errorf("failed to read cache: %w", err))
		}
	}

	resp, err := http.Get(everythingAPIEndpoint)
	if err != nil {
		return nil, err
	}
	defer tools.CloseReader(resp.Body, logger.Warn)

	gz, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = fscache.Network.Add(gz, gzFilename, everythingAPIEndpoint, 0)
	if err != nil {
		logger.Warn(err)
	}
	return readEverythingGz(gz)
}

func readEverythingGz(gz []byte) (e *everything, err error) {
	e = &everything{}
	gzReader, err := gzip.NewReader(bytes.NewBuffer(gz))
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	defer tools.CloseReader(gzReader, logger.Warn)
	data, err := io.ReadAll(gzReader)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, e)
	if err != nil {
		return nil, err
	}
	return e, nil
}
