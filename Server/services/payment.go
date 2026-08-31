package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// MinimumBidAmount is the fixed floor for every bid ($3.00), enforced here
// AND on the Dodo Payments dashboard via Pay What You Want minimum price.
const MinimumBidAmount = 3.00

type PaymentService struct {
	DodoAPIKey        string
	DodoWebhookSecret string
	DodoMode          string
	BaseURL           string
	HTTPClient        *http.Client
}

func NewPaymentService(apiKey, webhookSecret, mode string) *PaymentService {
	baseURL := "https://test.dodopayments.com"
	if strings.ToLower(mode) == "live" {
		baseURL = "https://live.dodopayments.com"
	}

	return &PaymentService{
		DodoAPIKey:        apiKey,
		DodoWebhookSecret: webhookSecret,
		DodoMode:          strings.ToLower(mode),
		BaseURL:           baseURL,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// ---------- Checkout Session ----------

// CheckoutProductItem matches Dodo Payments' actual product_cart schema.
// Amount is only used because the product has "Pay What You Want" enabled
// on the Dodo dashboard — it must NOT be sent for fixed-price products.
type CheckoutProductItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Amount    int    `json:"amount,omitempty"` // cents, PWYW products only
}

type CheckoutCustomer struct {
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
}

type CreateCheckoutRequest struct {
	ProductCart []CheckoutProductItem `json:"product_cart"`
	Customer    CheckoutCustomer      `json:"customer,omitempty"`
	Metadata    map[string]string     `json:"metadata,omitempty"`
	ReturnURL   string                `json:"return_url,omitempty"`
}

// CreateCheckoutResponse matches the real Dodo Payments API response shape.
type CreateCheckoutResponse struct {
	SessionID      string `json:"session_id"`
	CheckoutURL    string `json:"checkout_url"`
	ClientSecret   string `json:"client_secret"`
	PaymentID      string `json:"payment_id"`
	PublishableKey string `json:"publishable_key"`
}

// CreateCheckoutSession calls Dodo Payments API to create a hosted checkout
// session for a bid. productID must be a PWYW-enabled one-time product with
// a $3.00 minimum configured on the Dodo dashboard.
func (s *PaymentService) CreateCheckoutSession(
	bidID uint,
	productID string,
	bidAmount float64,
	customerEmail string,
	customerName string,
	returnURL string,
) (checkoutURL string, sessionID string, err error) {

	if productID == "" {
		return "", "", errors.New("productID is required")
	}
	if bidAmount < MinimumBidAmount {
		return "", "", fmt.Errorf("bid amount must be at least $%.2f", MinimumBidAmount)
	}
	if s.DodoAPIKey == "" {
		return "", "", errors.New("dodo API key is not configured")
	}

	amountInCents := int(bidAmount*100 + 0.5) // round to nearest cent

	reqBody := CreateCheckoutRequest{
		ProductCart: []CheckoutProductItem{
			{
				ProductID: productID,
				Quantity:  1,
				Amount:    amountInCents,
			},
		},
		Customer: CheckoutCustomer{
			Email: customerEmail,
			Name:  customerName,
		},
		Metadata: map[string]string{
			"bid_id":     strconv.FormatUint(uint64(bidID), 10),
			"amount_usd": fmt.Sprintf("%.2f", bidAmount),
		},
		ReturnURL: returnURL,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal checkout request: %w", err)
	}
	log.Printf("[Dodo] Request Body: %s", string(jsonData))

	url := fmt.Sprintf("%s/checkouts", s.BaseURL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.DodoAPIKey))

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		log.Printf("[Dodo Payments] request failed: %v", err)
		return "", "", fmt.Errorf("failed to reach dodo payments: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return "", "", fmt.Errorf("failed to read dodo response: %w", readErr)
	}

	if resp.StatusCode >= 400 {
		log.Printf("[Dodo Payments] API error status %d: %s", resp.StatusCode, string(bodyBytes))
		return "", "", fmt.Errorf("dodo payments error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var checkoutResp CreateCheckoutResponse
	if err := json.Unmarshal(bodyBytes, &checkoutResp); err != nil {
		return "", "", fmt.Errorf("failed to parse dodo checkout response: %w", err)
	}

	if checkoutResp.CheckoutURL == "" {
		return "", "", errors.New("dodo payments did not return a checkout_url")
	}

	sessionID = checkoutResp.SessionID
	if sessionID == "" {
		sessionID = checkoutResp.PaymentID
	}

	return checkoutResp.CheckoutURL, sessionID, nil
}

// ---------- Webhook Verification ----------

// VerifyWebhookSignature verifies Standard Webhooks HMAC-SHA256 signature.
func (s *PaymentService) VerifyWebhookSignature(
	payload []byte,
	signatureHeader string,
	webhookID string,
	webhookTimestamp string,
) bool {
	secret := strings.TrimSpace(s.DodoWebhookSecret)
	if secret == "" {
		log.Println("[Dodo Webhook] DODO_WEBHOOK_SECRET is empty — rejecting webhook")
		return false // never silently accept unsigned webhooks, even in test mode
	}

	if signatureHeader == "" || webhookID == "" || webhookTimestamp == "" {
		log.Println("[Dodo Webhook] Missing required webhook headers")
		return false
	}

	// Validate timestamp tolerance (5 minutes) to prevent replay attacks.
	ts, err := strconv.ParseInt(webhookTimestamp, 10, 64)
	if err != nil {
		log.Println("[Dodo Webhook] Invalid timestamp header")
		return false
	}
	diff := time.Now().Unix() - ts
	if diff < -300 || diff > 300 {
		log.Printf("[Dodo Webhook] Timestamp out of tolerance: now=%d received=%d", time.Now().Unix(), ts)
		return false
	}

	// Strip 'whsec_' prefix if present, then base64-decode per Standard Webhooks spec.
	cleanSecret := strings.TrimPrefix(secret, "whsec_")
	secretBytes, err := base64.StdEncoding.DecodeString(cleanSecret)
	if err != nil {
		secretBytes = []byte(secret)
	}

	// Signed content format: "{webhook_id}.{webhook_timestamp}.{raw_body}"
	payloadToSign := fmt.Sprintf("%s.%s.%s", webhookID, webhookTimestamp, string(payload))

	mac := hmac.New(sha256.New, secretBytes)
	mac.Write([]byte(payloadToSign))
	expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))


	// signatureHeader can contain space-separated signatures, e.g. "v1,sigA v1,sigB"
	for _, sig := range strings.Split(signatureHeader, " ") {
		parts := strings.SplitN(sig, ",", 2)
		candidate := sig
		if len(parts) == 2 && parts[0] == "v1" {
			candidate = parts[1]
		}
		if subtle.ConstantTimeCompare([]byte(candidate), []byte(expectedSignature)) == 1 {
			return true
		}
	}

	log.Println("[Dodo Webhook] Signature verification failed")
	return false
}

