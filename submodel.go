package coffee

import tea "github.com/charmbracelet/bubbletea"

type Submodel interface {
	Init() tea.Cmd
	Update(tea.Msg) (Submodel, tea.Cmd, string)
	View() string
}

type focusableSubmodel interface {
	Focusable() bool
}

func submodelFocusable(submodel Submodel) bool {
	focusable, ok := submodel.(focusableSubmodel)
	return !ok || focusable.Focusable()
}
