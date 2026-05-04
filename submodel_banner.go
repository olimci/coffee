package coffee

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func NewBanner(text string) *Banner {
	return &Banner{text: text}
}

type Banner struct {
	text  string
	width int
}

func (m *Banner) Focusable() bool {
	return false
}

func (m *Banner) Init() tea.Cmd {
	return nil
}

func (m *Banner) Update(msg tea.Msg) (Submodel, tea.Cmd, string) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil, ""
	case msgBannerSet:
		if msg.banner != m {
			return m, nil, ""
		}
		m.text = msg.text
		return m, nil, ""
	case msgBannerClear:
		if msg.banner != m {
			return m, nil, ""
		}
		return nil, nil, ""
	default:
		return m, nil, ""
	}
}

func (m *Banner) View() string {
	if m.text == "" {
		return ""
	}

	lines := strings.Split(m.text, "\n")
	for i, line := range lines {
		lines[i] = renderBannerLine(line, m.width)
	}
	return strings.Join(lines, "\n")
}

func renderBannerLine(line string, width int) string {
	if width <= 0 {
		width = ansi.StringWidth(line)
	}

	line = ansi.Truncate(line, width, "")
	return lipgloss.PlaceHorizontal(
		width,
		lipgloss.Left,
		BannerStyle.Render(line),
		lipgloss.WithWhitespaceBackground(lipgloss.Color("4")),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("0")),
		lipgloss.WithWhitespaceChars("\u00a0"),
	)
}

var _ Submodel = (*Banner)(nil)
