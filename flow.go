// Package flow provides the official Go SDK for Vylth Flow —
// self-custody crypto payment processing.
//
// Usage:
//
//	client := flow.New("vf_live_...", flow.WithWebhookSecret("whsec_..."))
//
//	invoice, err := client.Invoices.Create(ctx, &flow.CreateInvoiceParams{
//	    Amount:   100,
//	    Currency: "USDT",
//	    Network:  "tron",
//	})
package flow

import "time"

const (
	DefaultBaseURL = "https://flow.vylth.com/api/flow"
	DefaultTimeout = 30 * time.Second
	Version        = "0.1.0"
)

// Client is the Vylth Flow API client.
type Client struct {
	Invoices *InvoiceService
	Payouts  *PayoutService
	Wallets  *WalletService
	Swaps    *SwapService
	Webhooks *WebhookService

	http *httpClient
}

// Option configures the client.
type Option func(*Client)

// WithBaseURL sets a custom API base URL.
func WithBaseURL(url string) Option {
	return func(c *Client) { c.http.baseURL = url }
}

// WithTimeout sets the HTTP request timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.http.timeout = d }
}

// WithWebhookSecret sets the webhook signing secret for signature verification.
func WithWebhookSecret(secret string) Option {
	return func(c *Client) { c.Webhooks = NewWebhookService(secret) }
}

// New creates a new Vylth Flow client.
func New(apiKey string, opts ...Option) *Client {
	h := newHTTPClient(apiKey, DefaultBaseURL, DefaultTimeout)

	c := &Client{
		http:     h,
		Webhooks: NewWebhookService(""),
	}

	for _, opt := range opts {
		opt(c)
	}

	c.Invoices = &InvoiceService{http: h}
	c.Payouts = &PayoutService{http: h}
	c.Wallets = &WalletService{http: h}
	c.Swaps = &SwapService{http: h}

	return c
}
