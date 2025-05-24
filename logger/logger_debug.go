//go:build debug

package logger

func init() {
	UseConsoleOutput()
	UseDebug()
}
