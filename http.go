package flow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const maxRetries = 3

type httpClient struct {
	baseURL string
	apiKey  string
	timeout time.Duration
	client  *http.Client
}

func newHTTPClient(apiKey, baseURL string, timeout time.Duration) *httpClient {
	return &httpClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		timeout: timeout,
		client:  &http.Client{Timeout: timeout},
	}
}

func (h *httpClient) do(ctx context.Context, method, path string, body interface{}, params map[string]string) ([]byte, error) {
	fullURL := h.baseURL + path

	if len(params) > 0 {
		q := url.Values{}
		for k, v := range params {
			if v != "" {
				q.Set(k, v)
			}
		}
		if qs := q.Encode(); qs != "" {
			fullURL += "?" + qs
		}
	}

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("flow: marshal error: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
		if err != nil {
			return nil, fmt.Errorf("flow: request error: %w", err)
		}
		req.Header.Set("X-API-Key", h.apiKey)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "flow-go/"+Version)

		// Reset body reader for retries
		if body != nil {
			data, _ := json.Marshal(body)
			bodyReader = bytes.NewReader(data)
			req.Body = io.NopCloser(bodyReader)
		}

		resp, err := h.client.Do(req)
		if err != nil {
			lastErr = err
			if attempt < maxRetries-1 {
				time.Sleep(time.Duration(1<<attempt) * 500 * time.Millisecond)
				continue
			}
			return nil, &FlowError{Message: err.Error(), Code: "network_error"}
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("flow: read body error: %w", err)
		}

		if shouldRetry(resp.StatusCode) && attempt < maxRetries-1 {
			wait := time.Duration(1<<attempt) * 500 * time.Millisecond
			if resp.StatusCode == 429 {
				if ra := resp.Header.Get("Retry-After"); ra != "" {
					if secs, err := strconv.ParseFloat(ra, 64); err == nil {
						wait = time.Duration(secs * float64(time.Second))
					}
				}
			}
			time.Sleep(wait)
			continue
		}

		if resp.StatusCode >= 400 {
			return nil, parseError(resp.StatusCode, respBody, resp.Header)
		}

		if resp.StatusCode == 204 {
			return nil, nil
		}

		return respBody, nil
	}

	if lastErr != nil {
		return nil, &FlowError{Message: lastErr.Error(), Code: "network_error"}
	}
	return nil, &FlowError{Message: "max retries exceeded", Code: "network_error"}
}

func (h *httpClient) get(ctx context.Context, path string, params map[string]string) ([]byte, error) {
	return h.do(ctx, http.MethodGet, path, nil, params)
}

func (h *httpClient) post(ctx context.Context, path string, body interface{}) ([]byte, error) {
	return h.do(ctx, http.MethodPost, path, body, nil)
}

func shouldRetry(status int) bool {
	return status == 429 || status == 500 || status == 502 || status == 503 || status == 504
}

func parseError(status int, body []byte, headers http.Header) error {
	var errBody struct {
		Detail  string              `json:"detail"`
		Message string              `json:"message"`
		Error   string              `json:"error"`
		Code    string              `json:"code"`
		Errors  []map[string]string `json:"errors"`
	}
	_ = json.Unmarshal(body, &errBody)

	msg := errBody.Detail
	if msg == "" {
		msg = errBody.Message
	}
	if msg == "" {
		msg = errBody.Error
	}
	if msg == "" {
		msg = string(body)
	}

	base := FlowError{Message: msg, Code: errBody.Code, Status: status}

	switch status {
	case 401:
		return &AuthenticationError{FlowError: base}
	case 404:
		return &NotFoundError{FlowError: base}
	case 422:
		return &ValidationError{FlowError: base, Errors: errBody.Errors}
	case 429:
		var retryAfter float64
		if ra := headers.Get("Retry-After"); ra != "" {
			retryAfter, _ = strconv.ParseFloat(ra, 64)
		}
		return &RateLimitError{FlowError: base, RetryAfter: retryAfter}
	default:
		if status >= 500 {
			return &ServerError{FlowError: base}
		}
		return &base
	}
}
