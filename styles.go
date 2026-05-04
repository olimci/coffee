package coffee

import "github.com/charmbracelet/lipgloss"

var (
	MutedStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	InverseMutedStyle   = lipgloss.NewStyle().Background(lipgloss.Color("8")).Foreground(lipgloss.Color("7"))
	SuccessStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	InverseSuccessStyle = lipgloss.NewStyle().Background(lipgloss.Color("2")).Foreground(lipgloss.Color("0"))
	ErrorStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
	InverseErrorStyle   = lipgloss.NewStyle().Background(lipgloss.Color("1")).Foreground(lipgloss.Color("7"))
	InverseStyle        = lipgloss.NewStyle().Background(lipgloss.Color("4")).Foreground(lipgloss.Color("0"))
	AccentStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Bold(true)
	PromptStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	KeycapStyle         = lipgloss.NewStyle().Background(lipgloss.Color("4")).Foreground(lipgloss.Color("0")).Bold(true)
	BannerStyle         = lipgloss.NewStyle().Background(lipgloss.Color("4")).Foreground(lipgloss.Color("0"))
)
