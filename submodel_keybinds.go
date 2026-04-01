package coffee

import (
	"fmt"
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
)

const keybindBufferSize = 16

type Keybind struct {
	Key         string
	Event       string
	Description string
}

type KeybindEvent struct {
	Key   string
	Event string
}

type msgKeybindClear struct {
	keybinds *Keybinds
}

type msgKeybindSet struct {
	keybinds *Keybinds
	bindings []Keybind
}

func newKeybinds(bindings []Keybind, out chan KeybindEvent) *Keybinds {
	keybinds := &Keybinds{
		out: out,
	}
	keybinds.setBindings(bindings)
	return keybinds
}

type Keybinds struct {
	bindings      []Keybind
	bindingsByKey map[string]Keybind
	out           chan KeybindEvent
	closed        bool
}

func (m *Keybinds) Focusable() bool {
	return false
}

func (m *Keybinds) Init() tea.Cmd {
	return nil
}

func (m *Keybinds) Update(msg tea.Msg) (Submodel, tea.Cmd, string) {
	switch msg := msg.(type) {
	case msgKeybindClear:
		if msg.keybinds != m {
			return m, nil, ""
		}
		m.close()
		return nil, nil, ""
	case msgKeybindSet:
		if msg.keybinds != m {
			return m, nil, ""
		}
		m.setBindings(msg.bindings)
		return m, nil, ""
	case tea.KeyMsg:
		if binding, ok := m.match(msg.String()); ok {
			m.send(KeybindEvent{Key: binding.Key, Event: binding.Event})
		}
	}

	return m, nil, ""
}

func (m *Keybinds) View() string {
	if len(m.bindings) == 0 {
		return ""
	}

	parts := make([]string, 0, len(m.bindings))
	for _, binding := range m.bindings {
		parts = append(parts, renderKeybind(binding))
	}

	return strings.Join(parts, "  ")
}

func (m *Keybinds) send(ev KeybindEvent) {
	select {
	case m.out <- ev:
	default:
	}
}

func (m *Keybinds) close() {
	if m.closed {
		return
	}
	m.closed = true
	close(m.out)
}

func (m *Keybinds) setBindings(bindings []Keybind) {
	m.bindings = make([]Keybind, 0, len(bindings))
	m.bindingsByKey = make(map[string]Keybind, len(bindings))
	for _, binding := range bindings {
		if binding.Key == "" || binding.Event == "" {
			continue
		}

		binding.Key = normalizeKey(binding.Key)
		m.bindings = append(m.bindings, binding)
		m.bindingsByKey[binding.Key] = binding
	}
}

func (m *Keybinds) match(key string) (Keybind, bool) {
	binding, ok := m.bindingsByKey[normalizeKey(key)]
	return binding, ok
}

func normalizeKey(key string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return ""
	}

	runes := []rune(key)
	if len(runes) == 1 && unicode.IsLetter(runes[0]) {
		return strings.ToLower(key)
	}

	return key
}

func renderKeybind(binding Keybind) string {
	description := strings.TrimSpace(binding.Description)
	if description == "" {
		description = binding.Event
	}

	key := normalizeKey(binding.Key)
	if key == "" {
		return description
	}

	return fmt.Sprintf("%s %s", inverseStyle.Render(" "+binding.Key+" "), description)
}

var _ Submodel = (*Keybinds)(nil)
