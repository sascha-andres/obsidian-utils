package internal

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type textModel struct {
	label    string
	input    textinput.Model
	validate func(string) error
	err      error
}

func (m textModel) Init() tea.Cmd { return nil }

func (m textModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "ctrl+c", "esc":
			return m, tea.Interrupt
		case "enter":
			if m.validate != nil {
				if m.err = m.validate(m.input.Value()); m.err != nil {
					return m, nil
				}
			}
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.err = nil
	return m, cmd
}

func (m textModel) View() tea.View {
	s := m.label + ": " + m.input.View() + "\n"
	if m.err != nil {
		s += "✗ " + m.err.Error() + "\n"
	}
	return tea.NewView(s)
}

// PromptText runs a textual prompt
func PromptText(label, defaultValue string, val func(string) error) (string, error) {
	input := textinput.New()
	input.Prompt = ""
	input.SetValue(defaultValue)
	input.Focus()
	result, err := tea.NewProgram(textModel{label: label, input: input, validate: val}).Run()
	if err != nil {
		return "", err
	}
	return result.(textModel).input.Value(), nil
}

type selectModel struct {
	label  string
	items  []string
	cursor int
}

func (m selectModel) Init() tea.Cmd { return nil }

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "ctrl+c", "esc":
			return m, tea.Interrupt
		case "enter":
			return m, tea.Quit
		case "up", "k":
			m.cursor = max(m.cursor-1, 0)
		case "down", "j":
			m.cursor = min(m.cursor+1, len(m.items)-1)
		}
	}
	return m, nil
}

func (m selectModel) View() tea.View {
	var b strings.Builder
	b.WriteString(m.label + ":\n")
	for i, item := range m.items {
		marker := "  "
		if i == m.cursor {
			marker = "▸ "
		}
		b.WriteString(marker + item + "\n")
	}
	return tea.NewView(b.String())
}

// PromptSelect lets the user select one of items
func PromptSelect(label string, items []string) (string, error) {
	if len(items) == 0 {
		return "", fmt.Errorf("no items to select for %q", label)
	}
	result, err := tea.NewProgram(selectModel{label: label, items: items}).Run()
	if err != nil {
		return "", err
	}
	return result.(selectModel).items[result.(selectModel).cursor], nil
}
