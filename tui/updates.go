package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"quickmock/definition"
)

func handleListViewUpdate(m *model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		if len(*m.endpoints) != 0 {
			m.detailedView = (*m.endpoints)[m.selectedIndex]
			m.state = DetailView
		}
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
	return *m, nil
}

func handleDetailViewUpdate(m *model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.state = ListView
	case tea.KeyCtrlC:
		return m, tea.Quit
	}
	return *m, nil
}

func handleMethodSelectUpdate(m *model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
	return *m, nil
}

func handlePathInputUpdate(m *model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
	return *m, nil
}

func handleCodeInputUpdate(m *model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
	return *m, nil
}

func handleResponseInputUpdate(m *model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
	return m, nil
}
