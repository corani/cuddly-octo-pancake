package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	defaultBaseURL = "https://models.github.ai"
	apiVersion     = "2022-11-28"
)

// Client is the GitHub Models API client.
type Client struct {
	token      string
	model      string
	baseURL    string
	httpClient *http.Client
}

// Message represents a message in a chat conversation.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest is a request to create a chat completion.
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
}

// Choice represents a completion choice returned by the API.
type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// ChatResponse is the response from a chat completion request.
type ChatResponse struct {
	Choices []Choice `json:"choices"`
}

type ModelResponse struct {
	ID                        string   `json:"id"`
	Name                      string   `json:"name"`
	Publisher                 string   `json:"publisher"`
	Summary                   string   `json:"summary"`
	RateLimitTier             string   `json:"rate_limit_tier"`
	SupportedInputModalities  []string `json:"supported_input_modalities"`
	SupportedOutputModalities []string `json:"supported_output_modalities"`
	Tags                      []string `json:"tags"`
}

// StatusError is an error returned by the GitHub Models API.
type StatusError struct {
	StatusCode   int    `json:"status_code"`
	ErrorMessage string `json:"error"`
}

func (e StatusError) Error() string {
	return fmt.Sprintf("GitHub Models API error: %s (status code: %d)", e.ErrorMessage, e.StatusCode)
}

// NewClient creates a new GitHub Models API client.
func NewClient(token, model string, httpClient *http.Client) (*Client, error) {
	if token == "" {
		return nil, fmt.Errorf("missing GitHub token")
	}

	if model == "" {
		model = "openai/gpt-4.1"
	}

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		token:      token,
		model:      model,
		baseURL:    defaultBaseURL,
		httpClient: httpClient,
	}, nil
}

// CreateChat creates a chat completion.
func (c *Client) CreateChat(ctx context.Context, reqData *ChatRequest) (*ChatResponse, error) {
	if reqData.Model == "" {
		reqData.Model = c.model
	}

	url := fmt.Sprintf("%s/inference/chat/completions", c.baseURL)

	var reqBody io.Reader
	if reqData != nil {
		jsonData, err := json.Marshal(reqData)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("X-GitHub-Api-Version", apiVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		apiErr := StatusError{
			StatusCode: resp.StatusCode,
		}

		// Try to unmarshal error response
		if err := json.Unmarshal(respBody, &apiErr); err != nil {
			apiErr.ErrorMessage = string(respBody)
		}

		return nil, apiErr
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, err
	}

	return &chatResp, nil
}

func (c *Client) ListModels(ctx context.Context) ([]ModelResponse, error) {
	url := fmt.Sprintf("%s/catalog/models", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		apiErr := StatusError{
			StatusCode: resp.StatusCode,
		}

		// Try to unmarshal error response
		if err := json.Unmarshal(respBody, &apiErr); err != nil {
			apiErr.ErrorMessage = string(respBody)
		}

		return nil, apiErr
	}

	var modelsResp []ModelResponse
	if err := json.Unmarshal(respBody, &modelsResp); err != nil {
		return nil, err
	}

	return modelsResp, nil
}
