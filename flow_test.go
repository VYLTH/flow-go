package flow

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestServer(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	ts := httptest.NewServer(handler)
	client := New("test_key", WithBaseURL(ts.URL))
	return client, ts
}

func TestNewClient(t *testing.T) {
	c := New("vf_live_test123")
	if c.Invoices == nil || c.Payouts == nil || c.Wallets == nil || c.Swaps == nil || c.Webhooks == nil {
		t.Fatal("expected all services to be initialized")
	}
}

func TestNewClientWithOptions(t *testing.T) {
	c := New("vf_live_test123",
		WithBaseURL("https://custom.api.com"),
		WithWebhookSecret("whsec_test"),
	)
	if c.http.baseURL != "https://custom.api.com" {
		t.Fatalf("expected custom base URL, got %s", c.http.baseURL)
	}
}

func TestCreateInvoice(t *testing.T) {
	client, ts := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/invoices" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("X-API-Key") != "test_key" {
			t.Fatal("missing API key header")
		}
		json.NewEncoder(w).Encode(Invoice{
			ID: "inv_123", Amount: 100, Currency: "USDT", Network: "tron", Status: "pending",
		})
	})
	defer ts.Close()

	inv, err := client.Invoices.Create(context.Background(), &CreateInvoiceParams{
		Amount: 100, Currency: "USDT", Network: "tron",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inv.ID != "inv_123" {
		t.Fatalf("expected inv_123, got %s", inv.ID)
	}
}

func TestGetInvoice(t *testing.T) {
	client, ts := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/invoices/inv_123" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(Invoice{ID: "inv_123", Status: "paid"})
	})
	defer ts.Close()

	inv, err := client.Invoices.Get(context.Background(), "inv_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inv.Status != "paid" {
		t.Fatalf("expected paid, got %s", inv.Status)
	}
}

func TestListInvoices(t *testing.T) {
	client, ts := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "1" || r.URL.Query().Get("limit") != "10" {
			t.Fatal("missing query params")
		}
		json.NewEncoder(w).Encode(PaginatedList[Invoice]{
			Items: []Invoice{{ID: "inv_1"}, {ID: "inv_2"}}, Total: 2, Page: 1, PerPage: 10,
		})
	})
	defer ts.Close()

	list, err := client.Invoices.List(context.Background(), &ListParams{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(list.Items))
	}
}

func TestCreatePayout(t *testing.T) {
	client, ts := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(Payout{
			ID: "po_123", Amount: 50, Currency: "USDT", Status: "pending",
		})
	})
	defer ts.Close()

	po, err := client.Payouts.Create(context.Background(), &CreatePayoutParams{
		Amount: 50, Currency: "USDT", Network: "tron", Destination: "T...",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if po.ID != "po_123" {
		t.Fatalf("expected po_123, got %s", po.ID)
	}
}

func TestGenerateWallet(t *testing.T) {
	client, ts := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(Wallet{
			ID: "w_123", Address: "T...", Network: "tron", Currency: "USDT",
		})
	})
	defer ts.Close()

	w, err := client.Wallets.Generate(context.Background(), &GenerateWalletParams{
		Network: "tron", Currency: "USDT",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.ID != "w_123" {
		t.Fatalf("expected w_123, got %s", w.ID)
	}
}

func TestSwapQuote(t *testing.T) {
	client, ts := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(SwapQuote{
			FromCurrency: "USDT", ToCurrency: "BTC", FromAmount: 1000, ToAmount: 0.015, Rate: 0.000015,
		})
	})
	defer ts.Close()

	q, err := client.Swaps.Quote(context.Background(), &SwapParams{
		FromCurrency: "USDT", ToCurrency: "BTC", Amount: 1000,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.Rate != 0.000015 {
		t.Fatalf("expected rate 0.000015, got %f", q.Rate)
	}
}

func TestWebhookVerify(t *testing.T) {
	secret := "whsec_test_secret"
	payload := []byte(`{"id":"evt_123","type":"invoice.paid","data":{"invoice_id":"inv_123"}}`)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	signature := hex.EncodeToString(mac.Sum(nil))

	client := New("test_key", WithWebhookSecret(secret))

	event, err := client.Webhooks.Verify(payload, signature)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event.ID != "evt_123" || event.Type != "invoice.paid" {
		t.Fatalf("unexpected event: %+v", event)
	}
}

func TestWebhookInvalidSignature(t *testing.T) {
	client := New("test_key", WithWebhookSecret("whsec_secret"))
	_, err := client.Webhooks.Verify([]byte(`{}`), "bad_signature")
	if err == nil {
		t.Fatal("expected error for invalid signature")
	}
	if _, ok := err.(*InvalidSignatureError); !ok {
		t.Fatalf("expected InvalidSignatureError, got %T", err)
	}
}

func TestWebhookIsValid(t *testing.T) {
	secret := "whsec_test"
	payload := []byte(`{"test":true}`)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	sig := hex.EncodeToString(mac.Sum(nil))

	client := New("test_key", WithWebhookSecret(secret))
	if !client.Webhooks.IsValid(payload, sig) {
		t.Fatal("expected valid signature")
	}
	if client.Webhooks.IsValid(payload, "wrong") {
		t.Fatal("expected invalid signature")
	}
}

func TestErrorParsing(t *testing.T) {
	client, ts := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid API key"})
	})
	defer ts.Close()

	_, err := client.Invoices.Get(context.Background(), "inv_123")
	if err == nil {
		t.Fatal("expected error")
	}
	if _, ok := err.(*AuthenticationError); !ok {
		t.Fatalf("expected AuthenticationError, got %T: %v", err, err)
	}
}

func TestNotFoundError(t *testing.T) {
	client, ts := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]string{"message": "Not found"})
	})
	defer ts.Close()

	_, err := client.Invoices.Get(context.Background(), "inv_nonexist")
	if err == nil {
		t.Fatal("expected error")
	}
	if _, ok := err.(*NotFoundError); !ok {
		t.Fatalf("expected NotFoundError, got %T", err)
	}
}

func TestRateLimitError(t *testing.T) {
	attempt := 0
	client, ts := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		attempt++
		w.Header().Set("Retry-After", "0.1")
		w.WriteHeader(429)
		json.NewEncoder(w).Encode(map[string]string{"message": "Too many requests"})
	})
	defer ts.Close()

	_, err := client.Invoices.Get(context.Background(), "inv_123")
	if err == nil {
		t.Fatal("expected error")
	}
	if _, ok := err.(*RateLimitError); !ok {
		t.Fatalf("expected RateLimitError, got %T", err)
	}
	if attempt < 2 {
		t.Fatalf("expected retries, got %d attempts", attempt)
	}
}
