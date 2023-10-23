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
	MethodSelect
	PathInput
	CodeInput
	ResponseInput
)

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
			case tea.KeyEsc, tea.KeyCtrlC:
				return m, tea.Quit
			default:
				if msg.String() == "c" {
					m.state = MethodSelect
				} else if msg.String() == "d" {
					if len(*m.endpoints) > 0 {
						*m.endpoints = append((*m.endpoints)[:m.selectedIndex], (*m.endpoints)[m.selectedIndex+1:]...)
						if m.selectedIndex > 0 {
							m.selectedIndex--
						}
					}
				}
			}

		case MethodSelect:
			switch msg.Type {
			case tea.KeyDown:
				if m.selectedIndex < len(m.methods)-1 {
					m.selectedIndex++
				}
			case tea.KeyUp:
				if m.selectedIndex > 0 {
					m.selectedIndex--
				}
			case tea.KeyEnter:
				m.creating.Method = m.methods[m.selectedIndex]
				// Reset selected index for next state
				m.selectedIndex = 0
				m.state = PathInput

			case tea.KeyEsc, tea.KeyCtrlC:
				m.state = ListView
			}

		case PathInput:
			switch msg.Type {
			case tea.KeyBackspace:
				m.input = m.input[:len(m.input)-1]
			case tea.KeyEsc, tea.KeyCtrlC:
				m.state = ListView
			case tea.KeyEnter:
				m.errorMsg = ""
				if definition.ValidatePath(m.input) {
					m.errorMsg = "Invalid path. It has to start with / and be a valid URL path"
					m.input = ""
					break
				}
				m.creating.Path = m.input
				m.input = ""
				m.state = CodeInput
			default:
				m.input += msg.String()
			}

		case CodeInput:
			switch msg.Type {
			case tea.KeyEnter:
				m.creating.Response.Code = definition.ParseStatusCode(m.input)
				m.input = ""
				m.state = ResponseInput
			case tea.KeyBackspace:
				m.input = m.input[:len(m.input)-1]
			case tea.KeyEsc, tea.KeyCtrlC:
				m.state = ListView
			default:
				m.input += msg.String()
			}

		case ResponseInput:
			switch msg.Type {
			case tea.KeyEsc, tea.KeyCtrlC:
				m.state = ListView
			case tea.KeyBackspace:
				m.input = m.input[:len(m.input)-1]
			case tea.KeyEnter:
				m.input += "\n"
			case tea.KeyTab:
				m.input += "\t"
			case tea.KeyUp, tea.KeyDown, tea.KeyLeft, tea.KeyRight:
				// Ignore these keys
				break
			default:
				if msg.String() == "`" { // Using the backtick '`' to confirm multi-line input.
					m.creating.Response.Content = m.input
					*m.endpoints = append(*m.endpoints, m.creating)

					// Reset state
					m.creating = definition.Endpoint{}
					m.input = ""
					m.state = ListView
				} else {
					m.input += msg.String()
				}

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
		view := "Controls:\n" +
			"  c: Create new endpoint\n" +
			"  d: Delete endpoint\n" +
			"  ↑/↓: Navigate endpoints\n" +
			"  enter: View endpoint\n" +
			"  esc: Exit\n\n"
		for i, endpoint := range *m.endpoints {
			prefix := "  "
			if i == m.selectedIndex {
				prefix = "> "
			}
			view += fmt.Sprintf("%s%s %s\n", prefix, endpoint.Method, endpoint.Path)
		}
		return view

	case DetailView:
		view := fmt.Sprintf("Endpoint view\n"+
			"  esc: Back to list\n\n"+
			"Method: %s\nPath: %s\nCode: %d\nContent: %s\n",
			m.detailedView.Method,
			m.detailedView.Path,
			m.detailedView.Response.Code,
			m.detailedView.Response.Content,
		)
		return view

	case MethodSelect:
		view := "Select a method:\n" +
			"  ↑/↓: Navigate methods\n" +
			"  enter: Select method\n" +
			"  esc: Cancel\n\n"
		for i, method := range m.methods {
			prefix := "  "
			if i == m.selectedIndex {
				prefix = "> "
			}
			view += fmt.Sprintf("%s%s\n", prefix, method)
		}
		return view

	case PathInput:
		view := fmt.Sprintf("Enter path:\n%s", m.input)
		view += fmt.Sprintf("\n%s", m.errorMsg)
		return view

	case CodeInput:
		return fmt.Sprintf("Enter HTTP response status code: %s", m.input)

	case ResponseInput:
		return fmt.Sprintf("Enter response (confirm with `):\n%s", m.input)
	}

	return ""
}

func RunTui(endpoints *[]definition.Endpoint) {
	p := tea.NewProgram(model{endpoints: endpoints, state: ListView, methods: definition.SupportedHttpMethods()})
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running TUI: %v", err)
		return
	}
}
