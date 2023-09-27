package displaycase

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	display tea.Model
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, m.display.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	// Handle exit requests, then forward anything else to the displayed component
	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c":
			cmd = tea.Quit
			return m, cmd
		}
	}
	m.display, cmd = m.display.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return fmt.Sprintf("%s\n\n(ctrl+c to quit)\n", m.display.View())
}

func New(display tea.Model) Model {
	return Model{
		display: display,
	}
}
