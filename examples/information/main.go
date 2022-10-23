package main

import (
	"fmt"
	"github.com/fatih/color"
	"runtime"
	"time"
	"tool-boilerplate/pkg/ui"
)

func CUILoop(lastUpdated *time.Time) {
	if time.Now().Sub(*lastUpdated).Minutes() >= 1 {
		*lastUpdated = time.Now()

		if runtime.GOOS == "windows" {
			// Clear the UI on Windows because ui.WindowsResetCursor is not enough.
			// If you do not print anything in your program, it is not required but recommended.
			ui.Clear()
		}

		// Do some minute calculations, like operations per minute.
	}

	if runtime.GOOS == "windows" {
		// Smoother than calling ui.Clear.
		_ = ui.WindowsResetCursor()
	} else {
		ui.Clear()
	}

	sections := []*ui.Section{
		{
			Header: "Main",
			Fields: []*ui.Field{
				{
					Header:      "Next Clear In:",
					HeaderColor: color.BlueString,
					Value:       "10",
				},
			},
			Color: color.GreenString,
		},
	}

	printable, err := ui.ToPrintable(sections)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(printable)
	}
}

func main() {
	if runtime.GOOS == "windows" {
		// To make unicode characters work.
		ui.Clear()
	}

	lastUpdated := time.Now()
	for {
		CUILoop(&lastUpdated)
	}
}
