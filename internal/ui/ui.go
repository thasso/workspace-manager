// Package ui provides colored terminal output helpers.
package ui

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

var colorEnabled bool

func init() {
	colorEnabled = term.IsTerminal(int(os.Stdout.Fd())) && os.Getenv("NO_COLOR") == ""
}

// SetColor forces color on or off (useful for testing).
func SetColor(enabled bool) {
	colorEnabled = enabled
}

func wrap(code, s string) string {
	if !colorEnabled {
		return s
	}
	return code + s + "\033[0m"
}

func Bold(s string) string   { return wrap("\033[1m", s) }
func Green(s string) string  { return wrap("\033[32m", s) }
func Red(s string) string    { return wrap("\033[31m", s) }
func Yellow(s string) string { return wrap("\033[33m", s) }
func Cyan(s string) string   { return wrap("\033[36m", s) }

func Success(msg string) { fmt.Printf("  %s %s\n", Green("✓"), msg) }
func Fail(msg string)    { fmt.Printf("  %s %s\n", Red("✗"), msg) }
func Spin(msg string)    { fmt.Printf("  %s %s\n", Cyan("⟳"), msg) }
