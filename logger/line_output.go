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
