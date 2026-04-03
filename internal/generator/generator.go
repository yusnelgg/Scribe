package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yusnelgg/scribe/internal/ai"
	"github.com/yusnelgg/scribe/internal/formatter"
	"github.com/yusnelgg/scribe/internal/parser"
)

type Options struct {
	OutputDir  string
	AIEnabled  bool
	AIProvider string
	AIEndpoint string
	AIModel    string
	AIAPIKey   string
}

type Generator struct {
	opts      Options
	aiClient  ai.AIClient
	enhancer  *ai.Enhancer
	formatter *formatter.Formatter
}

func New(opts Options) *Generator {
	g := &Generator{
		opts:      opts,
		formatter: formatter.New(),
	}

	if opts.AIEnabled {
		g.aiClient = ai.NewClientWithConfig(opts.AIProvider, opts.AIEndpoint, opts.AIModel, opts.AIAPIKey)
		g.enhancer = ai.NewEnhancer(g.aiClient)
	}

	return g
}

func (g *Generator) Generate(routes []parser.Route, projectPath string) error {
	if err := g.generateOpenAPI(routes); err != nil {
		return fmt.Errorf("failed to generate OpenAPI: %w", err)
	}

	if err := g.generateTests(routes); err != nil {
		return fmt.Errorf("failed to generate tests: %w", err)
	}

	if err := g.generateReport(routes, projectPath); err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	return nil
}

func (g *Generator) generateOpenAPI(routes []parser.Route) error {
	openapi := map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":       "API Documentation",
			"description": "Auto-generated API documentation",
			"version":     "1.0.0",
		},
		"servers": []map[string]string{
			{"url": "http://localhost:8080"},
		},
		"paths": map[string]map[string]interface{}{},
	}

	paths := openapi["paths"].(map[string]map[string]interface{})

	for _, route := range routes {
		if _, exists := paths[route.Path]; !exists {
			paths[route.Path] = map[string]interface{}{}
		}

		method := strings.ToLower(route.Method)
		pathSpec := map[string]interface{}{
			"summary":     fmt.Sprintf("%s %s", route.Method, route.Path),
			"operationId": fmt.Sprintf("%s%s", method, sanitizeOperationID(route.Path)),
			"responses": map[string]interface{}{
				"200": map[string]interface{}{
					"description": "Successful response",
				},
			},
		}

		if g.enhancer != nil && route.Handler != "anonymous" {
			if code, err := ai.ReadHandlerCode(route.FilePath, route.Handler); err == nil {
				if enhancement, err := g.enhancer.EnhanceOpenAPI(route, code); err == nil {
					if enhancement.Summary != "" {
						pathSpec["summary"] = enhancement.Summary
					}
					if enhancement.Description != "" {
						pathSpec["description"] = enhancement.Description
					}
				}
			}
		}

		paths[route.Path][method] = pathSpec
	}

	outputPath := filepath.Join(g.opts.OutputDir, "openapi.json")
	return writeJSON(outputPath, openapi)
}

func (g *Generator) generateTests(routes []parser.Route) error {
	testsDir := filepath.Join(g.opts.OutputDir, "tests")
	if err := os.MkdirAll(testsDir, 0755); err != nil {
		return err
	}

	for _, route := range routes {
		var content string

		if g.enhancer != nil && route.Handler != "anonymous" {
			if code, err := ai.ReadHandlerCode(route.FilePath, route.Handler); err == nil {
				if enhancement, err := g.enhancer.EnhanceTest(route, code); err == nil {
					content = enhancement.TestBody
				}
			}
		}

		if content == "" {
			content = g.formatter.FormatTest(route)
		}

		filename := fmt.Sprintf("test_%s_%s_generated.go", strings.ToLower(route.Method), sanitizeFilename(route.Path))
		outputPath := filepath.Join(testsDir, filename)

		if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) generateReport(routes []parser.Route, projectPath string) error {
	var buf bytes.Buffer

	buf.WriteString("# Scribe Analysis Report\n\n")
	buf.WriteString(fmt.Sprintf("Generated: %s\n\n", time.Now().Format(time.RFC1123)))
	buf.WriteString("## Summary\n\n")
	buf.WriteString(fmt.Sprintf("- **Total Routes**: %d\n", len(routes)))
	buf.WriteString(fmt.Sprintf("- **AI Enhanced**: %v\n", g.opts.AIEnabled))
	buf.WriteString(fmt.Sprintf("- **Project Path**: %s\n\n", projectPath))

	methodCount := make(map[string]int)
	for _, route := range routes {
		methodCount[route.Method]++
	}

	buf.WriteString("## Routes by Method\n\n")
	buf.WriteString("| Method | Count |\n")
	buf.WriteString("|--------|-------|\n")
	for method, count := range methodCount {
		buf.WriteString(fmt.Sprintf("| %s | %d |\n", method, count))
	}
	buf.WriteString("\n")

	buf.WriteString("## Detailed Routes\n\n")
	buf.WriteString("| Method | Path | Handler | File |\n")
	buf.WriteString("|--------|------|---------|------|\n")
	for _, route := range routes {
		filename := filepath.Base(route.FilePath)
		buf.WriteString(fmt.Sprintf("| %s | %s | %s | %s:%d |\n",
			route.Method, route.Path, route.Handler, filename, route.Line))
	}

	outputPath := filepath.Join(g.opts.OutputDir, "report.md")
	return os.WriteFile(outputPath, buf.Bytes(), 0644)
}

func writeJSON(path string, data interface{}) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func sanitizeOperationID(path string) string {
	path = strings.TrimPrefix(path, "/")
	path = strings.ReplaceAll(path, "/", "_")
	path = strings.ReplaceAll(path, "-", "_")
	return path
}

func sanitizeFilename(path string) string {
	path = strings.TrimPrefix(path, "/")
	path = strings.ReplaceAll(path, "/", "_")
	path = strings.ReplaceAll(path, "-", "_")
	path = strings.ReplaceAll(path, ":", "_")
	return path
}
