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

func main() {
	filePath := flag.String("f", "", "Path to YAML configuration file")
	port := flag.Int64("p", 8080, "Port to listen on")
	flag.Parse()

	if *filePath != "" {
		definition.ReadYaml(*filePath, &endpoints)
	} else {
		endpoints = []definition.Endpoint{
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
		fmt.Println("Starting quickmock in interactive TUI mode. Press ctrl+c to exit.")
		go tui.RunTui(&endpoints)
	}

	http.HandleFunc("/", handler)
	log.Println("Starting quickmock server")
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
