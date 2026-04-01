package coffee

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/olimci/coffee/internal/promise"
)

func NewInput() *Input {
	in := textinput.New()
	in.Prompt = ": "
	in.PlaceholderStyle = mutedStyle
	in.Cursor.Style = inverseStyle

	return &Input{
		input: in,
	}
}

type Input struct {
	input           textinput.Model
	discardOnSubmit bool
	resolve         *promise.Resolver[string]
}

func (m *Input) WithPlaceholder(placeholder string) *Input {
	m.input.Placeholder = placeholder
	return m
}

func (m *Input) WithValue(value string) *Input {
	m.input.SetValue(value)
	return m
}

func (m *Input) WithValidate(validate func(string) error) *Input {
	m.input.Validate = validate
	return m
}

func (m *Input) WithCharLimit(limit int) *Input {
	m.input.CharLimit = limit
	return m
}

func (m *Input) WithSuggestions(suggestions []string) *Input {
	m.input.SetSuggestions(suggestions)
	m.input.ShowSuggestions = len(suggestions) > 0
	return m
}

func (m *Input) WithWidth(width int) *Input {
	m.input.Width = width
	return m
}

func (m *Input) WithDiscardSubmitted(discard bool) *Input {
	m.discardOnSubmit = discard
	return m
}

func (m *Input) withResolve(resolve promise.Resolver[string]) *Input {
	m.resolve = &resolve
	return m
}

func (m *Input) Init() tea.Cmd {
	return nil
}

func (m *Input) Update(msg tea.Msg) (Submodel, tea.Cmd, string) {
	switch msg := msg.(type) {
	case MsgFocusGained:
		return m, m.input.Focus(), ""
	case MsgFocusLost:
		m.input.Blur()
		return m, nil, ""
	case tea.KeyMsg:
		if !m.input.Focused() {
			return m, nil, ""
		}

		switch msg.Type {
		case tea.KeyEnter:
			if m.input.Err != nil {
				return m, nil, ""
			}
			m.resolve.Ok(m.input.Value())
			return nil, nil, m.final()
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd, ""
}

func (m *Input) View() string {
	if m.input.Err != nil {
		return fmt.Sprintf("%s %s", m.input.View(), errorStyle.Render(m.input.Err.Error()))
	}
	return m.input.View()
}

func (m *Input) final() string {
	if m.discardOnSubmit {
		return ""
	}

	return fmt.Sprintf("%s%s", m.input.Prompt, m.input.Value())
}

var _ Submodel = (*Input)(nil)
