package handlers

import (
	"Productbid/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


type BidHandler struct {
	DB *gorm.DB
}


func NewBidHandler(db *gorm.DB) *BidHandler {
	return &BidHandler{DB: db}
}


type BidPreviewInput struct {
	CategoryID      uint		`json:"category_id"`
	Amount			float64		`json:"amount"`
}


type PlaceBidInput struct {
	ProductID		uint		`json:"product_id"`
	CategoryID		uint		`json:"category_id"`
	Amount			float64		`json:"amount"`
}


// PreviewRank tells the user what rank they would get if they bid this amount
func (h* BidHandler) PreviewRank(c *fiber.Ctx) error {
	var input BidPreviewInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if input.CategoryID == 0 || input.Amount <=0 {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "category_id and a valid amount are required",
		})
	}

	// Count how many active bids in this category have a higher amount
	var higherBidCount int64
	h.DB.Model(&models.Bid{}).Where("category_id = ? AND is_active = ? AND amount > ?", input.CategoryID, true, input.Amount).Count(&higherBidCount)


	// Rank = number of higher bids + 1
	predictedRank := higherBidCount + 1

	return c.JSON(fiber.Map{
		"predicted_rank": predictedRank,
	})
}


// PlaceBid creates a new bid for a product
func (h *BidHandler) PlaceBid(c *fiber.Ctx) error {
	var input PlaceBidInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if input.ProductID == 0 || input.CategoryID == 0 || input.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "product_id, category_id and a valid amount are required",
		})
	}

	// Confirm the product actually exists
	var product models.Product

	if err := h.DB.First(&product, input.ProductID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	newBid := models.Bid {
		ProductID: input.ProductID,
		CategoryID: input.CategoryID,
		Amount: input.Amount,
		IsActive: true,  // temporary: directly active until payment integration added
	}

	if err := h.DB.Create(&newBid).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to place bid",
		})
	}

	return c.JSON(fiber.Map{
		"bid": newBid,
	})
}