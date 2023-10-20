package main

import (
	"flag"
	"fmt"
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

func main() {
	filePath := flag.String("f", "", "Path to YAML configuration file")
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
		// TODO: Start interactive mode here
		fmt.Println("Starting in interactive mode...")
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
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
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
	w.Write([]byte("Not Found"))
}
