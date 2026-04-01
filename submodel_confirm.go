package coffee

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/olimci/coffee/internal/promise"
)

func NewConfirm(prompt string, confirm bool) *Confirm {
	return &Confirm{
		prompt:  prompt,
		confirm: confirm,
	}
}

type Confirm struct {
	prompt  string
	confirm bool
	focused bool
	resolve *promise.Resolver[bool]
}

func (m *Confirm) Init() tea.Cmd {
	return nil
}

func (m *Confirm) withResolve(resolve promise.Resolver[bool]) *Confirm {
	m.resolve = &resolve
	return m
}

func (m *Confirm) Update(msg tea.Msg) (Submodel, tea.Cmd, string) {
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

		switch msg.Type {
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "y", "Y":
				m.confirm = true
				m.resolve.Ok(m.confirm)
				return nil, nil, m.final()
			case "n", "N":
				m.confirm = false
				m.resolve.Ok(m.confirm)
				return nil, nil, m.final()
			}
		case tea.KeyEnter:
			m.resolve.Ok(m.confirm)
			return nil, nil, m.final()
		}
	}

	return m, nil, ""
}

func (m *Confirm) View() string {
	if m.confirm {
		return fmt.Sprintf("%s Yn", m.prompt)
	}
	return fmt.Sprintf("%s yN", m.prompt)
}

func (m *Confirm) final() string {
	if m.confirm {
		return fmt.Sprintf("%s %s", m.prompt, successStyle.Render("YES"))
	}
	return fmt.Sprintf("%s %s", m.prompt, errorStyle.Render("NO"))
}

var _ Submodel = (*Confirm)(nil)
