# Vylth Flow Go SDK

Official Go SDK for [Vylth Flow](https://flow.vylth.com) — self-custody crypto payment processing.

## Installation

```bash
go get github.com/VYLTH/flow-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    flow "github.com/VYLTH/flow-go"
)

func main() {
    client := flow.New("vf_live_...")

    // Create an invoice
    invoice, err := client.Invoices.Create(context.Background(), &flow.CreateInvoiceParams{
        Amount:   100,
        Currency: "USDT",
        Network:  "tron",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Invoice:", invoice.ID, invoice.PaymentURL)
}
```

## Resources

### Invoices

```go
// Create
inv, _ := client.Invoices.Create(ctx, &flow.CreateInvoiceParams{
    Amount: 100, Currency: "USDT", Network: "tron",
    Description: "Order #123",
    CallbackURL: "https://example.com/webhook",
})

// Get
inv, _ := client.Invoices.Get(ctx, "inv_123")

// List
list, _ := client.Invoices.List(ctx, &flow.ListParams{Page: 1, Limit: 20})

// Cancel
inv, _ := client.Invoices.Cancel(ctx, "inv_123")
```

### Payouts

```go
// Create
payout, _ := client.Payouts.Create(ctx, &flow.CreatePayoutParams{
    Amount: 50, Currency: "USDT", Network: "tron",
    Destination: "TXyz...",
})

// Batch create
payouts, _ := client.Payouts.CreateBatch(ctx, []*flow.CreatePayoutParams{
    {Amount: 50, Currency: "USDT", Network: "tron", Destination: "TXyz..."},
    {Amount: 25, Currency: "USDT", Network: "tron", Destination: "TAbc..."},
})

// Get & List
payout, _ := client.Payouts.Get(ctx, "po_123")
list, _ := client.Payouts.List(ctx, &flow.ListParams{Status: "completed"})
```

### Wallets

```go
// Generate
wallet, _ := client.Wallets.Generate(ctx, &flow.GenerateWalletParams{
    Network: "tron", Currency: "USDT", Label: "deposits",
})

// Get, List, Balance
wallet, _ := client.Wallets.Get(ctx, "w_123")
list, _ := client.Wallets.List(ctx, &flow.ListParams{Network: "tron"})
wallet, _ := client.Wallets.Balance(ctx, "w_123")
```

### Swaps

```go
// Get quote
quote, _ := client.Swaps.Quote(ctx, &flow.SwapParams{
    FromCurrency: "USDT", ToCurrency: "BTC", Amount: 1000,
})

// Execute swap
swap, _ := client.Swaps.Create(ctx, &flow.SwapParams{
    FromCurrency: "USDT", ToCurrency: "BTC", Amount: 1000,
})

// Get & List
swap, _ := client.Swaps.Get(ctx, "sw_123")
list, _ := client.Swaps.List(ctx, nil)
```

### Webhooks

```go
// With webhook secret
client := flow.New("vf_live_...", flow.WithWebhookSecret("whsec_..."))

// Verify in HTTP handler
func webhookHandler(w http.ResponseWriter, r *http.Request) {
    event, err := client.Webhooks.VerifyRequest(r)
    if err != nil {
        http.Error(w, "Invalid signature", 401)
        return
    }
    switch event.Type {
    case "invoice.paid":
        // Handle payment
    case "payout.completed":
        // Handle payout
    }
    w.WriteHeader(200)
}

// Manual verification
event, err := client.Webhooks.Verify(payload, signature)

// Quick check
if client.Webhooks.IsValid(payload, signature) {
    // Process
}
```

## Error Handling

```go
inv, err := client.Invoices.Get(ctx, "inv_123")
if err != nil {
    switch e := err.(type) {
    case *flow.AuthenticationError:
        // Invalid API key
    case *flow.NotFoundError:
        // Resource not found
    case *flow.ValidationError:
        // Invalid parameters — check e.Errors
    case *flow.RateLimitError:
        // Too many requests — check e.RetryAfter
    case *flow.ServerError:
        // Server error — retry later
    default:
        // Other error
    }
}
```

## Configuration

```go
client := flow.New("vf_live_...",
    flow.WithBaseURL("https://custom-api.example.com"),
    flow.WithTimeout(60 * time.Second),
    flow.WithWebhookSecret("whsec_..."),
)
```

## License

MIT
