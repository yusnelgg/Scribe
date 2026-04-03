package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/yusnel/mt-scribe/internal/parser"
)

type Enhancer struct {
	client  AIClient
	timeout time.Duration
}

func NewEnhancer(client AIClient) *Enhancer {
	return &Enhancer{
		client:  client,
		timeout: 60 * time.Second,
	}
}

type TestEnhancement struct {
	PackageImports []string
	TestBody       string
	SetupCode      string
	Assertions     []string
}

func (e *Enhancer) EnhanceTest(route parser.Route, handlerCode string) (*TestEnhancement, error) {
	prompt := buildTestPrompt(route, handlerCode)

	resp, err := e.client.GenerateWithSystem(systemPromptTest, prompt)
	if err != nil {
		return nil, fmt.Errorf("AI enhancement failed: %w", err)
	}

	return parseTestResponse(resp)
}

const systemPromptTest = `You are a Go testing expert. Generate test code following these rules:
- Use table-driven tests when appropriate
- Include proper error handling assertions
- Mock external dependencies (DB, HTTP clients)
- Use testify/require for assertions
- Return only valid Go code wrapped in markdown code blocks
- Include realistic test data
- Include setup code for router configuration`

func buildTestPrompt(route parser.Route, handlerCode string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Generate a test for this Gin endpoint:\n"))
	sb.WriteString(fmt.Sprintf("- Method: %s\n", route.Method))
	sb.WriteString(fmt.Sprintf("- Path: %s\n", route.Path))
	sb.WriteString(fmt.Sprintf("- Handler: %s\n\n", route.Handler))
	sb.WriteString("Handler code:\n```go\n")
	sb.WriteString(handlerCode)
	sb.WriteString("\n```\n\n")
	sb.WriteString("Generate a complete, runnable test file with router setup.")
	return sb.String()
}

func parseTestResponse(resp string) (*TestEnhancement, error) {
	code := extractCodeBlock(resp)
	if code == "" {
		code = resp
	}

	enh := &TestEnhancement{
		PackageImports: []string{
			`"testing"`,
			`"net/http"`,
			`"net/http/httptest"`,
			`"github.com/gin-gonic/gin"`,
			`"github.com/stretchr/testify/assert"`,
		},
		TestBody: code,
	}

	return enh, nil
}

func extractCodeBlock(s string) string {
	lines := strings.Split(s, "\n")
	var inBlock bool
	var result []string

	for _, line := range lines {
		if strings.Contains(line, "```go") || strings.Contains(line, "```") {
			if !inBlock {
				inBlock = true
				continue
			}
			break
		}
		if inBlock {
			result = append(result, line)
		}
	}

	return strings.TrimSpace(strings.Join(result, "\n"))
}

type OpenAPIEnhancement struct {
	Summary     string             `json:"summary"`
	Description string             `json:"description"`
	RequestBody *RequestBodySchema `json:"requestBody,omitempty"`
	Responses   map[string]string  `json:"responses"`
}

type RequestBodySchema struct {
	Required   bool              `json:"required"`
	Properties map[string]string `json:"properties"`
	Example    string            `json:"example"`
}

func (e *Enhancer) EnhanceOpenAPI(route parser.Route, handlerCode string) (*OpenAPIEnhancement, error) {
	prompt := buildOpenAPIPrompt(route, handlerCode)

	resp, err := e.client.GenerateWithSystem(systemPromptOpenAPI, prompt)
	if err != nil {
		return nil, fmt.Errorf("AI enhancement failed: %w", err)
	}

	return parseOpenAPIResponse(resp, route)
}

const systemPromptOpenAPI = `You are an OpenAPI 3.0 specification expert. Return ONLY valid JSON with this structure:
{
  "summary": "short description (max 50 chars)",
  "description": "detailed description of what this endpoint does",
  "requestBody": { "required": true/false, "properties": {"field": "type"}, "example": "json example" },
  "responses": { "200": "description", "400": "description", "500": "description" }
}`

func buildOpenAPIPrompt(route parser.Route, handlerCode string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Analyze this Gin handler and generate OpenAPI 3.0 enhancement:\n"))
	sb.WriteString(fmt.Sprintf("- Method: %s\n", route.Method))
	sb.WriteString(fmt.Sprintf("- Path: %s\n", route.Path))
	sb.WriteString(fmt.Sprintf("- Handler: %s\n\n", route.Handler))
	sb.WriteString("Handler code:\n```go\n")
	sb.WriteString(handlerCode)
	sb.WriteString("\n```")
	return sb.String()
}

func parseOpenAPIResponse(resp string, route parser.Route) (*OpenAPIEnhancement, error) {
	resp = strings.TrimSpace(resp)
	resp = strings.TrimPrefix(resp, "```json")
	resp = strings.TrimPrefix(resp, "```")
	resp = strings.TrimSuffix(resp, "```")
	resp = strings.TrimSpace(resp)

	var result OpenAPIEnhancement
	if err := json.Unmarshal([]byte(resp), &result); err != nil {
		return &OpenAPIEnhancement{
			Summary:     fmt.Sprintf("%s %s", route.Method, route.Path),
			Description: "AI-generated description",
			Responses:   map[string]string{"200": "Successful response"},
		}, nil
	}

	if result.Responses == nil {
		result.Responses = map[string]string{
			"200": "Successful response",
		}
	}

	return &result, nil
}

func ReadHandlerCode(filePath string, handlerName string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(content), "\n")

	for i, line := range lines {
		if strings.Contains(line, "func "+handlerName) || strings.Contains(line, "func (") {
			start := i
			var end int
			braceCount := 0
			started := false

			for j := i; j < len(lines); j++ {
				for _, ch := range lines[j] {
					if ch == '{' {
						braceCount++
						started = true
					} else if ch == '}' {
						braceCount--
					}
				}
				if started && braceCount == 0 {
					end = j + 1
					break
				}
			}

			if end > start {
				return strings.Join(lines[start:end], "\n"), nil
			}
		}
	}

	return "", fmt.Errorf("handler %s not found in %s", handlerName, filePath)
}
