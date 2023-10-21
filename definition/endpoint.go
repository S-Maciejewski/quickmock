package definition

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Response struct {
	Code    int    `yaml:"code"`
	Content string `yaml:"content"`
}

type Endpoint struct {
	Method   string   `yaml:"method"`
	Path     string   `yaml:"path"`
	Response Response `yaml:"response"`
}

func ReadYaml(filePath string, endpoints *[]Endpoint) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	err = yaml.Unmarshal(content, &endpoints)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML content: %v", err)
	}
	log.Println("Loaded endpoints from the YAML file")
}
