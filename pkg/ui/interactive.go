package ui

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mitchellh/go-wordwrap"
	"github.com/nsf/termbox-go"
	"golang.org/x/term"
)

const (
	Padding = 3
)

var (
	NotArrowError = errors.New("not an arrow")
)

type ASCIIPrinter interface {
	GetAscii() string
}

type InteractiveField struct {
	Render       func(bool) (string, int)
	RenderLength int
	UI           *InteractiveUI
	Callback     func()
	Title        string
	TitleLength  int
	Description  string
}

type InteractiveUI struct {
	Printer     ASCIIPrinter
	Selected    int
	Fields      []*InteractiveField
	Notice      string
	Back        *InteractiveUI
	Title       string
	RGB         *RGB
	TitleLength int
}

func GetPressedArrow() (termbox.Key, error) {
	ev := termbox.PollEvent()
	if ev.Key == 0 {
		return 0, NotArrowError
	}

	return ev.Key, nil
}

func New(printer ASCIIPrinter, notice string, back *InteractiveUI, title string, titleLength int, rgb *RGB) *InteractiveUI {
	return &InteractiveUI{
		Printer:     printer,
		Notice:      notice,
		Back:        back,
		Title:       title,
		TitleLength: titleLength,
		RGB:         rgb,
	}
}

func (u *InteractiveUI) GetRenderedFields() ([]string, error) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return nil, err
	}

	var selectedField *InteractiveField
	if len(u.Fields) != 0 {
		selectedField = u.Fields[u.Selected]
	}

	fieldSectionWidth := u.TitleLength
	for _, f := range u.Fields {
		if fieldSectionWidth < f.RenderLength {
			fieldSectionWidth = f.RenderLength
		}
	}
	fieldSectionWidth += 2

	// Add the padding and 2 characters
	fieldSectionWidthTotal := fieldSectionWidth + 2 + Padding

	descriptionSectionWidthTotal := width - fieldSectionWidthTotal - Padding
	descriptionSectionWidth := descriptionSectionWidthTotal - 2

	descriptionTitlePadding := (descriptionSectionWidth - selectedField.TitleLength) / 2

	padding := strings.Repeat(" ", Padding)
	neededSpace := ""
	if (descriptionSectionWidth-selectedField.TitleLength)%2 == 1 {
		neededSpace = " "
	}

	fieldTitlePadding := (fieldSectionWidth - u.TitleLength) / 2
	fieldNeededSpace := ""
	if (fieldSectionWidth-u.TitleLength)%2 == 1 {
		fieldNeededSpace = " "
	}

	fieldSectionDash := strings.Repeat(Dash, fieldSectionWidth)
	descriptionSectionDash := strings.Repeat(Dash, descriptionSectionWidth)

	descriptionPadding := strings.Repeat(" ", descriptionTitlePadding)
	fieldPadding := strings.Repeat(" ", fieldTitlePadding)

	lines := []string{}

	lines = append(lines, u.RGB.Color(padding+CornerLeft+fieldSectionDash+IntersectDown+descriptionSectionDash+CornerRight))
	lines = append(lines, padding+u.RGB.Color(VerticalDash)+fieldPadding+u.Title+fieldPadding+fieldNeededSpace+u.RGB.Color(VerticalDash)+descriptionPadding+selectedField.Title+descriptionPadding+neededSpace+u.RGB.Color(VerticalDash))

	lastCharacter := IntersectRight
	if selectedField.Description == "" {
		lastCharacter = CornerDownRight
	}

	lines = append(lines, u.RGB.Color(padding+IntersectLeft+fieldSectionDash+Intersect+descriptionSectionDash+lastCharacter))

	descriptionSplitted := []string{}
	if selectedField.Description != "" {
		descriptionSplitted = strings.Split(wordwrap.WrapString(selectedField.Description, uint(descriptionSectionWidth-2)), "\n")
	}

	length := len(descriptionSplitted)
	if len(u.Fields) > length {
		length = len(u.Fields)
	}

	renderFieldEnd := false
	renderDescriptionEnd := false
	renderedEnd := false

	for i := 0; i < length; i++ {
		var field *InteractiveField
		var currentDescription string

		fieldSelected := i == u.Selected
		fieldUsed := true
		descriptionUsed := true

		if i >= len(u.Fields) {
			fieldUsed = false
		} else {
			field = u.Fields[i]
		}

		if i >= len(descriptionSplitted) {
			descriptionUsed = false
		} else {
			currentDescription = descriptionSplitted[i]
		}

		if i == 0 && !descriptionUsed {
			renderDescriptionEnd = true
			renderedEnd = true
		}

		fieldText := strings.Repeat(" ", fieldSectionWidthTotal-1) + VerticalDash
		descriptionText := ""

		if fieldUsed {
			text, length := field.Render(fieldSelected)
			fieldText = padding + u.RGB.Color(VerticalDash) + " " + text + strings.Repeat(" ", fieldSectionWidth-1-length)

			if !renderDescriptionEnd || renderedEnd {
				fieldText += u.RGB.Color(VerticalDash)
			}
		}

		if !renderedEnd {
			if renderFieldEnd {
				fieldText = padding + u.RGB.Color(CornerDownLeft+fieldSectionDash+IntersectRight)
				renderedEnd = true
			}

			if renderDescriptionEnd {
				descriptionText = u.RGB.Color(IntersectLeft + descriptionSectionDash + CornerDownRight)
				renderedEnd = true
			}
		}

		if descriptionUsed {
			descriptionText = " " + currentDescription + strings.Repeat(" ", descriptionSectionWidth-len(currentDescription)-2) + " " + u.RGB.Color(VerticalDash)
		}

		lines = append(lines, fieldText+descriptionText)

		if !renderedEnd {
			if i == len(descriptionSplitted)-1 {
				renderDescriptionEnd = true
			}

			if i == len(u.Fields)-1 {
				renderFieldEnd = true
			}

		}
	}

	finalLine := padding + CornerDownLeft + fieldSectionDash + IntersectUp + descriptionSectionDash + CornerDownRight

	if renderedEnd {
		if renderFieldEnd {
			finalLine = padding + strings.Repeat(" ", fieldSectionWidthTotal-4) + CornerDownLeft + descriptionSectionDash + CornerDownRight
		} else if renderDescriptionEnd {
			finalLine = padding + CornerDownLeft + fieldSectionDash + CornerDownRight
		}
	}
	lines = append(lines, u.RGB.Color(finalLine))

	return lines, nil
}

