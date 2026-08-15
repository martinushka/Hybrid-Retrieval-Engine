package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPProvider struct {
	baseURL string
	client  *http.Client
}

func NewHTTPProvider(baseURL string) *HTTPProvider {
	return &HTTPProvider{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

type embedRequest struct {
	Text string `json:"text"`
}

type embedResponse struct {
	Embedding []float32 `json:"embedding"`
}

func (p *HTTPProvider) Embed(ctx context.Context, text string) ([]float32, error) {
	body, err := json.Marshal(embedRequest{
		Text: text,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal embedding request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		p.baseURL+"/embed",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("create embedding request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request embedding service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("embedding service returned status %d", resp.StatusCode)
	}

	var result embedResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode embedding response: %w", err)
	}

	if len(result.Embedding) != 384 {
		return nil, fmt.Errorf(
			"invalid embedding dimension: got %d, want 384",
			len(result.Embedding),
		)
	}

	return result.Embedding, nil
}
