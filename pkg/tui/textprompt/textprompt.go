// Copyright Nitric Pty Ltd.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package textprompt

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/nitrictech/pearls/pkg/tui"
	"github.com/nitrictech/pearls/pkg/tui/validation"
	"github.com/nitrictech/pearls/pkg/tui/view"
)

type (
	errMsg error
)

type Model struct {
	ID               string
	textInput        textinput.Model
	Prompt           string
	Tag              string
	validate         validation.StringValidator
	validateInFlight validation.StringValidator
	focus            bool
	previous         string

	err error
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

type CompleteMsg struct {
	ID    string
	Value string
}

func (m *Model) submit() tea.Cmd {
	return func() tea.Msg {
		return CompleteMsg{
			ID:    m.ID,
			Value: m.textInput.Value(),
		}
	}
}

func (m Model) UpdateTextPrompt(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, tui.KeyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, tui.KeyMap.Enter):
			if m.textInput.Value() == "" {
				m.textInput.SetValue(m.textInput.Placeholder)
			}

			m.err = m.validate(m.textInput.Value())

			if m.err == nil {
				m.textInput.Blur()
				return m, m.submit()
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)

	// only clear/update error messages if the input has changed
	if m.previous != m.textInput.Value() {
		if m.textInput.Value() != "" {
			m.err = m.validateInFlight(m.textInput.Value())
		} else {
			m.err = nil
		}
	}

	m.previous = m.textInput.Value()

	return m, cmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.UpdateTextPrompt(msg)
}

var (
	labelStyle      = lipgloss.NewStyle().MarginTop(1)
	tagStyle        = lipgloss.NewStyle().Background(tui.Colors.Purple).Foreground(tui.Colors.White).Width(8).Align(lipgloss.Center)
	promptStyle     = lipgloss.NewStyle().MarginLeft(2)
	shiftRightStyle = lipgloss.NewStyle().MarginLeft(10)
	textStyle       = lipgloss.NewStyle().Foreground(tui.Colors.Gray)
	errorStyle      = lipgloss.NewStyle().Foreground(tui.Colors.Red).Italic(true).MarginTop(1)
)

func (m Model) View() string {
	renderer := view.New()

	renderer.AddRow(
		view.NewFragment(m.Tag).WithStyle(tagStyle),
		view.NewFragment(m.Prompt).WithStyle(promptStyle),
		view.Break(),
	).WithStyle(labelStyle)

	renderer.AddRow(view.WhenOr(
		m.textInput.Focused(),
		view.NewFragment(m.textInput.View()),
		view.NewFragment(m.textInput.Value()).WithStyle(textStyle),
	), view.When(
		m.err != nil,
		view.NewFragment(m.err).WithStyle(errorStyle),
	)).WithStyle(shiftRightStyle)

	return renderer.Render()
}

// Focus sets the focus state on the model. When the model is in focus it can
// receive keyboard input and the cursor will be shown.
func (m *Model) Focus() tea.Cmd {
	m.focus = true
	return m.textInput.Focus()
}

// Blur removes the focus state on the model.  When the model is blurred it can
// not receive keyboard input and the cursor will be hidden.
func (m *Model) Blur() {
	m.focus = false
	m.textInput.Blur()
}

func (m *Model) SetValue(value string) {
	m.textInput.SetValue(value)
}

func (m Model) Value() string {
	return m.textInput.Value()
}

type TextPromptArgs struct {
	ID                string
	Placeholder       string
	Validator         validation.StringValidator
	InFlightValidator validation.StringValidator
	Prompt            string
	Tag               string
}

func NewTextPrompt(id string, args TextPromptArgs) Model {
	ti := textinput.New()
	ti.CharLimit = 156
	ti.Width = 20
	ti.Placeholder = args.Placeholder

	return Model{
		ID:               id,
		textInput:        ti,
		Prompt:           args.Prompt,
		Tag:              args.Tag,
		validate:         args.Validator,
		validateInFlight: args.InFlightValidator,
		err:              nil,
	}
}
