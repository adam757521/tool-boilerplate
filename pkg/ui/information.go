package ui

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/term"
	"os"
	"strings"
)

const (
	Dash            = "━"
	VerticalDash    = "┃"
	CornerLeft      = "┏"
	CornerRight     = "┓"
	CornerDownLeft  = "┗"
	CornerDownRight = "┛"
)

type colorFunc func(string, ...any) string

type Field struct {
	Header      string
	HeaderColor colorFunc
	Value       string
	ValueColor  colorFunc
}

type Section struct {
	Header string
	Fields []*Field
	Width  int
	Color  colorFunc
}

func (s *Section) GetWidth() (int, int) {
	start := 0
	width := s.Width

	for _, field := range s.Fields {
		fieldStart := len(field.Header) + 1
		if fieldStart > start {
			start = fieldStart
		}
	}

	for _, field := range s.Fields {
		fieldLength := start + len(field.Value) + 7
		if fieldLength > width {
			width = fieldLength
		}
	}

	return start, width
}

func (s *Section) ToStringArray() []string {
	start, width := s.GetWidth()

	var printable []string

	categoryHeader := CornerLeft + strings.Repeat(Dash, width-2) + CornerRight
	if s.Header != "" {
		dashes := width - 5 - len(s.Header)
		categoryHeader = CornerLeft + Dash + " " + s.Header + " " + strings.Repeat(Dash, dashes) + CornerRight
	}

	printable = append(printable, s.Color(categoryHeader))

	for _, field := range s.Fields {
		padding := start - len(field.Header)
		prefix := s.Color(VerticalDash+"[") + ">" + s.Color("] ")

		headerColored := field.Header
		if field.HeaderColor != nil {
			headerColored = field.HeaderColor(field.Header)
		}

		valueColored := field.Value
		if field.ValueColor != nil {
			valueColored = field.ValueColor(field.Value)
		}

		if padding < 0 {
			padding = 0
		}
		fieldStr := headerColored + strings.Repeat(" ", padding) + valueColored

		padding = width - (6 + len(field.Header+field.Value) + padding)
		if padding < 0 {
			padding = 0
		}
		printable = append(printable, prefix+fieldStr+strings.Repeat(" ", padding)+s.Color(VerticalDash))
	}

	printable = append(printable, s.Color(CornerDownLeft+strings.Repeat(Dash, width-2)+CornerDownRight))

	return printable
}

func FitSections(sections []*Section, width int) int {
	totalWidth := 0
	for i, section := range sections {
		_, sectionWidth := section.GetWidth()

		totalWidth += sectionWidth
		if totalWidth > width {
			return i
		}
	}

	return -1
}

func LongestLength(arrays [][]string) int {
	longest := -1

	for _, array := range arrays {
		arrayLength := len(array)
		if arrayLength > longest {
			longest = arrayLength
		}
	}

	return longest
}

func SectionsToString(sections []*Section) string {
	var arrays [][]string

	for _, section := range sections {
		arrays = append(arrays, section.ToStringArray())
	}

	var sb strings.Builder

	longest := LongestLength(arrays)
	for i := 0; i < longest; i++ {
		for j, array := range arrays {
			if len(array) <= i {
				section := sections[j]
				_, width := section.GetWidth()

				sb.WriteString(strings.Repeat(" ", width))
			} else {
				sb.WriteString(array[i])
			}
		}

		if i != longest-1 {
			sb.WriteRune('\n')
		}
	}

	return sb.String()
}

func insertSubstring(str string, substr string, color colorFunc, index int) string {
	return str[:index] + color(substr) + str[index+len(substr):]
}

func ProgressBar(header string, percent int, pColor colorFunc) (string, error) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return "", err
	}

	headerStr := fmt.Sprintf(" %s (%d%%) ", header, percent)

	space := width - 2

	var sb strings.Builder
	filler := strings.Repeat(Dash, space)

	headerS := pColor(CornerLeft + Dash + headerStr + strings.Repeat(Dash, width-len(headerStr)-3) + CornerRight)
	sb.WriteString(headerS)

	bgWhite := color.New(color.BgHiWhite)
	filled := 0
	if percent != 0 {
		filled = space / (100 / percent)
	}
	remaining := space - filled
	progressFiller := bgWhite.Sprint(strings.Repeat(" ", filled)) + strings.Repeat(" ", remaining)

	sb.WriteString(pColor(VerticalDash) + progressFiller + pColor(VerticalDash))

	sb.WriteString(pColor(CornerDownLeft + filler + CornerDownRight))

	return sb.String(), nil
}

func ToPrintable(sections []*Section) (string, error) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return "", err
	}

	var sb strings.Builder

	for len(sections) > 0 {
		if sb.Len() != 0 {
			sb.WriteRune('\n')
		}

		index := FitSections(sections, width)
		if index == 0 {
			return "", errors.New("terminal width is to small")
		}

		if index != -1 {
			sb.WriteString(SectionsToString(sections[:index]))
			sections = sections[index:]
		} else {
			sb.WriteString(SectionsToString(sections))
			break
		}
	}

	return sb.String(), nil
}
