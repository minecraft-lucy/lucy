package mcdr

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"

	"lucy/logger"
	"lucy/tools"
	"lucy/util"
)

// Everything in this context is an extremely long json file of plugins info
// used by the official mcdr catalogue website.

var getEverything = tools.MemoizeE(fetchEverything)

const EverythingAPIEndpoint = "https://raw.githubusercontent.com/MCDReforged/PluginCatalogue/meta/everything.json.gz"

func fetchEverything() (everything *everything, err error) {
	if exist, err := checkEverythingCache(); err != nil && exist {
		return getEverythingCache()
	}

	resp, err := http.Get(EverythingAPIEndpoint)
	if err != nil {
		return nil, err
	}
	tools.CloseReader(resp.Body, logger.Warn)

	var gz []byte
	_, err = resp.Body.Read(gz)
	if err != nil {
		return nil, err
	}

	err = cacheEverythingGz(gz)
	if err != nil {
		logger.Warn(err)
	}
	return readEverythingGz(gz)
}

func readEverythingGz(gz []byte) (everything *everything, err error) {
	gzReader, err := gzip.NewReader(bytes.NewBuffer(gz))
	if err != nil {
		return nil, err
	}
	tools.CloseReader(gzReader, logger.Warn)

	var buff []byte
	_, err = gzReader.Read(buff)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buff, everything)
	if err != nil {
		return nil, err
	}
	return everything, nil
}

func cacheEverythingGz(gz []byte) error {
	filepath := path.Join(util.CacheDir, "everything.json.gz")
	_, err := os.Create(filepath)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath, gz, 0o644)
	if err != nil {
		return err
	}
	return nil
}

func checkEverythingCache() (bool, error) {
	filepath := path.Join(util.CacheDir, "everything.json.gz")
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func getEverythingCache() (*everything, error) {
	filepath := path.Join(util.CacheDir, "everything.json.gz")
	if exist, err := checkEverythingCache(); exist && err == nil {
		var gz []byte
		gz, err = os.ReadFile(filepath)
		return readEverythingGz(gz)
	}
	return nil, fmt.Errorf("cache not found")
}
