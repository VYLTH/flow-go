package flow

// Invoice represents a payment invoice.
type Invoice struct {
	ID            string                 `json:"id"`
	MerchantID    string                 `json:"merchant_id"`
	Amount        float64                `json:"amount"`
	Currency      string                 `json:"currency"`
	Network       string                 `json:"network"`
	Status        string                 `json:"status"`
	WalletAddress string                 `json:"wallet_address,omitempty"`
	Description   string                 `json:"description,omitempty"`
	CustomerEmail string                 `json:"customer_email,omitempty"`
	CallbackURL   string                 `json:"callback_url,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	PaymentURL    string                 `json:"payment_url,omitempty"`
	TxHash        string                 `json:"tx_hash,omitempty"`
	PaidAmount    float64                `json:"paid_amount,omitempty"`
	ExpiresAt     string                 `json:"expires_at,omitempty"`
	CreatedAt     string                 `json:"created_at,omitempty"`
	UpdatedAt     string                 `json:"updated_at,omitempty"`
}

// CreateInvoiceParams are parameters for creating an invoice.
type CreateInvoiceParams struct {
	Amount        float64                `json:"amount"`
	Currency      string                 `json:"currency"`
	Network       string                 `json:"network"`
	Description   string                 `json:"description,omitempty"`
	CustomerEmail string                 `json:"customer_email,omitempty"`
	CallbackURL   string                 `json:"callback_url,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	ExpiryMinutes int                    `json:"expiry_minutes,omitempty"`
}

// Payout represents a vendor payout.
type Payout struct {
	ID          string                 `json:"id"`
	Amount      float64                `json:"amount"`
	Currency    string                 `json:"currency"`
	Network     string                 `json:"network"`
	Destination string                 `json:"destination_address"`
	Status      string                 `json:"status"`
	TxHash      string                 `json:"tx_hash,omitempty"`
	Fee         float64                `json:"fee,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   string                 `json:"created_at,omitempty"`
	CompletedAt string                 `json:"completed_at,omitempty"`
}

// CreatePayoutParams are parameters for creating a payout.
type CreatePayoutParams struct {
	Amount      float64                `json:"amount"`
	Currency    string                 `json:"currency"`
	Network     string                 `json:"network"`
	Destination string                 `json:"destination_address"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CallbackURL string                 `json:"callback_url,omitempty"`
}

// Wallet represents a deposit wallet.
type Wallet struct {
	ID        string  `json:"id"`
	Address   string  `json:"address"`
	Network   string  `json:"network"`
	Currency  string  `json:"currency"`
	Balance   float64 `json:"balance"`
	Label     string  `json:"label,omitempty"`
	CreatedAt string  `json:"created_at,omitempty"`
}

// GenerateWalletParams are parameters for generating a wallet.
type GenerateWalletParams struct {
	Network  string `json:"network"`
	Currency string `json:"currency"`
	Label    string `json:"label,omitempty"`
}

// SwapQuote represents a swap price quote.
type SwapQuote struct {
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	FromAmount   float64 `json:"from_amount"`
	ToAmount     float64 `json:"to_amount"`
	Rate         float64 `json:"rate"`
	ExpiresAt    string  `json:"expires_at,omitempty"`
}

// Swap represents a swap.
type Swap struct {
	ID           string  `json:"id"`
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	FromAmount   float64 `json:"from_amount"`
	ToAmount     float64 `json:"to_amount"`
	Rate         float64 `json:"rate"`
	Status       string  `json:"status"`
	TxHash       string  `json:"tx_hash,omitempty"`
	CreatedAt    string  `json:"created_at,omitempty"`
}

// SwapParams are parameters for creating a swap.
type SwapParams struct {
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	Amount       float64 `json:"amount"`
}

// WebhookEvent represents a verified webhook event.
type WebhookEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt string                 `json:"created_at,omitempty"`
}

// ListParams are common parameters for list endpoints.
type ListParams struct {
	Page     int    `json:"page,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Status   string `json:"status,omitempty"`
	Currency string `json:"currency,omitempty"`
	Network  string `json:"network,omitempty"`
}

// PaginatedList is a paginated response.
type PaginatedList[T any] struct {
	Items   []T  `json:"items"`
	Total   int  `json:"total"`
	Page    int  `json:"page"`
	PerPage int  `json:"per_page"`
	HasMore bool `json:"has_more"`
}
