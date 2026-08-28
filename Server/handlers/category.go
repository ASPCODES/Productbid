package handlers

import(
	"Productbid/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


type CategoryHandler struct {
	DB *gorm.DB
}


func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{DB: db}
}

func (h *CategoryHandler) GetCategories(c *fiber.Ctx) error {
	var categories []models.Category

	result := h.DB.Find(&categories)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch categories",
		})
	}

	return c.JSON(fiber.Map{
		"categories": categories,
	})
}

