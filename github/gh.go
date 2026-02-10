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
	res, err := http.Get(apiEndpoint)
	if err != nil {
		return err, nil, nil
	}
	defer tools.CloseReader(res.Body, logger.Warn)
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return err, nil, nil
	}

	err = json.Unmarshal(data, &msg)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotDecode, err), nil, nil
	}
	if msg != nil && msg.Message != "" {
		return nil, msg, data
	}
	return nil, nil, data
}
