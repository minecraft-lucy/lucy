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

type logItem struct {
	Level   logLevel
	Content any
}

type logLevel uint8

const (
	lInfo logLevel = iota
	lWarn
	lError
	lFatal
	lDebug
)

func (level logLevel) String() string {
	switch level {
	case lInfo:
		return "INFO"
	case lWarn:
		return "WARN"
	case lError:
		return "ERROR"
	case lFatal:
		return "FATAL"
	case lDebug:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}
