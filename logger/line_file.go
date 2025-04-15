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

package logger

import (
	"os"
	"path/filepath"
	"runtime"
)

var LogFile = logFile()

func logDir() string {
	var logDir string

	switch runtime.GOOS {
	case "windows":
		logDir = filepath.Join(os.Getenv("APPDATA"), "lucy", "logs")
	case "darwin":
		logDir = filepath.Join(os.Getenv("HOME"), "Library", "Logs", "lucy")
	case "linux":
		logDir = filepath.Join(
			os.Getenv("HOME"),
			".local",
			"share",
			"lucy",
			"logs",
		)
	default:
		logDir = "./logs"
	}

	return logDir
}

func logFile() *os.File {
	logDir := logDir()
	err := os.MkdirAll(logDir, 0o755)
	if err != nil {
		devNull, _ := os.Open(os.DevNull)
		return devNull
	}

	logFilePath := filepath.Join(logDir, "lucy.log")
	logFile, err := os.OpenFile(
		logFilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o755,
	)
	if err != nil {
		println(err.Error())
		devNull, _ := os.Open(os.DevNull)
		return devNull
	}

	return logFile
}
