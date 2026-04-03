package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type AIClient interface {
	Generate(prompt string) (string, error)
	GenerateWithSystem(systemPrompt, userPrompt string) (string, error)
}

func NewClient(provider, endpoint string) AIClient {
	return NewClientWithConfig(provider, endpoint, "", "")
}

func NewClientWithConfig(provider, endpoint, model, apiKey string) AIClient {
	switch provider {
	case "ollama":
		if model == "" {
			model = "llama2"
		}
		return &OllamaClient{
			endpoint: endpoint,
			model:    model,
			client:   &http.Client{},
		}
	case "openai":
		if model == "" {
			model = "gpt-4o-mini"
		}
		if apiKey == "" {
			apiKey = os.Getenv("OPENAI_API_KEY")
		}
		return &OpenAIClient{
			endpoint: endpoint,
			model:    model,
			apiKey:   apiKey,
			client:   &http.Client{},
		}
	case "claude":
		if model == "" {
			model = "claude-sonnet-4-20250514"
		}
		if apiKey == "" {
			apiKey = os.Getenv("ANTHROPIC_API_KEY")
		}
		return &ClaudeClient{
			model:  model,
			apiKey: apiKey,
			client: &http.Client{},
		}
	case "opencode":
		if model == "" {
			model = "qwen3-8b"
		}
		if apiKey == "" {
			apiKey = os.Getenv("OPENCODE_API_KEY")
		}
		return &OpenCodeClient{
			model:  model,
			apiKey: apiKey,
			client: &http.Client{},
		}
	default:
		return &MockClient{}
	}
}

type OllamaClient struct {
	endpoint string
	model    string
	client   *http.Client
}

func (c *OllamaClient) Generate(prompt string) (string, error) {
	return c.GenerateWithSystem("", prompt)
}

func (c *OllamaClient) GenerateWithSystem(systemPrompt, userPrompt string) (string, error) {
	url := fmt.Sprintf("%s/api/generate", c.endpoint)

	requestBody := map[string]interface{}{
		"model":  c.model,
		"prompt": userPrompt,
		"stream": false,
	}

	if systemPrompt != "" {
		requestBody["system"] = systemPrompt
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama returned status %d", resp.StatusCode)
	}

	var result struct {
		Response string `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Response, nil
}

type OpenAIClient struct {
	endpoint string
	model    string
	apiKey   string
	client   *http.Client
}

func (c *OpenAIClient) Generate(prompt string) (string, error) {
	return c.GenerateWithSystem("", prompt)
}

func (c *OpenAIClient) GenerateWithSystem(systemPrompt, userPrompt string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"
	if c.endpoint != "" {
		url = c.endpoint + "/v1/chat/completions"
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY environment variable is required")
	}

	messages := []map[string]string{}
	if systemPrompt != "" {
		messages = append(messages, map[string]string{
			"role":    "system",
			"content": systemPrompt,
		})
	}
	messages = append(messages, map[string]string{
		"role":    "user",
		"content": userPrompt,
	})

	requestBody := map[string]interface{}{
		"model":    c.model,
		"messages": messages,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call OpenAI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI returned status %d", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return result.Choices[0].Message.Content, nil
}

type ClaudeClient struct {
	model  string
	apiKey string
	client *http.Client
}

func (c *ClaudeClient) Generate(prompt string) (string, error) {
	return c.GenerateWithSystem("", prompt)
}

func (c *ClaudeClient) GenerateWithSystem(systemPrompt, userPrompt string) (string, error) {
	url := "https://api.anthropic.com/v1/messages"

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("ANTHROPIC_API_KEY environment variable is required")
	}

	system := systemPrompt
	if system == "" {
		system = "You are a helpful assistant."
	}

	requestBody := map[string]interface{}{
		"model":  c.model,
		"system": system,
		"messages": []map[string]string{
			{"role": "user", "content": userPrompt},
		},
		"max_tokens": 4096,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Claude: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Claude returned status %d", resp.StatusCode)
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Content) == 0 {
		return "", fmt.Errorf("no response from Claude")
	}

	return result.Content[0].Text, nil
}

type OpenCodeClient struct {
	model  string
	apiKey string
	client *http.Client
}

func (c *OpenCodeClient) Generate(prompt string) (string, error) {
	return c.GenerateWithSystem("", prompt)
}

func (c *OpenCodeClient) GenerateWithSystem(systemPrompt, userPrompt string) (string, error) {
	url := "https://opencode.ai/api/chat/completions"

	apiKey := os.Getenv("OPENCODE_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENCODE_API_KEY environment variable is required")
	}

	messages := []map[string]string{}
	if systemPrompt != "" {
		messages = append(messages, map[string]string{
			"role":    "system",
			"content": systemPrompt,
		})
	}
	messages = append(messages, map[string]string{
		"role":    "user",
		"content": userPrompt,
	})

	requestBody := map[string]interface{}{
		"model":    c.model,
		"messages": messages,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call OpenCode: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenCode returned status %d", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenCode")
	}

	return result.Choices[0].Message.Content, nil
}

type MockClient struct{}

func (c *MockClient) Generate(prompt string) (string, error) {
	return "AI generation is disabled. Enable with --ai flag.", nil
}

func (c *MockClient) GenerateWithSystem(systemPrompt, userPrompt string) (string, error) {
	return "AI generation is disabled. Enable with --ai flag.", nil
}
