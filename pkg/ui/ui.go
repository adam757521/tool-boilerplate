package ui

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
)

type ColorInfo struct {
	Ascii  *color.Color
	Error  *color.Color
	Parent *color.Color
	Child  *color.Color
}

type RenderInfo struct {
	Prefix string
	Color  *color.Color
}

type Renderable struct {
	Label string
}

type Option struct {
	Renderable
	Callback func()
}

type Category struct {
	Renderable
	Options []Option
}

type UI struct {
	Title      string
	Ascii      string
	Color      ColorInfo
	Categories []Category
	Version    string
}

func (r *Renderable) Render(info *RenderInfo) {
	var printf func(f string, a ...interface{}) (n int, err error)
	if info.Color == nil {
		printf = fmt.Printf
	} else {
		printf = info.Color.Printf
	}

	_, _ = printf("%s%s\n", info.Prefix, r.Label)
}

func (ui *UI) WaitForInput(options int) (int, bool) {
	var input int
	_, _ = fmt.Scan(&input)

	index := input - 1
	if index < 0 || index >= options+1 {
		return -1, false
	} else {
		return index, true
	}
}

func (ui *UI) RenderInformation() {
	_, _ = ui.Color.Parent.Print("\nVersion: ")
	_, _ = ui.Color.Child.Print(ui.Version + "\n")
}

func (ui *UI) DefaultRender(err string, parent string, renderables []Renderable) {
	Clear()

	_, _ = ui.Color.Ascii.Println(ui.Ascii)

	if err != "" {
		_, _ = ui.Color.Error.Println(err)
	}

	ui.RenderInformation()

	_, _ = ui.Color.Parent.Println("\n" + parent)
	for i, r := range renderables {
		r.Render(&RenderInfo{
			Prefix: fmt.Sprintf("  [%d] ", i+1),
			Color:  ui.Color.Child,
		})
	}
}

func (ui *UI) ChangeTitle(title string) {
	cmd := exec.Command("title", title)
	_ = cmd.Run()
}

func (category *Category) Render(parentUI *UI, err string) {
	parentUI.DefaultRender(err, "Options", ConvertOptionsToRenderables(category.Options))

	back := Renderable{Label: "Back"}
	back.Render(&RenderInfo{
		Prefix: fmt.Sprintf("  [%d] ", len(category.Options)+1),
		Color:  parentUI.Color.Child,
	})

	_, _ = parentUI.Color.Parent.Print("\nSelect an option > ")

	index, ok := parentUI.WaitForInput(len(category.Options))
	if !ok {
		category.Render(parentUI, "[ERROR] Invalid option selected.")
	} else {
		if index == len(category.Options) {
			parentUI.RenderCategories("")
		} else {
			parentUI.DefaultRender("", "Running option...", []Renderable{})
			option := category.Options[index]

			parentUI.ChangeTitle(fmt.Sprintf("%s - %s", parentUI.Title, option.Renderable.Label))
			option.Callback()

			_, _ = parentUI.Color.Parent.Print("\nOption finished.")
			_, _ = parentUI.Color.Child.Print(" Press enter to return to menu.")

			_, _ = fmt.Scanln()

			category.Render(parentUI, "")
		}
	}
}

func (ui *UI) RenderAuth(authValidator func(string) bool) {
	ui.DefaultRender("", "Enter license >", []Renderable{})

	var license string
	_, _ = fmt.Scan(&license)

	if authValidator(license) {
		ui.RenderCategories("")
	}
}

func (ui *UI) RenderCategories(err string) {
	ui.DefaultRender(err, "Categories", ConvertCategoriesToRenderables(ui.Categories))

	exit := Renderable{Label: "Exit"}
	exit.Render(&RenderInfo{
		Prefix: fmt.Sprintf("  [%d] ", len(ui.Categories)+1),
		Color:  ui.Color.Child,
	})

	_, _ = ui.Color.Parent.Print("\nSelect a category > ")

	index, ok := ui.WaitForInput(len(ui.Categories))
	if !ok {
		ui.RenderCategories("[ERROR] Invalid category selected.")
	} else {
		if index == len(ui.Categories) {
			os.Exit(0)
		} else {
			ui.Categories[index].Render(ui, "")
		}
	}
}
