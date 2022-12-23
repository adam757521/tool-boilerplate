package main

import (
	"fmt"

	"github.com/adam757521/tool-boilerplate/pkg/ui"
)

type Printer struct{}

func (p Printer) GetAscii() string {
	return "hello"
}

func RenderFunc(selected bool) (string, int) {
	if selected {
		return "> ss", 4
	}

	return "ss", 2
}

func Callback() {
	fmt.Println("enter name:")

	var name string
	_, _ = fmt.Scanln(&name)

	fmt.Println(name)
}

func main() {
	fmt.Println(ui.Dash + "➡️" + ui.Dash + "\n " + ui.VerticalDash)

	interactive := ui.New(Printer{}, "Hello", nil, "Twitch AIO", 10, &ui.RGB{255, 0, 0})
	interactive.Fields = []*ui.InteractiveField{
		&ui.InteractiveField{RenderLength: 4, Render: RenderFunc, Title: "Title", TitleLength: 5, Callback: Callback, Description: "Fuck"},
		&ui.InteractiveField{RenderLength: 4, Render: RenderFunc, Title: "Title", TitleLength: 5, Callback: func() { fmt.Println("hello world") }, Description: "Wish i was dead"},
	}

	interactive.Start()
}
