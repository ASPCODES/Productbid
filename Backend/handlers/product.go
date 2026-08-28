package handlers

import (
	"Productbid/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProductHandler struct {
	DB *gorm.DB
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{DB: db}
}

type ResolveProductInput struct {
	HandleOrURL  string `json:"handle_or_url"`
	Name         string `json:"name"`
	Tagline      string `json:"tagline"`
	LogoURL      string `json:"logo_url"`
	CategoryID   string `json:"category_id"`
	ContactEmail string `json:"contact_email"`
}

func (h *ProductHandler) ResolveProduct(c *fiber.Ctx) error {
	var input ResolveProductInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if input.HandleOrURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "handle_or_url is required",
		})
	}

	var existingProduct models.Product
	result := h.DB.Where("handle_or_url = ?", input.HandleOrURL).First(&existingProduct)

	// If the product already exist, it returns that
	if result.RowsAffected > 0 {
		return c.JSON(fiber.Map{
			"product": existingProduct,
			"existed": true,
		})
	}

	// validation of making new product
	categoryID, err := strconv.ParseUint(input.CategoryID, 10, 64)

	if input.Name == "" || input.CategoryID == "" || err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name and category_id are required for new products",
		})
	}

	newProduct := models.Product{
		HandleOrURL:  input.HandleOrURL,
		Name:         input.Name,
		Tagline:      input.Tagline,
		LogoURL:      input.LogoURL,
		CategoryID:   uint(categoryID),
		ContactEmail: input.ContactEmail,
	}

	if err := h.DB.Create(&newProduct).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create product",
		})
	}

	return c.JSON(fiber.Map{
		"product": newProduct,
		"existed": false,
	})
}
