package internal

import (
	"errors"
	"testing"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

func press(code rune) tea.KeyPressMsg { return tea.KeyPressMsg{Code: code} }

func TestSelectModelCursorStaysInRange(t *testing.T) {
	var m tea.Model = selectModel{items: []string{"a", "b"}}
	for _, code := range []rune{tea.KeyUp, tea.KeyDown, tea.KeyDown, tea.KeyDown} {
		m, _ = m.Update(press(code))
	}
	if got := m.(selectModel).cursor; got != 1 {
		t.Fatalf("cursor = %d, want 1", got)
	}
}

func TestTextModelEnterBlockedByValidation(t *testing.T) {
	input := textinput.New()
	input.Focus()
	var m tea.Model = textModel{input: input, validate: func(s string) error {
		if s == "" {
			return errors.New("empty")
		}
		return nil
	}}
	m, cmd := m.Update(press(tea.KeyEnter))
	if cmd != nil || m.(textModel).err == nil {
		t.Fatalf("enter on invalid input: cmd = %v, err = %v, want nil cmd and error", cmd, m.(textModel).err)
	}
	m, _ = m.Update(tea.KeyPressMsg{Code: 'x', Text: "x"})
	if _, cmd = m.Update(press(tea.KeyEnter)); cmd == nil {
		t.Fatal("enter on valid input: want quit cmd")
	}
}
