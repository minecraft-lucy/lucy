package mcdr

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"lucy/logger"
	"lucy/tools"
	"lucy/util"
)

// Everything in this context is an extremely long json file of plugins info
// used by the official mcdr catalogue website.

var getEverything = tools.MemoizeE(fetchEverything)

const (
	everythingAPIEndpoint = "https://raw.githubusercontent.com/MCDReforged/PluginCatalogue/meta/everything.json.gz"
	cacheExpiration       = 24 * time.Hour
)

func fetchEverything() (everything *everything, err error) {
	if exist, err := checkEverythingCache(); err != nil && exist {
		everything, err = getEverythingCache()
		if err == nil {
			return everything, nil
		}
		logger.Warn(fmt.Errorf("failed to read cache: %w", err))
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

	err = cacheEverythingGz(gz)
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

func cacheEverythingGz(gz []byte) error {
	if !util.DoCache {
		return nil
	}
	filepath := path.Join(util.CacheDir, "everything.json.gz")
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	_, err = file.Write(gz)
	if err != nil {
		return err
	}
	return nil
}

func checkEverythingCache() (bool, error) {
	if !util.DoCache {
		return false, nil
	}
	filepath := path.Join(util.CacheDir, "everything.json.gz")
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	file, err := os.Open(filepath)
	if err != nil {
		return false, err
	}
	gz, err := io.ReadAll(file)
	if err != nil {
		return false, err
	}
	everything, err := readEverythingGz(gz)
	if err != nil {
		return false, err
	}
	if time.Since(time.Unix(int64(everything.Timestamp), 0)) > cacheExpiration {
		return false, nil
	}
	return true, nil
}

func getEverythingCache() (*everything, error) {
	if !util.DoCache {
		return nil, fmt.Errorf("cache not available")
	}
	filepath := path.Join(util.CacheDir, "everything.json.gz")
	if exist, err := checkEverythingCache(); exist && err == nil {
		var gz []byte
		gz, err = os.ReadFile(filepath)
		return readEverythingGz(gz)
	}
	return nil, fmt.Errorf("cache not found")
}
