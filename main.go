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

func getDummyEndpoints() []definition.Endpoint {
	return []definition.Endpoint{
		{
			Method: "GET",
			Path:   "/",
			Response: definition.Response{
				Code:    204,
				Content: "",
			},
		},
		{
			Method: "POST",
			Path:   "/post-test",
			Response: definition.Response{
				Code:    200,
				Content: "quickmock default POST response",
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
		definition.ReadYaml(*filePath, &endpoints)
	} else {
		endpoints = getDummyEndpoints()
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
