package tui

import "fmt"

func getListView(m model) string {
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
}

func getDetailView(m model) string {
	view := fmt.Sprintf("Endpoint view\n"+
		"  esc: Back to list\n\n"+
		"Method: %s\nPath: %s\nCode: %d\nContent:\n%s\n",
		m.detailedView.Method,
		m.detailedView.Path,
		m.detailedView.Response.Code,
		m.detailedView.Response.Content,
	)
	return view
}

func getMethodSelect(m model) string {
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
}

func getPathInput(m model) string {
	view := fmt.Sprintf("Enter path: %s", m.input)
	view += fmt.Sprintf("\n%s", m.errorMsg)
	return view
}

func getCodeInput(m model) string {
	return fmt.Sprintf("Enter HTTP response status code: %s", m.input)
}

func getResponseInput(m model) string {
	return fmt.Sprintf("Enter response (confirm with `):\n%s", m.input)
}