// ---------- Webhook Event Parsing ----------

type DodoWebhookData struct {
	PaymentID   string            `json:"payment_id"`
	ID          string            `json:"id"`
	Status      string            `json:"status"`
	TotalAmount int               `json:"total_amount"`
	Currency    string            `json:"currency"`
	Metadata    map[string]string `json:"metadata"`
}

type DodoWebhookPayload struct {
	Type      string           `json:"type"`
	Timestamp string           `json:"timestamp"`
	Data      DodoWebhookData  `json:"data"`
}

// ParseWebhookEvent extracts bidID, paymentID, and a normalized status
// ("success" | "failed" | raw status) from a verified Dodo webhook body.
func (s *PaymentService) ParseWebhookEvent(payload []byte) (bidID uint, paymentID string, status string, err error) {
	var event DodoWebhookPayload
	if err := json.Unmarshal(payload, &event); err != nil {
		return 0, "", "", fmt.Errorf("failed to parse webhook payload: %w", err)
	}

	switch strings.ToLower(event.Type) {
	case "payment.succeeded", "payment.successful", "payment_succeeded":
		status = "success"
	case "payment.failed", "payment_failed":
		status = "failed"
	default:
		switch strings.ToLower(event.Data.Status) {
		case "succeeded", "success":
			status = "success"
		case "failed":
			status = "failed"
		default:
			status = event.Data.Status
		}
	}

	paymentID = event.Data.PaymentID
	if paymentID == "" {
		paymentID = event.Data.ID
	}

	if bidIDStr, ok := event.Data.Metadata["bid_id"]; ok && bidIDStr != "" {
		if parsed, parseErr := strconv.ParseUint(bidIDStr, 10, 64); parseErr == nil {
			bidID = uint(parsed)
		}
	}

	return bidID, paymentID, status, nil
}

