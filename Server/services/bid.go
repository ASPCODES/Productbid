package services

import (
	"Productbid/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type BidService struct {
	DB *gorm.DB
}

func NewBidService(db *gorm.DB) *BidService {
	return &BidService{DB: db}
}

// CreatePendingBid validates input and persists a bid with is_active = false
func (s *BidService) CreatePendingBid(productID, categoryID uint, amount float64) (*models.Bid, error) {
	if productID == 0 {
		return nil, errors.New("product_id is required")
	}
	if categoryID == 0 {
		return nil, errors.New("category_id is required")
	}
	if amount < 3.0 {
		return nil, errors.New("bid amount must be at least $3.00")
	}

	// Verify product exists
	var product models.Product
	if err := s.DB.First(&product, productID).Error; err != nil {
		return nil, fmt.Errorf("product not found with id %d: %w", productID, err)
	}

	// Verify category exists
	var category models.Category
	if err := s.DB.First(&category, categoryID).Error; err != nil {
		return nil, fmt.Errorf("category not found with id %d: %w", categoryID, err)
	}

	newBid := models.Bid{
		ProductID:  productID,
		CategoryID: categoryID,
		Amount:     amount,
		IsActive:   false, // Only activated upon payment confirmation
	}

	if err := s.DB.Create(&newBid).Error; err != nil {
		return nil, fmt.Errorf("failed to create bid: %w", err)
	}

	return &newBid, nil
}

// ActivateBid marks the bid as active after webhook confirms payment
func (s *BidService) ActivateBid(bidID uint) error {
	result := s.DB.Model(&models.Bid{}).Where("id = ?", bidID).Update("is_active", true)

	if result.Error != nil {
		return fmt.Errorf("failed to activate bid %d: %w", bidID, result.Error)
	}
	
	if result.RowsAffected == 0 {
		return fmt.Errorf("no bid found with id %d", bidID)
	}
	return nil
}

// PreviewRank calculates the predicted rank for a given category and amount
func (s *BidService) PreviewRank(categoryID uint, amount float64) (int64, error) {
	if categoryID == 0 || amount <= 0 {
		return 1, errors.New("category_id and valid amount are required")
	}

	var higherBidCount int64
	err := s.DB.Model(&models.Bid{}).
		Where("category_id = ? AND is_active = ? AND amount > ?", categoryID, true, amount).
		Count(&higherBidCount).Error

	if err != nil {
		return 1, err
	}

	return higherBidCount + 1, nil
}

// DeleteBid permanently deletes a bid and its associated payment records from the database
func (s *BidService) DeleteBid(bidID uint) error {
	if bidID == 0 {
		return errors.New("invalid bid id")
	}

	// Permanently delete associated payment records
	if err := s.DB.Unscoped().Where("bid_id = ?", bidID).Delete(&models.Payment{}).Error; err != nil {
		return fmt.Errorf("failed to delete payment records for bid %d: %w", bidID, err)
	}

	// Permanently delete the bid row
	result := s.DB.Unscoped().Where("id = ?", bidID).Delete(&models.Bid{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete bid %d: %w", bidID, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("bid with id %d not found", bidID)
	}

	return nil
}
