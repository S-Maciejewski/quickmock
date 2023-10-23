package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"quickmock/definition"
)

// viewState is an enum representing the current view of the TUI.
type viewState int

const (
	ListView viewState = iota
	DetailView
	MethodSelect
	PathInput
	CodeInput
	ResponseInput
)

// model is the state of the TUI. It is updated by the Update method and rendered by the View method.
type model struct {
	state         viewState
	endpoints     *[]definition.Endpoint
	selectedIndex int                 // Index of the currently selected endpoint in ListView
	detailedView  definition.Endpoint // The endpoint in focus in DetailView
	creating      definition.Endpoint // The endpoint being created
	methods       [9]string           // Supported HTTP methods for selection
	input         string              // General input buffer
	errorMsg      string              // Error message to display in TUI
}

func (m model) Init() tea.Cmd {
	return nil
}

// Update handles messages sent to the TUI. Called by the Bubble Tea runtime.
// This method is the main logic of the TUI, so it has been divided into separate updates called for each view.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case ListView:
			return handleListViewUpdate(&m, msg)
		case DetailView:
			return handleDetailViewUpdate(&m, msg)
		case MethodSelect:
			return handleMethodSelectUpdate(&m, msg)
		case PathInput:
			return handlePathInputUpdate(&m, msg)
		case CodeInput:
			return handleCodeInputUpdate(&m, msg)
		case ResponseInput:
			return handleResponseInputUpdate(&m, msg)
		}
	}
	return m, nil
}

// View returns the current view of the TUI as a string. Called by the Bubble Tea runtime.
func (m model) View() string {
	switch m.state {
	case ListView:
		return getListView(m)
	case DetailView:
		return getDetailView(m)
	case MethodSelect:
		return getMethodSelect(m)
	case PathInput:
		return getPathInput(m)
	case CodeInput:
		return getCodeInput(m)
	case ResponseInput:
		return getResponseInput(m)
	}
	return ""
}

// RunTui starts the TUI. Called asynchronously from main.go.
func RunTui(endpoints *[]definition.Endpoint) {
	p := tea.NewProgram(model{endpoints: endpoints, state: ListView, methods: definition.SupportedHttpMethods()})
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running TUI: %v", err)
		return
	}
}
