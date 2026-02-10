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
	if toConsole {
		println("The following is written to " + LogFile.Name() + ":")
	}
	for queue.Empty() == false {
		pop()
	}
}

var queue = singlylinkedlist.New()

func factoryQueuedLog(level logLevel) func(content error) {
	return func(content error) {
		if content == nil {
			return
		}
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
