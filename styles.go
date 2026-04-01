package coffee

import "github.com/charmbracelet/lipgloss"

var (
	mutedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	inverseStyle = lipgloss.NewStyle().Reverse(true)
)
