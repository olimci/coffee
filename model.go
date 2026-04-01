package coffee

import (
	"slices"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	header []item
	body   []item
	footer []item

	submodels map[*submodelEntry]struct{}
	focus     []*submodelEntry
	width     int
	height    int
	altScreen bool
}

func newModel(o *options) *model {
	return &model{
		submodels: make(map[*submodelEntry]struct{}),
		altScreen: o.altScreen,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case msgLog:
		m.appendText(msg.section, msg.message)
		return m, nil

	case msgClear:
		m.clearSection(msg.section)
		return m, nil

	case msgSubmodel:
		entry := &submodelEntry{submodel: msg.submodel}
		m.submodels[entry] = struct{}{}
		m.appendSubmodel(msg.section, entry)
		if !submodelFocusable(entry.submodel) {
			return m, entry.submodel.Init()
		}
		return m, tea.Batch(entry.submodel.Init(), m.pushFocus(entry, msg.behind))

	case msgWindowTitle:
		return m, tea.SetWindowTitle(msg.title)
	}

	if _, ok := msg.(tea.KeyMsg); ok {
		focused, hasFocus := m.focused()
		if hasFocus {
			return m, m.updateSubmodel(focused, msg)
		}
	}

	cmds := make([]tea.Cmd, 0, len(m.submodels))
	for entry := range m.submodels {
		cmd := m.updateSubmodel(entry, msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *model) appendText(section Section, text string) {
	items := m.section(section)
	*items = append(*items, item{
		text: text,
	})
}

func (m *model) appendSubmodel(section Section, entry *submodelEntry) {
	items := m.section(section)
	*items = append(*items, item{
		entry: entry,
	})
}

func (m *model) clearSection(section Section) {
	items := m.section(section)
	*items = slices.DeleteFunc(*items, func(item item) bool {
		return item.entry == nil
	})
}

func (m *model) section(section Section) *[]item {
	switch section {
	case SectionHeader:
		return &m.header
	case SectionFooter:
		return &m.footer
	default:
		return &m.body
	}
}

func (m *model) finishSubmodel(entry *submodelEntry, final string) tea.Cmd {
	oldFocus, _ := m.focused()
	delete(m.submodels, entry)
	m.focus = slices.DeleteFunc(m.focus, func(e *submodelEntry) bool {
		return e == entry
	})

	for _, section := range []*[]item{&m.header, &m.body, &m.footer} {
		if i := slices.IndexFunc(*section, func(item item) bool {
			return item.entry == entry
		}); i != -1 {
			if final == "" {
				*section = slices.Delete(*section, i, i+1)
			} else {
				(*section)[i] = item{text: final}
			}
			break
		}
	}

	newFocus, _ := m.focused()
	return m.notifyFocusChange(oldFocus, newFocus)
}

func (m *model) pushFocus(entry *submodelEntry, behind bool) tea.Cmd {
	oldFocus, _ := m.focused()
	if behind {
		m.focus = append([]*submodelEntry{entry}, m.focus...)
	} else {
		m.focus = append(m.focus, entry)
	}

	newFocus, _ := m.focused()
	return m.notifyFocusChange(oldFocus, newFocus)
}

func (m *model) focused() (*submodelEntry, bool) {
	if len(m.focus) == 0 {
		return nil, false
	}

	return m.focus[len(m.focus)-1], true
}

func (m *model) notifyFocusChange(old, new *submodelEntry) tea.Cmd {
	if old == new {
		return nil
	}

	cmds := make([]tea.Cmd, 0, 2)
	if old != nil {
		if cmd := m.updateSubmodel(old, MsgFocusLost{}); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	if new != nil {
		if cmd := m.updateSubmodel(new, MsgFocusGained{}); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

func (m *model) updateSubmodel(entry *submodelEntry, msg tea.Msg) tea.Cmd {
	next, cmd, final := entry.submodel.Update(msg)
	if next != nil {
		entry.submodel = next
		return cmd
	}

	return tea.Batch(cmd, m.finishSubmodel(entry, final))
}
