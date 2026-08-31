package handlers

import (
	"Productbid/models"
	"Productbid/services"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type WebhookHandler struct {
	DB             *gorm.DB
	PaymentService *services.PaymentService
	BidService     *services.BidService
}

func NewWebhookHandler(
	db *gorm.DB,
	paymentService *services.PaymentService,
	bidService *services.BidService,
) *WebhookHandler {
	return &WebhookHandler{
		DB:             db,
		PaymentService: paymentService,
		BidService:     bidService,
	}
}

// HandleDodoWebhook verifies Dodo Payments webhook signatures, updates payment records, and activates bids upon successful payment
func (h *WebhookHandler) HandleDodoWebhook(c *fiber.Ctx) error {
	rawBody := c.Body()

	// Extract Standard Webhooks headers
	signature := c.Get("webhook-signature")
	if signature == "" {
		signature = c.Get("Webhook-Signature")
	}

	webhookID := c.Get("webhook-id")
	if webhookID == "" {
		webhookID = c.Get("Webhook-Id")
	}

	timestamp := c.Get("webhook-timestamp")
	if timestamp == "" {
		timestamp = c.Get("Webhook-Timestamp")
	}

	// 1. Verify webhook signature
	if !h.PaymentService.VerifyWebhookSignature(rawBody, signature, webhookID, timestamp) {
		log.Printf("[Dodo Webhook] Unauthorized webhook attempt. Webhook-Id: %s", webhookID)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid webhook signature",
		})
	}

	// 2. Parse the webhook payload
	bidID, dodoPaymentID, status, err := h.PaymentService.ParseWebhookEvent(rawBody)
	if err != nil {
		log.Printf("[Dodo Webhook] Error parsing webhook payload: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse webhook event payload",
		})
	}

	log.Printf("[Dodo Webhook] Received event for BidID: %d, PaymentID: %s, Status: %s", bidID, dodoPaymentID, status)

	// 3. Find and update Payment record
	var payment models.Payment
	var paymentFound bool

	if dodoPaymentID != "" {
		if err := h.DB.Where("dodo_session_id = ?", dodoPaymentID).First(&payment).Error; err == nil {
			paymentFound = true
		}
	}

	if !paymentFound && bidID != 0 {
		if err := h.DB.Where("bid_id = ?", bidID).Order("created_at DESC").First(&payment).Error; err == nil {
			paymentFound = true
		}
	}

	if paymentFound {
		payment.Status = status
		if dodoPaymentID != "" {
			payment.DodoPaymentID = dodoPaymentID
		}
		h.DB.Save(&payment)
		if bidID == 0 {
			bidID = payment.BidID
		}
	}

	// 4. If payment succeeded, activate the bid
	if status == "success" && bidID != 0 {
		if err := h.BidService.ActivateBid(bidID); err != nil {
			log.Printf("[Dodo Webhook] Error activating bid %d: %v", bidID, err)
		} else {
			log.Printf("[Dodo Webhook] Successfully activated Bid %d on leaderboard", bidID)
		}
	}

	// Return 200 OK promptly so Dodo does not retry
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "processed",
		"bid_id":  bidID,
		"payment": status,
	})
}
