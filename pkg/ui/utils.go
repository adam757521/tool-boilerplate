package ui

import (
	"fmt"
	"runtime"
)

func OptionEnd() {
	if runtime.GOOS == "windows" {
		_, _ = fmt.Scanln()
	}
}

func OptionStart(text string) {
	print(text)

	OptionEnd()
}

func ConvertOptionsToRenderables(options []Option) []Renderable {
	renderables := make([]Renderable, len(options))
	for index, v := range options {
		renderables[index] = v.Renderable
	}
	return renderables
}

func ConvertCategoriesToRenderables(categories []Category) []Renderable {
	renderables := make([]Renderable, len(categories))
	for index, category := range categories {
		renderables[index] = category.Renderable
	}
	return renderables
}
