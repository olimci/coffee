package coffee

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/olimci/coffee/internal/promise"
)

func TestActionPromptSelectsDefaultOnEnter(t *testing.T) {
	model, err := NewAction("run hook?", []Action{
		{Value: "run_once", Label: "run once"},
		{Value: "skip", Label: "skip"},
		{Value: "trust", Label: "trust and run"},
	}, "skip")
	if err != nil {
		t.Fatalf("NewAction() error = %v", err)
	}
	p, resolve := promise.New[string]()
	model = model.withResolve(resolve)
	model.Update(MsgFocusGained{})
	done := make(chan string, 1)
	go func() {
		next, _, final := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		if next != nil {
			t.Errorf("Update() model = %#v, want nil", next)
		}
		done <- final
	}()
	got, err := p.Await(t.Context())
	if err != nil {
		t.Fatalf("Await() error = %v", err)
	}
	if final := <-done; final == "" {
		t.Fatal("Update() final output is empty")
	}
	if got != "skip" {
		t.Fatalf("selected = %q, want skip", got)
	}
}

func TestActionPromptUsesUniqueShortcut(t *testing.T) {
	model, err := NewAction("run hook?", []Action{
		{Value: "run_once", Label: "run once"},
		{Value: "skip", Label: "skip"},
		{Value: "trust", Label: "trust and run"},
	}, "run_once")
	if err != nil {
		t.Fatalf("NewAction() error = %v", err)
	}
	p, resolve := promise.New[string]()
	model = model.withResolve(resolve)
	model.Update(MsgFocusGained{})
	done := make(chan struct{}, 1)
	go func() {
		next, _, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("t")})
		if next != nil {
			t.Errorf("Update() model = %#v, want nil", next)
		}
		done <- struct{}{}
	}()
	got, err := p.Await(t.Context())
	if err != nil {
		t.Fatalf("Await() error = %v", err)
	}
	<-done
	if got != "trust" {
		t.Fatalf("selected = %q, want trust", got)
	}
}

func TestActionPromptRejectsInvalidDefault(t *testing.T) {
	_, err := NewAction("run hook?", []Action{
		{Value: "run_once", Label: "run once"},
	}, "skip")
	if err == nil {
		t.Fatal("NewAction() error = nil, want invalid default")
	}
}
