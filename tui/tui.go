package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"quickmock/definition"
)

type model struct {
	count     int
	endpoints *[]definition.Endpoint
}

type message struct{}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		if msg.String() == "a" {
			// Add a mock endpoint as an example
			*m.endpoints = append(*m.endpoints, definition.Endpoint{
				Method: "GET",
				Path:   fmt.Sprintf("/mock%d", m.count),
				Response: struct {
					Code    int    `yaml:"code"`
					Content string `yaml:"content"`
				}{
					Code:    200,
					Content: fmt.Sprintf("Mock response %d", m.count),
				},
			})
		}
		m.count++
	}
	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("You've pressed a key %d times. Press any key to continue.", m.count)
}

func RunTui(endpoints *[]definition.Endpoint) {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running TUI: %v", err)
		return
	}
}
