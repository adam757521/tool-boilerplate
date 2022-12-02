package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ColorInfo struct {
	Error  *RGB
	Parent *RGB
	Child  *RGB
}

type RenderInfo struct {
	Prefix string
	Color  *RGB
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
	Ascii      func()
	Color      ColorInfo
	Categories []Category
	Version    string
	Title      string
	Notice     string
}

func (r *Renderable) Render(info *RenderInfo) {
	if info.Color == nil {
		fmt.Printf("%s%s\n", info.Prefix, r.Label)
	} else {
		info.Color.Foreground(fmt.Sprintf("%s%s\n", info.Prefix, r.Label))
	}
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
	term := ""
	if ui.Notice != "" {
		term = "\n"
	}

	ui.Color.Parent.Foreground("\nVersion: ")
	ui.Color.Child.Foreground(ui.Version + term)

	if ui.Notice != "" {
		ui.Color.Error.Foreground("\nNotice: ")
		ui.Color.Child.Foreground(ui.Notice)
	}
}

func (ui *UI) DefaultRender(err string, parent string, renderables []Renderable) {
	Clear()

	if ui.Ascii != nil {
		ui.Ascii()
	}

	if err != "" {
		ui.Color.Error.Foreground(err + "\n")
	}

	ui.RenderInformation()

	ui.Color.Parent.Foreground("\n" + parent + "\n")
	for i, r := range renderables {
		r.Render(&RenderInfo{
			Prefix: fmt.Sprintf("  [%d] ", i+1),
			Color:  ui.Color.Child,
		})
	}
}

func (ui *UI) ChangeTitle(title string) {
	cmd := exec.Command("cmd", append([]string{"/c", "title"}, strings.Split(title, " ")...)...)
	_ = cmd.Run()
}

func (category *Category) Render(parentUI *UI, err string) {
	parentUI.ChangeTitle(parentUI.Title)

	parentUI.DefaultRender(err, "Options", ConvertOptionsToRenderables(category.Options))

	back := Renderable{Label: "Back"}
	back.Render(&RenderInfo{
		Prefix: fmt.Sprintf("  [%d] ", len(category.Options)+1),
		Color:  parentUI.Color.Child,
	})

	parentUI.Color.Parent.Foreground("\nSelect an option > ")

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

			parentUI.Color.Parent.Foreground("\nOption finished.")
			parentUI.Color.Child.Foreground(" Press enter to return to menu.")

			_, _ = fmt.Scanln()

			category.Render(parentUI, "")
		}
	}
}

func (ui *UI) RenderAuth(authValidator func(string) bool) {
	ui.ChangeTitle(ui.Title)

	ui.DefaultRender("", "Enter license >", []Renderable{})

	var license string
	_, _ = fmt.Scan(&license)

	if authValidator(license) {
		ui.RenderCategories("")
	}
}

func (ui *UI) RenderCategories(err string) {
	ui.ChangeTitle(ui.Title)

	ui.DefaultRender(err, "Categories", ConvertCategoriesToRenderables(ui.Categories))

	exit := Renderable{Label: "Exit"}
	exit.Render(&RenderInfo{
		Prefix: fmt.Sprintf("  [%d] ", len(ui.Categories)+1),
		Color:  ui.Color.Child,
	})

	ui.Color.Parent.Foreground("\nSelect a category > ")

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
