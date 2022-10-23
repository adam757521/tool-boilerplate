package ui

import (
	"fmt"
	"golang.org/x/sys/windows"
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
	hOut, err := windows.GetStdHandle(windows.STD_OUTPUT_HANDLE)
	if err != nil {
		return err
	}

	fmt.Println("\033[?25l")

	return windows.SetConsoleCursorPosition(hOut, windows.Coord{
		X: 0,
		Y: 0,
	})
}

func Clear() {
	arguments := clear[runtime.GOOS]
	cmd := exec.Command(arguments[0], arguments[1:]...)
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}
