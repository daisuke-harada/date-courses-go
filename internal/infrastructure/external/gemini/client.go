package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GeminiClient は Gemini API クライアントのインターフェースです。
type GeminiClient interface {
	GenerateContent(ctx context.Context, prompt string) (string, error)
}

type geminiClient struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewGeminiClient は Gemini REST API クライアントを生成します。
// standard net/http を使用し、外部 SDK への依存を持ちません。
func NewGeminiClient(apiKey, model string) GeminiClient {
	return &geminiClient{
		apiKey:     apiKey,
		model:      model,
		httpClient: &http.Client{},
	}
}

// generateContentRequest は Gemini API のリクエストボディです。
type generateContentRequest struct {
	Contents []contentItem `json:"contents"`
}

type contentItem struct {
	Parts []part `json:"parts"`
}

type part struct {
	Text string `json:"text"`
}

// generateContentResponse は Gemini API のレスポンスボディです。
type generateContentResponse struct {
	Candidates []candidate `json:"candidates"`
}

type candidate struct {
	Content contentItem `json:"content"`
}

// GenerateContent はプロンプトを Gemini API に送信し、生成されたテキストを返します。
func (c *geminiClient) GenerateContent(ctx context.Context, prompt string) (string, error) {
	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		c.model, c.apiKey,
	)

	reqBody := generateContentRequest{
		Contents: []contentItem{
			{Parts: []part{{Text: prompt}}},
		},
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("gemini: failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return "", fmt.Errorf("gemini: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("gemini: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gemini: API returned status %d", resp.StatusCode)
	}

	var result generateContentResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("gemini: failed to decode response: %w", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini: empty response")
	}

	return strings.TrimSpace(result.Candidates[0].Content.Parts[0].Text), nil
}
