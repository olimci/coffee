package coffee

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/olimci/coffee/internal/promise"
)

type Action struct {
	Value string
	Label string
}

func NewAction(prompt string, actions []Action, defaultValue string) (*ActionPrompt, error) {
	if len(actions) == 0 {
		return nil, fmt.Errorf("actions cannot be empty")
	}

	index := 0
	seen := make(map[string]bool, len(actions))
	for i, action := range actions {
		if strings.TrimSpace(action.Value) == "" {
			return nil, fmt.Errorf("actions[%d]: value is required", i)
		}
		if strings.TrimSpace(action.Label) == "" {
			return nil, fmt.Errorf("actions[%d]: label is required", i)
		}
		if seen[action.Value] {
			return nil, fmt.Errorf("actions[%d]: duplicate value %q", i, action.Value)
		}
		seen[action.Value] = true
		if action.Value == defaultValue {
			index = i
		}
	}
	if defaultValue != "" && !seen[defaultValue] {
		return nil, fmt.Errorf("default value not found in actions")
	}

	return &ActionPrompt{
		prompt:  prompt,
		actions: actions,
		index:   index,
	}, nil
}

type ActionPrompt struct {
	prompt  string
	actions []Action
	index   int
	focused bool
	resolve *promise.Resolver[string]
}

func (m *ActionPrompt) Init() tea.Cmd {
	return nil
}

func (m *ActionPrompt) withResolve(resolve promise.Resolver[string]) *ActionPrompt {
	m.resolve = &resolve
	return m
}

func (m *ActionPrompt) Update(msg tea.Msg) (Submodel, tea.Cmd, string) {
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
		case "left", "shift+tab", "h":
			m.index--
			if m.index < 0 {
				m.index = len(m.actions) - 1
			}
		case "right", "tab", "l":
			m.index++
			if m.index >= len(m.actions) {
				m.index = 0
			}
		case "enter":
			selected := m.actions[m.index]
			m.resolve.Ok(selected.Value)
			return nil, nil, m.final(selected)
		default:
			if msg.Type == tea.KeyRunes {
				selected, ok := m.shortcut(string(msg.Runes))
				if ok {
					m.resolve.Ok(selected.Value)
					return nil, nil, m.final(selected)
				}
			}
		}
	}

	return m, nil, ""
}

func (m *ActionPrompt) View() string {
	parts := make([]string, 0, len(m.actions)+1)
	if m.prompt != "" {
		parts = append(parts, PromptStyle.Render(m.prompt))
	}
	for i, action := range m.actions {
		label := " " + action.Label + " "
		if i == m.index {
			parts = append(parts, InverseStyle.Render(label))
			continue
		}
		parts = append(parts, MutedStyle.Render(action.Label))
	}
	return strings.Join(parts, " ")
}

func (m *ActionPrompt) final(selected Action) string {
	if m.prompt == "" {
		return selected.Label
	}
	return fmt.Sprintf("%s %s", PromptStyle.Render(m.prompt), AccentStyle.Render(selected.Label))
}

func (m *ActionPrompt) shortcut(runes string) (Action, bool) {
	if len(runes) != 1 {
		return Action{}, false
	}
	needle := strings.ToLower(runes)
	var match Action
	var matched bool
	for _, action := range m.actions {
		label := strings.TrimSpace(action.Label)
		if label == "" || strings.ToLower(label[:1]) != needle {
			continue
		}
		if matched {
			return Action{}, false
		}
		match = action
		matched = true
	}
	return match, matched
}

var _ Submodel = (*ActionPrompt)(nil)
