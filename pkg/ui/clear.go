//go:build !windows

package ui

import (
	"errors"
	"os"
	"os/exec"
)

func WindowsResetCursor() error {
	return errors.New("windows only")
}

func Clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}
