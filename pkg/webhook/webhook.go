package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"fleetify/pkg/errors"
)

type WebhookPayload struct {
	Event     string      `json:"event"`
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data"`
}

type Client struct {
	URL     string
	Timeout time.Duration
}

func NewClient(url string) *Client {
	return &Client{
		URL:     url,
		Timeout: 5 * time.Second,
	}
}

func (c *Client) Send(ctx context.Context, event string, data interface{}) error {
	if c.URL == "" {
		return nil
	}

	payload := WebhookPayload{
		Event:     event,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      data,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	reqCtx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, "POST", c.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Fleetify-Webhook/1.0")

	client := &http.Client{
		Timeout: c.Timeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) SendAsync(ctx context.Context, event string, data interface{}) {
	go func() {
		if err := c.Send(ctx, event, data); err != nil {
			errors.LogError("Webhook send error", err)
		}
	}()
}

