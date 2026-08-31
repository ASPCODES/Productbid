package handlers

import (
	"Productbid/models"
	"Productbid/services"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BidHandler struct {
	DB             *gorm.DB
	BidService     *services.BidService
	PaymentService *services.PaymentService
	FrontendURL    string
	DodoProductID  string 
}

func NewBidHandler(
	db *gorm.DB,
	bidService *services.BidService,
	paymentService *services.PaymentService,
	frontendURL string,
	dodoProductID string, 
) *BidHandler {
	return &BidHandler{
		DB:             db,
		BidService:     bidService,
		PaymentService: paymentService,
		FrontendURL:    frontendURL,
		DodoProductID:  dodoProductID,
	}
}

type BidPreviewInput struct {
	CategoryID   uint    `json:"category_id"`
	CategorySlug string  `json:"category_slug"`
	Amount       float64 `json:"amount"`
}

type InitiateBidInput struct {
	ProductID     uint    `json:"product_id"`
	CategoryID    uint    `json:"category_id"`
	CategorySlug  string  `json:"category_slug"`
	Amount        float64 `json:"amount"`
	CustomerEmail string  `json:"customer_email"`
	ReturnURL     string  `json:"return_url"`

	// Direct resolve fields for flexibility
	HandleOrURL string `json:"handle_or_url"`
	Name        string `json:"name"`
	Tagline     string `json:"tagline"`
	LogoURL     string `json:"logo_url"`
}

// PreviewRank calculates what rank the product would achieve with the given bid
func (h *BidHandler) PreviewRank(c *fiber.Ctx) error {
	var input BidPreviewInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	categoryID := input.CategoryID
	if categoryID == 0 && input.CategorySlug != "" {
		var category models.Category
		if err := h.DB.Where("slug = ?", input.CategorySlug).First(&category).Error; err == nil {
			categoryID = category.ID
		}
	}

	if categoryID == 0 || input.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "category_id (or category_slug) and a valid amount are required",
		})
	}

	predictedRank, err := h.BidService.PreviewRank(categoryID, input.Amount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate rank preview",
		})
	}

	return c.JSON(fiber.Map{
		"predicted_rank": predictedRank,
		"category_id":    categoryID,
		"amount":         input.Amount,
	})
}

// InitiateBid validates product/category, creates a pending bid, generates a Dodo Payments checkout session, and returns the checkout URL
func (h *BidHandler) InitiateBid(c *fiber.Ctx) error {
	var input InitiateBidInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if input.Amount < 3.0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Minimum bid amount is $3.00",
		})
	}

	customerEmail := strings.TrimSpace(input.CustomerEmail)
	if customerEmail == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "customer_email is required",
		})
	}

	// Resolve Category if slug provided
	categoryID := input.CategoryID
	if categoryID == 0 && input.CategorySlug != "" {
		var category models.Category
		if err := h.DB.Where("slug = ?", input.CategorySlug).First(&category).Error; err == nil {
			categoryID = category.ID
		}
	}

	// Resolve or find Product
	productID := input.ProductID
	var product models.Product

	if productID != 0 {
		if err := h.DB.First(&product, productID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Sprintf("Product with ID %d not found", productID),
			})
		}
	} else if input.HandleOrURL != "" {
		result := h.DB.Where("handle_or_url = ?", input.HandleOrURL).First(&product)
		if result.RowsAffected > 0 {
			productID = product.ID
		} else {
			if input.Name == "" || categoryID == 0 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "name and category are required for new products",
				})
			}
			newProduct := models.Product{
				HandleOrURL:  input.HandleOrURL,
				Name:         input.Name,
				Tagline:      input.Tagline,
				LogoURL:      input.LogoURL,
				CategoryID:   categoryID,
				ContactEmail: customerEmail,
			}
			if err := h.DB.Create(&newProduct).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to create product record",
				})
			}
			product = newProduct
			productID = product.ID
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "product_id or handle_or_url is required",
		})
	}

	if categoryID == 0 {
		categoryID = product.CategoryID
	}

	// 1. Create Pending Bid (is_active: false)
	pendingBid, err := h.BidService.CreatePendingBid(productID, categoryID, input.Amount)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 2. Return URL destination
	returnURL := input.ReturnURL
	if returnURL == "" {
		returnURL = fmt.Sprintf("%s/?payment=success&bid_id=%d", h.FrontendURL, pendingBid.ID)
	}

	// 3. Create Checkout Session with Dodo Payments
	checkoutURL, sessionID, err := h.PaymentService.CreateCheckoutSession(
		pendingBid.ID,
		h.DodoProductID,   
		pendingBid.Amount,
		customerEmail,
		product.Name,   
		returnURL,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to initiate payment session: %v", err),
		})
	}

	// 4. Create Payment row in database
	payment := models.Payment{
		BidID:         pendingBid.ID,
		Amount:        pendingBid.Amount,
		Status:        "pending",
		DodoSessionID: sessionID,
		PayerEmail:    customerEmail,
	}
	h.DB.Create(&payment)

	return c.JSON(fiber.Map{
		"success":         true,
		"checkout_url":    checkoutURL,
		"bid_id":          pendingBid.ID,
		"dodo_session_id": sessionID,
		"product_id":      productID,
		"amount":          pendingBid.Amount,
	})
}

// DeleteBid permanently deletes a bid and its associated payment records (DELETE /api/bids/:id)
func (h *BidHandler) DeleteBid(c *fiber.Ctx) error {
	idParam := c.Params("id")
	bidID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || bidID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid bid ID",
		})
	}

	if err := h.BidService.DeleteBid(uint(bidID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Bid deleted successfully",
		"bid_id":  bidID,
	})
}
