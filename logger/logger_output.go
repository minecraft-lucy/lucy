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
	"fmt"
	"os"
	"time"

	"lucy/tools"
)

func pop() {
	item, _ := queue.Get(0)
	writeItem(item.(*logItem))
	queue.Remove(0)
}

var color = map[logLevel]func(any) string{
	lInfo:  tools.Cyan,
	lWarn:  tools.Yellow,
	lError: tools.Red,
	lFatal: tools.Red,
	lDebug: tools.Green,
}

func (level logLevel) prefix(colored bool) string {
	return "[" + tools.Ternary(
		colored,
		color[level](level.String()),
		level.String(),
	) + "]"
}

func writeItem(message *logItem) {
	if toConsole {
		writeToConsole(message)
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	_, _ = fmt.Fprintln(
		LogFile,
		timestamp,
		message.Level.prefix(false),
		message.Content,
	)
}

func writeToConsole(message *logItem) {
	_, _ = fmt.Fprintln(
		os.Stderr,
		message.Level.prefix(true),
		message.Content,
	)
}
