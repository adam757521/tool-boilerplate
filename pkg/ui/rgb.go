package ui

import (
	"fmt"
	"golang.org/x/term"
	"os"
	"strings"
)

type RGB struct {
	R int
	G int
	B int
}

func (r *RGB) Color(text string) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", r.R, r.G, r.B, text)
}

func (r *RGB) Foreground(text string) {
	fmt.Print(r.Color(text))
}

func PrintCenteredAscii(ascii []string, length int, colors []RGB) {
	var padding int

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err == nil {
		padding = (width / 2) - (length / 2)
	}

	for i, part := range ascii {
		index := i / (len(ascii) / len(colors))
		if index > len(colors)-1 {
			index = len(colors) - 1
		}

		colors[index].Foreground(strings.Repeat(" ", padding) + part + "\n")
	}

	fmt.Println()
}
