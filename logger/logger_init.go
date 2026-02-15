package logger

// This file contains initialization and global state for the logger package.

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
