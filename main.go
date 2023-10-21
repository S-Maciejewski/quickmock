package main

import (
	"flag"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
)

type Endpoint struct {
	Method   string `yaml:"method"`
	Path     string `yaml:"path"`
	Response struct {
		Code    int    `yaml:"code"`
		Content string `yaml:"content"`
	} `yaml:"response"`
}

var endpoints []Endpoint

type model struct {
	count int
}

type message struct{}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		// Close the application on ctrl+c
		if msg.(tea.KeyMsg).Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		m.count++
	}
	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("You've pressed a key %d times. Press any key to continue.", m.count)
}

func main() {
	filePath := flag.String("f", "", "Path to YAML configuration file")
	port := flag.Int64("p", 8080, "Port to listen on")
	flag.Parse()

	if *filePath != "" {
		content, err := ioutil.ReadFile(*filePath)
		if err != nil {
			log.Fatalf("Error reading YAML file: %v", err)
		}

		err = yaml.Unmarshal(content, &endpoints)
		if err != nil {
			log.Fatalf("Error unmarshalling YAML content: %v", err)
		}
	} else {
		endpoints = []Endpoint{
			{
				Method: "GET",
				Path:   "/",
				Response: struct {
					Code    int    `yaml:"code"`
					Content string `yaml:"content"`
				}{
					Code:    204,
					Content: "quickmock default response",
				},
			},
		}
		fmt.Println("Starting quickmock. Press ctrl+c to exit the interactive TUI mode.")
		p := tea.NewProgram(model{})
		//TODO: Make this run in parallel with the HTTP server so that the TUI can be used to control the server
		if _, err := p.Run(); err != nil {
			log.Fatalf("Error running TUI: %v", err)
			return
		}
	}

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
		return
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	for _, endpoint := range endpoints {
		if r.URL.Path == endpoint.Path && r.Method == endpoint.Method {
			w.WriteHeader(endpoint.Response.Code)
			w.Write([]byte(endpoint.Response.Content))
			return
		}
	}

	w.WriteHeader(404)
	_, _ = w.Write([]byte("Not Found"))
}
