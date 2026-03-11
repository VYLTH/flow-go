package flow

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// WebhookService handles webhook signature verification.
type WebhookService struct {
	secret string
}

// NewWebhookService creates a new WebhookService with the given signing secret.
func NewWebhookService(secret string) *WebhookService {
	return &WebhookService{secret: secret}
}

// Verify checks the webhook signature and returns the parsed event.
func (s *WebhookService) Verify(payload []byte, signature string) (*WebhookEvent, error) {
	if s.secret == "" {
		return nil, &InvalidSignatureError{Message: "webhook secret not configured"}
	}

	mac := hmac.New(sha256.New, []byte(s.secret))
	mac.Write(payload)
	expected := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expected), []byte(signature)) {
		return nil, &InvalidSignatureError{}
	}

	var event WebhookEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("flow: unmarshal webhook event: %w", err)
	}
	return &event, nil
}

// IsValid checks if a webhook signature is valid without parsing the event.
func (s *WebhookService) IsValid(payload []byte, signature string) bool {
	if s.secret == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(s.secret))
	mac.Write(payload)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

// VerifyRequest reads and verifies a webhook from an http.Request.
// It expects the signature in the X-Flow-Signature header.
func (s *WebhookService) VerifyRequest(r *http.Request) (*WebhookEvent, error) {
	signature := r.Header.Get("X-Flow-Signature")
	if signature == "" {
		return nil, &InvalidSignatureError{Message: "missing X-Flow-Signature header"}
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("flow: read request body: %w", err)
	}
	defer r.Body.Close()

	return s.Verify(body, signature)
}
