package main

import (
	"fmt"
	"github.com/adam757521/tool-boilerplate/pkg/ui"
)

func test() {
	defer ui.OptionEnd()

	// When taking input, you must use ui.OptionStart to print information before it. only for the first time. Windows issue.
	fmt.Println("Hello world")
}

func main() {
	myInterface := &ui.UI{
		Color: ui.ColorInfo{
			Error:  &ui.RGB{R: 255},
			Parent: &ui.RGB{B: 255},
			Child:  &ui.RGB{B: 255},
		},
		Categories: nil,
		Version:    "1.0.0",
		Title:      "Example",
	}

	category := ui.Category{Renderable: ui.Renderable{Label: "Test Category"}, Options: []ui.Option{
		{Renderable: ui.Renderable{Label: "Test"}, Callback: test},
	}}
	category.Render(myInterface, "")
}
