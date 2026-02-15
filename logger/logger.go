// Package logger provides structured logging with clear separation between
// log-file entries (operational diagnostics) and user-facing messages
// (displayed on stderr).
//
// # Function sets
//
// There are three tiers of logging functions plus a fatal shortcut:
//
//	File-only      Info  Warn  Error  Debug   → written to log file; echoed on console only in verbose mode
//	User-display   ShowInfo  ShowWarn  ShowError   → printed to stderr for the user; NOT persisted to log file
//	Both           ReportInfo  ReportWarn  ReportError → written to log file AND printed to stderr
//	Fatal          Fatal   → logged + displayed + os.Exit(1)
//
// All writes to the log file are synchronous (no queue). A history buffer
// records every file-written entry so that [DumpHistory] can replay them to
// the console at program exit for post-mortem inspection.
package logger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"lucy/tools"
)

// Level represents the severity of a log entry.
// Levels are ordered from least to most severe: Debug < Info < Warn < Error < Fatal.
type Level uint8

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// levelColor maps each level to a styling function.
var levelColor = map[Level]func(any) string{
	LevelDebug: tools.Green,
	LevelInfo:  tools.Cyan,
	LevelWarn:  tools.Yellow,
	LevelError: tools.Red,
	LevelFatal: tools.Red,
}

// prefix returns the bracketed level tag, optionally colored.
func (l Level) prefix(colored bool) string {
	if colored {
		return "[" + levelColor[l](l.String()) + "]"
	}
	return "[" + l.String() + "]"
}

// Entry represents a single log item with its timestamp, level, and content. This is used internally for recording history and is not exposed to users of the logger package.
type entry struct {
	Time    time.Time
	Level   Level
	Content any
}

var (
	verbose bool // when true, file-only entries are also printed to console
	debug   bool // when true, Debug() entries are recorded

	mu      sync.Mutex // write lock for history
	history []*entry
)

// VerboseLevel controls which levels are echoed to the console in verbose
// mode. Everything at or above this level is shown. Set to LevelDebug so
// that all entries are visible.
const VerboseLevel = LevelDebug

// SetVerbose enables echoing of file-only log entries to the console.
func SetVerbose() { verbose = true }

// SetDebug enables Debug-level logging and implies verbose.
func SetDebug() {
	debug = true
	verbose = true
}

// IO

func writeToFile(e *entry) {
	timestamp := e.Time.Format("2006-01-02 15:04:05")
	_, _ = fmt.Fprintln(LogFile, timestamp, e.Level.prefix(false), e.Content)
}

func writeToConsole(e *entry) {
	_, _ = fmt.Fprintln(os.Stderr, e.Level.prefix(true), e.Content)
}

func record(e *entry) {
	mu.Lock()
	history = append(history, e)
	mu.Unlock()
}

// Info logs an informational entry to the log file.
// In verbose mode the entry is also printed to the console.
func Info(content any) {
	e := &entry{Time: time.Now(), Level: LevelInfo, Content: content}
	record(e)
	writeToFile(e)
	if verbose && LevelInfo >= VerboseLevel {
		writeToConsole(e)
	}
}

// Warn logs a warning to the log file.
// In verbose mode the entry is also printed to the console.
func Warn(content error) {
	if content == nil {
		return
	}
	e := &entry{Time: time.Now(), Level: LevelWarn, Content: content}
	record(e)
	writeToFile(e)
	if verbose && LevelWarn >= VerboseLevel {
		writeToConsole(e)
	}
}

// Error logs an error to the log file.
// In verbose mode the entry is also printed to the console.
func Error(content error) {
	if content == nil {
		return
	}
	e := &entry{Time: time.Now(), Level: LevelError, Content: content}
	record(e)
	writeToFile(e)
	if verbose && LevelError >= VerboseLevel {
		writeToConsole(e)
	}
}

// Debug logs a debug entry to the log file. No-op unless debug mode is on.
// In verbose mode (implied by debug) the entry is also printed to the console.
func Debug(content any) {
	if !debug {
		return
	}
	e := &entry{Time: time.Now(), Level: LevelDebug, Content: content}
	record(e)
	writeToFile(e)
	if verbose {
		writeToConsole(e)
	}
}

// ---------------------------------------------------------------------------
// User-display only (console, no file)
// ---------------------------------------------------------------------------

// ShowInfo displays an informational message to the user on stderr.
// The message is NOT written to the log file.
func ShowInfo(content any) {
	writeToConsole(&entry{Time: time.Now(), Level: LevelInfo, Content: content})
}

// ShowWarn displays a warning to the user on stderr.
// The message is NOT written to the log file.
func ShowWarn(content error) {
	writeToConsole(&entry{Time: time.Now(), Level: LevelWarn, Content: content})
}

// ShowError displays an error to the user on stderr.
// The message is NOT written to the log file.
func ShowError(content error) {
	writeToConsole(&entry{Time: time.Now(), Level: LevelError, Content: content})
}

// ---------------------------------------------------------------------------
// Both file + user-display
// ---------------------------------------------------------------------------

// ReportInfo logs an informational message to the file AND displays it to
// the user on stderr.
func ReportInfo(content any) {
	e := &entry{Time: time.Now(), Level: LevelInfo, Content: content}
	record(e)
	writeToFile(e)
	writeToConsole(e)
}

// ReportWarn logs a warning to the file AND displays it to the user on
// stderr.
func ReportWarn(content error) {
	if content == nil {
		return
	}
	e := &entry{Time: time.Now(), Level: LevelWarn, Content: content}
	record(e)
	writeToFile(e)
	writeToConsole(e)
}

// ReportError logs an error to the file AND displays it to the user on
// stderr.
func ReportError(content error) {
	if content == nil {
		return
	}
	e := &entry{Time: time.Now(), Level: LevelError, Content: content}
	record(e)
	writeToFile(e)
	writeToConsole(e)
}

// ---------------------------------------------------------------------------
// Fatal
// ---------------------------------------------------------------------------

// Fatal logs a fatal error to the file, displays it to the user, then
// calls os.Exit(1). Pending history is dumped before exit.
func Fatal(content error) {
	e := &entry{Time: time.Now(), Level: LevelFatal, Content: content}
	record(e)
	writeToFile(e)
	writeToConsole(e)
	DumpHistory()
	os.Exit(1)
}

// ---------------------------------------------------------------------------
// History
// ---------------------------------------------------------------------------

// DumpHistory replays all recorded log entries to the console. This is
// intended to be called from a deferred function in main for post-mortem
// inspection in verbose/debug mode.
//
// Entries already shown via verbose mode will appear again — this is
// intentional so that the dump provides a complete, uninterrupted
// chronological view.
func DumpHistory() {
	if !verbose || len(history) == 0 {
		return
	}
	_, _ = fmt.Fprintln(os.Stderr)
	_, _ = fmt.Fprintln(os.Stderr, tools.Dim("── Log history ("+LogFile.Name()+") ──"))
	for _, e := range history {
		timestamp := tools.Dim(e.Time.Format("15:04:05"))
		_, _ = fmt.Fprintln(os.Stderr, timestamp, e.Level.prefix(true), e.Content)
	}
}
