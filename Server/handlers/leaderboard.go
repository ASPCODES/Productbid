package handlers

import (
	"Productbid/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type LeaderboardHandler struct {
	DB *gorm.DB
}


func NewLeaderboardHandler(db *gorm.DB) *LeaderboardHandler {
	return &LeaderboardHandler{DB: db}
}


type LeaderboardEntry struct {
	Rank 		uint		`json:"rank"`
	BidID		uint		`json:"bid_id"`
	ProductID	uint		`json:"product_id"`
	Name 		string		`json:"name"`
	Tagline		string		`json:"tagline"`
	Logo_URL	string		`json:"logo_url"`
	HandleOrURL string		`json:"handle_or_url"`
	CategoryID	uint		`json:"category_id"`
	Amount		float64		`json:"amount"`
}


func (h* LeaderboardHandler) GetLeaderboard(c *fiber.Ctx) error {
	categoryParam := c.Query("category")

	var bids []models.Bid

	query := h.DB.Preload("Product").Where("is_active = ?", true).Order("amount DESC")


	// if category filter is given, then give me only that category bids
	if categoryParam != "" {
		categoryID, err := strconv.ParseInt(categoryParam, 10, 64)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid category id",
			})
		}

		query = query.Where("category_id = ?", categoryID)
	}

	if err := query.Find(&bids).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch leaderboard",
		})
	}

	
	// Convert each bid into clean leaderboard entry according to rank
	entries := make([]LeaderboardEntry, 0, len(bids))

	for i, bid := range bids {
		entries = append(entries, LeaderboardEntry{
			Rank: 			uint(i + 1),
			BidID: 			bid.ID,
			ProductID: 		bid.ProductID,
			Name: 			bid.Product.Name,
			Tagline: 		bid.Product.Tagline,
			Logo_URL: 		bid.Product.LogoURL,
			HandleOrURL: 	bid.Product.HandleOrURL,
			CategoryID: 	bid.CategoryID,
			Amount: 		bid.Amount,
		})
	}

	return  c.JSON(fiber.Map{
		"leaderboard": entries,
	})
}