package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"lucy/logger"
	"lucy/tools"
)

func GetFileFromGitHub(apiEndpoint string) (
	err error,
	msg *GhApiMessage,
	data []byte,
) {
	resp, err := http.Get(apiEndpoint)
	if err != nil {
		return err, nil, nil
	}
	defer tools.CloseReader(resp.Body, logger.Warn)
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return err, nil, nil
	}

	// Check if the response is an error message from GitHub API
	err = json.Unmarshal(data, &msg)
	if err == nil && msg != nil && msg.Message != "" {
		return nil, msg, data
	}

	var item GhItem
	err = json.Unmarshal(data, &item)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotDecode, err), nil, nil
	}
	resp, err = http.Get(item.DownloadUrl)
	if err != nil {
		return err, nil, nil
	}
	defer tools.CloseReader(resp.Body, logger.Warn)
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return err, nil, nil
	}

	return nil, nil, data
}
