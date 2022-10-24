package main

import (
	"fmt"
	"github.com/adam757521/tool-boilerplate/pkg/ui"
	"github.com/fatih/color"
	"runtime"
	"time"
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
			Header: "AdBot",
			Fields: []*ui.Field{
				{
					Header:      "Threads:",
					HeaderColor: color.BlueString,
					Value:       "3",
				},
				{
					Header:      "Revenue",
					HeaderColor: color.BlueString,
					Value:       "0.000/1.340$",
				},
				{
					Header:      "Ads (Hour):",
					HeaderColor: color.BlueString,
					Value:       "0",
				},
				{
					Header:      "Revenue (Hour):",
					HeaderColor: color.BlueString,
					Value:       "0.000$",
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
