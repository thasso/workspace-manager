package ui

import (
	"strings"
	"testing"
)

func TestColorDisabled(t *testing.T) {
	SetColor(false)
	defer SetColor(true)

	if strings.Contains(Bold("test"), "\033") {
		t.Error("Bold should not contain escape codes when color is disabled")
	}
	if strings.Contains(Green("test"), "\033") {
		t.Error("Green should not contain escape codes when color is disabled")
	}
	if strings.Contains(Red("test"), "\033") {
		t.Error("Red should not contain escape codes when color is disabled")
	}
}

func TestColorEnabled(t *testing.T) {
	SetColor(true)

	if !strings.Contains(Bold("test"), "\033[1m") {
		t.Error("Bold should contain escape codes when color is enabled")
	}
	if !strings.Contains(Green("test"), "\033[32m") {
		t.Error("Green should contain escape codes when color is enabled")
	}
}

func TestPlainTextPreserved(t *testing.T) {
	SetColor(false)
	defer SetColor(true)

	if Bold("hello") != "hello" {
		t.Errorf("Bold without color should return original text")
	}
}
