package hooksniff

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	defaultBaseURL = "https://hooksniff-api-499907444852.europe-west1.run.app"
	defaultTimeout = 30 * time.Second
	defaultRetries = 3
	sdkVersion     = "0.5.0"
)

// ClientConfig holds optional configuration for the HookSniff client.
type ClientConfig struct {
	BaseURL string
	Timeout time.Duration
	Retries int
	Headers map[string]string
}

// httpClient handles HTTP communication with the HookSniff API.
type httpClient struct {
	apiKey       string
	baseURL      string
	client       *http.Client
	retries      int
	extraHeaders map[string]string
}

func newHTTPClient(apiKey string, cfg *ClientConfig) *httpClient {
	if apiKey == "" {
		panic("HookSniff API key is required")
	}

	baseURL := defaultBaseURL
	timeout := defaultTimeout
	retries := defaultRetries
	headers := make(map[string]string)

	if cfg != nil {
		if cfg.BaseURL != "" {
			baseURL = cfg.BaseURL
		}
		if cfg.Timeout > 0 {
			timeout = cfg.Timeout
		}
		if cfg.Retries > 0 {
			retries = cfg.Retries
		}
		if cfg.Headers != nil {
			headers = cfg.Headers
		}
	}

	return &httpClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: timeout,
		},
		retries:      retries,
		extraHeaders: headers,
	}
}

func (c *httpClient) request(method, path string, body interface{}, opts ...requestOption) (json.RawMessage, error) {
	var options requestOptions
	for _, opt := range opts {
		opt(&options)
	}

	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBytes)
	}

	var lastErr error

	for attempt := 0; attempt <= c.retries; attempt++ {
		req, err := http.NewRequest(method, url, bodyReader)
		if err != nil {
			return nil, fmt.Errorf("create request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "hooksniff-sdk/go/"+sdkVersion)
		for k, v := range c.extraHeaders {
			req.Header.Set(k, v)
		}
		if options.idempotencyKey != "" {
			req.Header.Set("Idempotency-Key", options.idempotencyKey)
		}

		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = err
			if attempt < c.retries {
				time.Sleep(backoff(attempt))
				continue
			}
			return nil, fmt.Errorf("request failed after %d retries: %w", c.retries, err)
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			if attempt < c.retries {
				time.Sleep(backoff(attempt))
				continue
			}
			return nil, fmt.Errorf("read response: %w", err)
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if len(respBody) == 0 {
				return nil, nil
			}
			return respBody, nil
		}

		// Parse error body
		var errBody map[string]interface{}
		json.Unmarshal(respBody, &errBody)

		// Don't retry client errors (except 429 and 408)
		if resp.StatusCode >= 400 && resp.StatusCode < 500 && resp.StatusCode != 429 && resp.StatusCode != 408 {
			return nil, mapError(resp.StatusCode, errBody)
		}

		// Handle rate limiting
		if resp.StatusCode == 429 {
			retryAfter := 60
			if ra := resp.Header.Get("Retry-After"); ra != "" {
				if v, err := strconv.Atoi(ra); err == nil {
					retryAfter = v
				}
			}
			if attempt < c.retries {
				time.Sleep(time.Duration(retryAfter) * time.Second)
				continue
			}
			err := mapError(429, errBody)
			if re, ok := err.(*RateLimitError); ok {
				re.RetryAfter = retryAfter
			}
			return nil, err
		}

		// Server errors — retry with backoff
		lastErr = mapError(resp.StatusCode, errBody)
		if attempt < c.retries {
			time.Sleep(backoff(attempt))
			continue
		}
	}

	return nil, lastErr
}

type requestOptions struct {
	idempotencyKey string
}

type requestOption func(*requestOptions)

// WithIdempotencyKey sets the Idempotency-Key header.
func WithIdempotencyKey(key string) requestOption {
	return func(o *requestOptions) {
		o.idempotencyKey = key
	}
}

func backoff(attempt int) time.Duration {
	base := math.Pow(2, float64(attempt)) * 1000
	jitter := rand.Float64() * 1000
	return time.Duration(base+jitter) * time.Millisecond
}
