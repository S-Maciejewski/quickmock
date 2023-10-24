package definition

import (
	"encoding/json"
	"github.com/go-openapi/spec"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

// ValidatePath is a simple URL path validator.
func ValidatePath(path string) bool {
	if len(path) > 0 && path[0] != '/' {
		return true
	}
	return false
}

// LoadEndpointsFromFile is the main function for loading endpoints from a file, based on the file extension.
// It supports YAML and JSON files in both quickmock's custom format and OpenAPI format.
func LoadEndpointsFromFile(filePath string, endpoints *[]Endpoint) {
	ext := filepath.Ext(filePath)
	content := readDefinitionFileContent(filePath)
	switch ext {
	case ".yaml", ".yml":
		if strings.Contains(string(content), "openapi:") || strings.Contains(string(content), "swagger:") {
			readOpenApiYaml(content, endpoints)
		} else {
			readYaml(content, endpoints)
		}
	case ".json":
		if strings.Contains(string(content), "\"openapi\":") || strings.Contains(string(content), "\"swagger\":") {
			readOpenApiJson(content, endpoints)
		} else {
			readJson(content, endpoints)
		}
	default:
		log.Fatalf("Unsupported file type: %s", ext)
	}
}

func readDefinitionFileContent(filePath string) (content []byte) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading definition file: %v", err)
	}
	return
}

func readYaml(fileContent []byte, endpoints *[]Endpoint) {
	err := yaml.Unmarshal(fileContent, &endpoints)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML content: %v", err)
	}
	log.Printf("Loaded %d endpoints from the YAML file\n", len(*endpoints))
}

func readJson(fileContent []byte, endpoints *[]Endpoint) {
	err := json.Unmarshal(fileContent, endpoints)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON content: %v", err)
	}
	log.Printf("Loaded %d endpoints from the JSON file\n", len(*endpoints))
}

func readOpenApiYaml(fileContent []byte, endpoints *[]Endpoint) {
	var swagger spec.Swagger
	if err := yaml.Unmarshal(fileContent, &swagger); err != nil {
		log.Fatalf("Error unmarshalling OpenAPI YAML content: %v", err)
	}
	transformOpenApiToCustomFormat(&swagger, endpoints)
}

func readOpenApiJson(fileContent []byte, endpoints *[]Endpoint) {
	var swagger spec.Swagger
	if err := json.Unmarshal(fileContent, &swagger); err != nil {
		log.Fatalf("Error unmarshalling OpenAPI JSON content: %v", err)
	}
	transformOpenApiToCustomFormat(&swagger, endpoints)
}

// transformOpenApiToCustomFormat transforms an OpenAPI spec to quickmock's custom format.
// There are some limitations to this transformation, as OpenAPI spec is more complex than quickmock's.
func transformOpenApiToCustomFormat(swagger *spec.Swagger, endpoints *[]Endpoint) {
	for path, pathItem := range swagger.Paths.Paths {
		// Connect and Trace are not supported by go-openapi spec for some reason
		// TODO: Find out why and fix it in go-openapi?
		methods := map[string]*spec.Operation{
			"GET":     pathItem.Get,
			"POST":    pathItem.Post,
			"PUT":     pathItem.Put,
			"PATCH":   pathItem.Patch,
			"DELETE":  pathItem.Delete,
			"HEAD":    pathItem.Head,
			"OPTIONS": pathItem.Options,
		}

		for method, operation := range methods {
			if operation != nil {
				endpoint := Endpoint{
					Method: method,
					Path:   path,
					Response: Response{
						// Using 200 as the default response code for OpenAPI endpoints
						Code: 200,
						// TODO: Add support for OpenAPI response content
						Content: "Mocked response for OpenAPI",
					},
				}
				*endpoints = append(*endpoints, endpoint)
			}
		}
	}
}
