package db

import (
	"log"
	"Productbid/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func ConnectDB(dsn string) *gorm.DB {
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	log.Println("Database connected successfully")

	// Auto-migrate all models — creates tables if they don't exist
	err = database.AutoMigrate(
		&models.Category{},
		&models.Product{},
		&models.Bid{},
		models.Payment{},
	)

	if err != nil {
		log.Fatal("Migration failed: ", err)
	}

	log.Println("Migration completed successfully!!")

	SeedCategories(database)

	return database
}

func SeedCategories(database *gorm.DB) {
	categories := []models.Category{
		{Name: "AI Agents & Infrastructure", Slug: "ai-agents-infrastructure"},
		{Name: "SEO & AI Visibility", Slug: "seo-ai-visibility"},
		{Name: "Marketing & Advertising", Slug: "marketing-advertising"},
		{Name: "Crypto, Web3 & Investing", Slug: "crypto-web3-investing"},
		{Name: "Developer Tools", Slug: "developer-tools"},
		{Name: "Business, Finance & Legal", Slug: "business-finance-legal"},
		{Name: "Security, Privacy & Compliance", Slug: "security-privacy-compliance"},
		{Name: "Health, Fitness & Wellness", Slug: "health-fitness-wellness"},
		{Name: "Social Media & Creator Tools", Slug: "social-media-creator-tools"},
		{Name: "Leaderboards & Attention Markets", Slug: "leaderboards-attention-markets"},
		{Name: "Hiring, Jobs & Careers", Slug: "hiring-jobs-careers"},
		{Name: "Education & Learning", Slug: "education-learning"},
		{Name: "Agencies, Studios & Services", Slug: "agencies-studios-services"},
		{Name: "Ecommerce & Retail", Slug: "ecommerce-retail"},
		{Name: "Domains & Web Assets", Slug: "domains-web-assets"},
		{Name: "Games & Entertainment", Slug: "games-entertainment"},
		{Name: "People & Profiles", Slug: "people-profiles"},
		{Name: "Productivity & Personal Tools", Slug: "productivity-personal-tools"},
		{Name: "Design & Creative", Slug: "design-creative"},
		{Name: "Writing & Content", Slug: "writing-content"},
		{Name: "Directories, Launch & Discovery", Slug: "directories-launch-discovery"},
		{Name: "AI Media Generation", Slug: "ai-media-generation"},
		{Name: "Audio, Voice & Podcasting", Slug: "audio-voice-podcasting"},
		{Name: "Sales & Lead Generation", Slug: "sales-lead-generation"},
		{Name: "Travel, Local & Lifestyle", Slug: "travel-local-lifestyle"},
		{Name: "Real Estate & Property", Slug: "real-estate-property"},
		{Name: "Media & News", Slug: "media-news"},
		{Name: "Other", Slug: "other"},
	}

	for _, category := range categories {
		var existing models.Category
		result := database.Where("slug = ?", category.Slug).First(&existing)

		if result.RowsAffected == 0 {
			database.Create(&category)
			log.Printf("Seeded category: %s", category.Name)
		}
	}
}