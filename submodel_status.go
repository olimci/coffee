package coffee

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type statusState int

const (
	statusIdle statusState = iota
	statusWorking
	statusProgress
	statusSuccess
	statusError
)

const statusProgressMinWidth = 8

type msgStatusIdle struct {
	status  *Status
	message string
}

type msgStatusWorking struct {
	status  *Status
	message string
}

type msgStatusProgress struct {
	status  *Status
	message string
	percent float64
}

type msgStatusProgressValue struct {
	status  *Status
	percent float64
}

type msgStatusMessage struct {
	status  *Status
	message string
}

type msgStatusSuccess struct {
	status  *Status
	message string
}

type msgStatusError struct {
	status  *Status
	message string
}

type msgStatusClear struct {
	status *Status
}

func NewStatus(message string) *Status {
	spin := spinner.New()
	spin.Spinner = spinner.Pulse

	return &Status{
		state:   statusWorking,
		message: message,
		spinner: spin,
	}
}

type Status struct {
	state    statusState
	message  string
	progress float64
	spinner  spinner.Model
	width    int
}

func (m *Status) Focusable() bool {
	return false
}

func (m *Status) Init() tea.Cmd {
	if m.state == statusWorking {
		return m.spinner.Tick
	}
	return nil
}

func (m *Status) Update(msg tea.Msg) (Submodel, tea.Cmd, string) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil, ""
	case msgStatusIdle:
		if msg.status != m {
			return m, nil, ""
		}
		m.state = statusIdle
		m.message = msg.message
		return m, nil, ""
	case msgStatusWorking:
		if msg.status != m {
			return m, nil, ""
		}
		m.state = statusWorking
		m.message = msg.message
		return m, m.spinner.Tick, ""
	case msgStatusProgress:
		if msg.status != m {
			return m, nil, ""
		}
		m.state = statusProgress
		m.message = msg.message
		m.progress = clampProgress(msg.percent)
		return m, nil, ""
	case msgStatusProgressValue:
		if msg.status != m {
			return m, nil, ""
		}
		m.state = statusProgress
		m.progress = clampProgress(msg.percent)
		return m, nil, ""
	case msgStatusMessage:
		if msg.status != m {
			return m, nil, ""
		}
		m.message = msg.message
		return m, nil, ""
	case msgStatusSuccess:
		if msg.status != m {
			return m, nil, ""
		}
		m.state = statusSuccess
		m.message = msg.message
		return m, nil, ""
	case msgStatusError:
		if msg.status != m {
			return m, nil, ""
		}
		m.state = statusError
		m.message = msg.message
		return m, nil, ""
	case msgStatusClear:
		if msg.status != m {
			return m, nil, ""
		}
		return nil, nil, m.final()
	}

	if m.state != statusWorking {
		return m, nil, ""
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd, ""
}

func (m *Status) View() string {
	switch m.state {
	case statusIdle:
		return m.message
	case statusWorking:
		if m.message == "" {
			return m.spinner.View()
		}
		return fmt.Sprintf("%s %s", m.spinner.View(), m.message)
	case statusProgress:
		return renderProgressLine(m.width, m.message, m.progress)
	case statusSuccess:
		return SuccessStyle.Render(m.message)
	case statusError:
		return ErrorStyle.Render(m.message)
	default:
		return m.message
	}
}

func (m *Status) final() string {
	switch m.state {
	case statusSuccess, statusError:
		return m.View()
	default:
		return ""
	}
}

func renderProgressLine(width int, message string, percent float64) string {
	percentLabel := fmt.Sprintf("%3.0f%%", percent*100)
	message = strings.TrimSpace(message)

	if width <= 0 {
		width = 80
	}

	fixedWidth := lipgloss.Width(percentLabel)
	if message != "" {
		fixedWidth += lipgloss.Width(message) + 1
	}

	barWidth := width - fixedWidth - 1
	if barWidth < statusProgressMinWidth {
		barWidth = statusProgressMinWidth
	}

	bar := renderProgressBar(barWidth, percent)
	if message == "" {
		return fmt.Sprintf("%s %s", bar, percentLabel)
	}
	return fmt.Sprintf("%s %s %s", bar, percentLabel, message)
}

func renderProgressBar(width int, percent float64) string {
	percent = clampProgress(percent)
	filled := int(percent * float64(width))
	empty := width - filled
	return strings.Repeat("#", filled) + strings.Repeat(" ", empty)
}

func clampProgress(percent float64) float64 {
	switch {
	case percent < 0:
		return 0
	case percent > 1:
		return 1
	default:
		return percent
	}
}

var _ Submodel = (*Status)(nil)
