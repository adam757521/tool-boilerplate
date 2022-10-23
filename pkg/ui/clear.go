package ui

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
)

var (
	clear = map[string][]string{
		"windows": {"cmd", "/c", "cls"},
		"linux":   {"clear"},
		"darwin":  {"clear"},
	}
)

func WindowsResetCursor() error {
	return errors.New("windows only")
}

func Clear() {
	arguments := clear[runtime.GOOS]
	cmd := exec.Command(arguments[0], arguments[1:]...)
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}
