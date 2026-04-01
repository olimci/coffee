package coffee

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/olimci/coffee/internal/promise"
)

func NewSelect(prompt string, options []string, index int) *Select {
	if index < 0 {
		index = 0
	}
	if len(options) > 0 && index >= len(options) {
		index = len(options) - 1
	}

	return &Select{
		prompt:  prompt,
		options: options,
		index:   index,
	}
}

type Select struct {
	prompt  string
	options []string
	index   int
	focused bool
	resolve *promise.Resolver[string]
}

func (m *Select) Init() tea.Cmd {
	return nil
}

func (m *Select) withResolve(resolve promise.Resolver[string]) *Select {
	m.resolve = &resolve
	return m
}

func (m *Select) Update(msg tea.Msg) (Submodel, tea.Cmd, string) {
	switch msg := msg.(type) {
	case MsgFocusGained:
		m.focused = true
		return m, nil, ""
	case MsgFocusLost:
		m.focused = false
		return m, nil, ""
	case tea.KeyMsg:
		if !m.focused {
			return m, nil, ""
		}

		switch msg.String() {
		case "up", "k":
			if len(m.options) == 0 {
				return m, nil, ""
			}
			m.index--
			if m.index < 0 {
				m.index = len(m.options) - 1
			}
		case "down", "j":
			if len(m.options) == 0 {
				return m, nil, ""
			}
			m.index++
			if m.index >= len(m.options) {
				m.index = 0
			}
		case "enter":
			if len(m.options) == 0 {
				return m, nil, ""
			}
			selected := m.options[m.index]
			m.resolve.Ok(selected)
			return nil, nil, m.final(selected)
		}
	}

	return m, nil, ""
}

func (m *Select) View() string {
	lines := make([]string, 0, len(m.options)+1)
	if m.prompt != "" {
		lines = append(lines, m.prompt)
	}

	for i, option := range m.options {
		prefix := "  "
		if i == m.index {
			prefix = "> "
		}
		lines = append(lines, prefix+option)
	}

	return strings.Join(lines, "\n")
}

func (m *Select) final(selected string) string {
	if m.prompt == "" {
		return selected
	}
	return fmt.Sprintf("%s %s", m.prompt, selected)
}

var _ Submodel = (*Select)(nil)
