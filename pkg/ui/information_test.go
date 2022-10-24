package ui

import (
	"fmt"
	"github.com/fatih/color"
	"testing"
)

func TestInformation(t *testing.T) {
	Clear()
	s, err := ProgressBar("Streaming", 50, color.HiRedString)
	fmt.Println(s)

	sections := []*Section{
		{
			Header: "Main",
			Fields: []*Field{
				{
					Header:      "Hello",
					HeaderColor: color.BlueString,
					Value:       "World",
				},
			},
			Color: color.GreenString,
		},
	}

	printable, err := ToPrintable(sections)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(printable)
	}
}
