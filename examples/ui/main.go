package main

import (
	"fmt"
	"github.com/fatih/color"
	"tool-boilerplate/pkg/ui"
)

func test() {
	defer ui.OptionEnd()

	// When taking input, you must use ui.OptionStart to print information before it. only for the first time. Windows issue.
	fmt.Println("Hello world")
}

func main() {
	myInterface := &ui.UI{
		Ascii: "TestAscii",
		Color: ui.ColorInfo{
			Ascii:  color.New(color.FgHiBlue),
			Error:  color.New(color.FgHiRed),
			Parent: color.New(color.FgHiBlue),
			Child:  color.New(color.FgHiBlue),
		},
		Categories: nil,
		Version:    "1.0.0",
		Author:     "...",
	}

	category := ui.Category{Renderable: ui.Renderable{Label: "Test Category"}, Options: []ui.Option{
		{Renderable: ui.Renderable{Label: "Test"}, Callback: test},
	}}
	category.Render(myInterface, "")
}
