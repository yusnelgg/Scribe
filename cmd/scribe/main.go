package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yusnelgg/scribe/internal/config"
	"github.com/yusnelgg/scribe/internal/generator"
	"github.com/yusnelgg/scribe/internal/parser"
	"github.com/yusnelgg/scribe/internal/scanner"
)

func main() {
	path := flag.String("path", ".", "Path to the Go project to scan")
	framework := flag.String("framework", "gin", "Web framework to scan (gin, echo, fiber)")
	enableAI := flag.Bool("ai", false, "Enable AI-powered enhancements")
	aiProvider := flag.String("ai-provider", "ollama", "AI provider (ollama, openai, anthropic)")
	aiEndpoint := flag.String("ai-endpoint", "http://localhost:11434", "AI provider endpoint")
	aiModel := flag.String("ai-model", "", "AI model (optional, overrides config)")
	aiKey := flag.String("ai-key", "", "AI API key (optional, overrides config/env)")
	outputDir := flag.String("output", ".", "Output directory for generated files")
	verbose := flag.Bool("v", false, "Verbose output")

	flag.Parse()

	cfg, err := config.LoadDefault()
	if err != nil {
		if *verbose {
			fmt.Fprintf(os.Stderr, "Warning: could not load config: %v\n", err)
		}
	}

	provider := *aiProvider
	endpoint := *aiEndpoint
	model := *aiModel
	apiKey := *aiKey

	if cfg != nil {
		if provider == "ollama" && cfg.GetAIProvider() != "ollama" {
			provider = cfg.GetAIProvider()
		}
		if endpoint == "http://localhost:11434" && cfg.GetAIEndpoint() != "http://localhost:11434" {
			endpoint = cfg.GetAIEndpoint()
		}
		if model == "" {
			model = cfg.GetAIModel()
		}
		if apiKey == "" {
			apiKey = cfg.GetAIAPIKey()
		}
	}

	if *verbose {
		fmt.Printf("Scanning project at: %s\n", *path)
		fmt.Printf("Framework: %s\n", *framework)
		fmt.Printf("AI enabled: %v\n", *enableAI)
		fmt.Printf("AI provider: %s\n", provider)
		if apiKey != "" {
			fmt.Printf("AI API key: configured\n")
		}
	}

	absPath, err := filepath.Abs(*path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		os.Exit(1)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Path does not exist: %s\n", absPath)
		os.Exit(1)
	}

	s := scanner.New()
	files, err := s.Scan(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning project: %v\n", err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("Found %d Go files\n", len(files))
	}

	var p parser.Parser
	switch *framework {
	case "echo":
		p = parser.NewEchoParser()
	case "fiber":
		p = parser.NewFiberParser()
	default:
		p = parser.NewGinParser()
	}
	var routes []parser.Route
	for _, file := range files {
		r, err := p.Parse(file)
		if err != nil {
			if *verbose {
				fmt.Fprintf(os.Stderr, "Warning: error parsing %s: %v\n", file, err)
			}
			continue
		}
		routes = append(routes, r...)
	}

	if *verbose {
		fmt.Printf("Found %d routes\n", len(routes))
	}

	g := generator.New(generator.Options{
		OutputDir:  *outputDir,
		AIEnabled:  *enableAI,
		AIProvider: provider,
		AIEndpoint: endpoint,
		AIModel:    model,
		AIAPIKey:   apiKey,
	})

	if *verbose {
		fmt.Printf("Output directory: %s\n", *outputDir)
		fmt.Printf("Generating files...\n")
	}

	if err := g.Generate(routes, absPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating files: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Generated files:")
	fmt.Printf("  - openapi.json\n")
	fmt.Printf("  - report.md\n")
	fmt.Printf("  - tests/\n")
}
