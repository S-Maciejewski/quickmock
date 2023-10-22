package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"quickmock/definition"
)

type viewState int

const (
	ListView viewState = iota
	DetailView
)

type model struct {
	state         viewState
	endpoints     *[]definition.Endpoint
	selectedIndex int                 // Index of the currently selected endpoint in ListView
	detailedView  definition.Endpoint // The endpoint in focus in DetailView
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case ListView:
			switch msg.Type {
			case tea.KeyDown:
				if m.selectedIndex < len(*m.endpoints)-1 {
					m.selectedIndex++
				}
			case tea.KeyUp:
				if m.selectedIndex > 0 {
					m.selectedIndex--
				}
			case tea.KeyEnter:
				m.detailedView = (*m.endpoints)[m.selectedIndex]
				m.state = DetailView
			case tea.KeyCtrlC:
				return m, tea.Quit
			}

		case DetailView:
			switch msg.Type {
			case tea.KeyEsc:
				m.state = ListView
			case tea.KeyCtrlC:
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	switch m.state {
	case ListView:
		view := ""
		for i, endpoint := range *m.endpoints {
			prefix := "  "
			if i == m.selectedIndex {
				prefix = "> "
			}
			view += fmt.Sprintf("%s%s %s\n", prefix, endpoint.Method, endpoint.Path)
		}
		return view

	case DetailView:
		return fmt.Sprintf("Endpoint view, press esc to exit\n"+
			"Method: %s\nPath: %s\nCode: %d\nContent: %s\n",
			m.detailedView.Method,
			m.detailedView.Path,
			m.detailedView.Response.Code,
			m.detailedView.Response.Content,
		)
	}
	return ""
}

func RunTui(endpoints *[]definition.Endpoint) {
	p := tea.NewProgram(model{endpoints: endpoints, state: ListView})
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running TUI: %v", err)
		return
	}
}
