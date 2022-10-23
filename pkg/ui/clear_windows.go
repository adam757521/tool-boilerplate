package ui

import (
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"os/exec"
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
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}
