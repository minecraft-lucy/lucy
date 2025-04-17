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

	"github.com/emirpasic/gods/lists/singlylinkedlist"
)

var (
	debug     = false
	toConsole = false
)

func UseDebug() {
	debug = true
}

func UseConsoleOutput() {
	toConsole = true
}

func WriteAll() {
	for queue.Empty() == false {
		pop()
	}
}

var queue = singlylinkedlist.New()

func factoryQueuedLog(level logLevel) func(content error) {
	return func(content error) {
		queue.Add(&logItem{Level: level, Content: content})
	}
}

func factoryInstantLog(level logLevel) func(content error) {
	return func(content error) {
		writeToConsole(&logItem{Level: level, Content: content})
	}
}

var (
	Info = func(content any) {
		queue.Add(&logItem{Level: lInfo, Content: content})
	}
	Warn  = factoryQueuedLog(lWarn)
	Error = factoryQueuedLog(lError)
	Fatal = FatalNow
	Debug = func(content any) {
		if debug {
			queue.Add(&logItem{Level: lDebug, Content: content})
		}
	}
)

// These functions bypass the toConsole flag and write directly to the console.
// Use them to communicate with the user in a way that cannot be suppressed.
var (
	InfoNow = func(content any) {
		writeToConsole(&logItem{Level: lInfo, Content: content})
	}
	WarnNow  = factoryInstantLog(lWarn)
	ErrorNow = factoryInstantLog(lError)
	FatalNow = func(content error) {
		defer os.Exit(1)
		factoryInstantLog(lFatal)(content)
		WriteAll()
	}
)
