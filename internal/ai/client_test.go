package ai

import "testing"

func TestNewClient_AllProviders(t *testing.T) {
	providers := []string{"ollama", "openai", "claude", "opencode", "unknown"}

	for _, provider := range providers {
		t.Run(provider, func(t *testing.T) {
			client := NewClient(provider, "http://localhost:8080")
			if client == nil {
				t.Errorf("NewClient returned nil for provider: %s", provider)
			}
		})
	}
}

func TestOllamaClient_Generate(t *testing.T) {
	client := NewClient("ollama", "http://localhost:11434")

	resp, err := client.Generate("test prompt")
	if err == nil {
		t.Logf("Ollama response (expected failure without server): %s", resp)
	}
}

func TestOpenAIClient_MissingAPIKey(t *testing.T) {
	client := NewClient("openai", "")

	_, err := client.Generate("test")
	if err == nil {
		t.Error("Expected error when OPENAI_API_KEY is not set")
	}
	if err != nil && err.Error() != "OPENAI_API_KEY environment variable is required" {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestClaudeClient_MissingAPIKey(t *testing.T) {
	client := NewClient("claude", "")

	_, err := client.Generate("test")
	if err == nil {
		t.Error("Expected error when ANTHROPIC_API_KEY is not set")
	}
	if err != nil && err.Error() != "ANTHROPIC_API_KEY environment variable is required" {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestOpenCodeClient_MissingAPIKey(t *testing.T) {
	client := NewClient("opencode", "")

	_, err := client.Generate("test")
	if err == nil {
		t.Error("Expected error when OPENCODE_API_KEY is not set")
	}
	if err != nil && err.Error() != "OPENCODE_API_KEY environment variable is required" {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestMockClient(t *testing.T) {
	client := NewClient("unknown", "")

	resp, err := client.Generate("test")
	if err != nil {
		t.Errorf("MockClient should not return error: %v", err)
	}
	expected := "AI generation is disabled. Enable with --ai flag."
	if resp != expected {
		t.Errorf("Expected '%s', got '%s'", expected, resp)
	}
}
