package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nitrictech/pearls/internal/displaycase"
	"github.com/nitrictech/pearls/pkg/tui/inlinelist"
)

func main() {
	list := inlinelist.New(inlinelist.Args{
		Items: []string{
			"item one",
			"item two",
			"item three",
			"item four",
			"item five",
			"item six",
			"item seven",
			"item eight",
			"item nine",
			"item ten",
		},
		MaxDisplayedItems: 5,
	})

	displaycase := displaycase.New(list)

	tea.NewProgram(displaycase).Run()

}