func formatDuration(duration time.Duration) string {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	milliseconds := int(duration.Milliseconds()) % 1000

	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	} else if seconds > 0 {
		return fmt.Sprintf("%ds %dms", seconds, milliseconds)
	} else {
		return fmt.Sprintf("%dms", milliseconds)
	}
}

func (u *InteractiveUI) StartNonInitialize() {
	for {
		Clear()

		fmt.Println(u.Printer.GetAscii() + "\n")
		fmt.Println(u.Notice + "\n")

		rendered, err := u.GetRenderedFields()
		if err != nil {
			panic(err)
		}

		for _, line := range rendered {
			fmt.Println(line)
		}

		key, err := GetPressedArrow()
		if err != nil {
			continue
		}

		if key == 65517 {
			if u.Selected > 0 {
				u.Selected--
			} else {
				u.Selected = len(u.Fields) - 1
			}
		}

		if key == 65516 {
			if u.Selected < len(u.Fields)-1 {
				u.Selected++
			} else {
				u.Selected = 0
			}
		}

		if key == 65514 || key == 13 {
			field := u.Fields[u.Selected]

			if field.UI != nil {
				field.UI.StartNonInitialize()
				break
			} else if field.Callback != nil {
				termbox.Close()

				Clear()
				width, _, err := term.GetSize(int(os.Stdout.Fd()))
				if err != nil {
					width = 50
				}

				fmt.Println("\033[1m" + u.RGB.Color("Running option: ") + "\033[1m" + field.Title + "\033[0m")
				fmt.Println("\033[1m" + u.RGB.Color(strings.Repeat(Dash, width)))

				start := time.Now()
				field.Callback()
				dur := time.Since(start)

				fmt.Println("\033[1m" + u.RGB.Color(strings.Repeat(Dash, width)))
				fmt.Println("\033[1m" + u.RGB.Color("Option has been ran. Execution time: ") + "\033[1m" + formatDuration(dur) + "\nPress Enter to return to the menu.\033[0m")

				_, _ = fmt.Scanln()

				_ = termbox.Init()
			}
		}

		if key == 65515 {
			if u.Back != nil {
				u.Back.StartNonInitialize()
				break
			}
		}

		if key == 3 {
			return
		}
	}

}

func (u *InteractiveUI) Start() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	defer termbox.Close()

	u.StartNonInitialize()
}
