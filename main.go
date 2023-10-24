package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"quickmock/definition"
	"quickmock/tui"
)

var endpoints []definition.Endpoint

// defaultEndpoints returns a slice of default endpoints to be used if no YAML file is provided.
// It's here to make it easier for user to understand how the TUI works.
func defaultEndpoints() []definition.Endpoint {
	return []definition.Endpoint{
		{
			Method: "GET",
			Path:   "/",
			Response: definition.Response{
				Code:    204,
				Content: "",
			},
		},
	}
}

func main() {
	filePath := flag.String("f", "", "Path to YAML configuration file")
	port := flag.Int64("p", 8080, "Port to listen on")
	detachedMode := flag.Bool("d", false, "Run in detached mode (no TUI)")
	flag.Parse()

	if *filePath != "" {
		definition.LoadEndpointsFromFile(*filePath, &endpoints)
	} else {
		endpoints = defaultEndpoints()
	}
	if *detachedMode {
		log.Printf("Starting quickmock in detached mode on port %d", *port)
	} else {
		log.Printf("Starting quickmock in interactive TUI mode on port %d", *port)
		go tui.RunTui(&endpoints)
	}

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
		return
	}
}

// handler handles all HTTP requests to the server. It mocks the routing and endpoints defined in the file or by TUI.
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
