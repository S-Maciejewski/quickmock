package definition

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"strconv"
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

func ValidHttpStatusCodes() [63]string {
	return [63]string{
		"100", "101", "102", "103",
		"200", "201", "202", "203", "204", "205", "206", "207", "208", "226",
		"300", "301", "302", "303", "304", "305", "306", "307", "308",
		"400", "401", "402", "403", "404", "405", "406", "407", "408", "409",
		"410", "411", "412", "413", "414", "415", "416", "417", "418", "421",
		"422", "423", "424", "425", "426", "428", "429", "431", "451",
		"500", "501", "502", "503", "504", "505", "506", "507", "508", "510",
		"511",
	}
}

func SupportedHttpMethods() [9]string {
	return [9]string{
		"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS", "CONNECT", "TRACE",
	}
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

// ParseStatusCode converts a code as a string to an int, and returns 200 if the conversion fails or the code isn't valid HTTP response status code.
func ParseStatusCode(code string) int {
	codeInt, err := strconv.Atoi(code)
	if err != nil {
		return 200
	}
	for _, validCode := range ValidHttpStatusCodes() {
		if code == validCode {
			return codeInt
		}
	}
	return 200
}
