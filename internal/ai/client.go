package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AIClient interface {
	Generate(prompt string) (string, error)
	GenerateWithSystem(systemPrompt, userPrompt string) (string, error)
}

type OllamaClient struct {
	endpoint string
	client   *http.Client
}

func NewClient(provider, endpoint string) AIClient {
	switch provider {
	case "ollama":
		return &OllamaClient{
			endpoint: endpoint,
			client:   &http.Client{},
		}
	default:
		return &MockClient{}
	}
}

func (c *OllamaClient) Generate(prompt string) (string, error) {
	return c.GenerateWithSystem("", prompt)
}

func (c *OllamaClient) GenerateWithSystem(systemPrompt, userPrompt string) (string, error) {
	url := fmt.Sprintf("%s/api/generate", c.endpoint)

	requestBody := map[string]interface{}{
		"model":  "llama2",
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

type MockClient struct{}

func (c *MockClient) Generate(prompt string) (string, error) {
	return "AI generation is disabled. Enable with --ai flag.", nil
}

func (c *MockClient) GenerateWithSystem(systemPrompt, userPrompt string) (string, error) {
	return "AI generation is disabled. Enable with --ai flag.", nil
}
