package coffee

import "github.com/charmbracelet/lipgloss"

var (
	MutedStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	InverseMutedStyle   = lipgloss.NewStyle().Reverse(true).Foreground(lipgloss.Color("8"))
	SuccessStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	InverseSuccessStyle = lipgloss.NewStyle().Reverse(true).Foreground(lipgloss.Color("2"))
	ErrorStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	InverseErrorStyle   = lipgloss.NewStyle().Reverse(true).Foreground(lipgloss.Color("1"))
	InverseStyle        = lipgloss.NewStyle().Reverse(true)
)
