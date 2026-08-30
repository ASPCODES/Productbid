package config


import(
	"log"
	"os"

	"github.com/joho/godotenv"
)


type Config struct {
	DatabaseURL       string
	Port              string
	DodoAPIKey        string
	DodoWebhookSecret string
	DodoMode          string
	FrontendURL       string
}

func LoadConfig() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found, reading from system env")
	}

	dodoMode := os.Getenv("DODO_MODE")
	if dodoMode == "" {
		dodoMode = "test"
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "https://productbid.space"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	cfg := &Config{
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		Port:              port,
		DodoAPIKey:        os.Getenv("DODO_API_KEY"),
		DodoWebhookSecret: os.Getenv("DODO_WEBHOOK_SECRET"),
		DodoMode:          dodoMode,
		FrontendURL:       frontendURL,
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is not set in .env")
	}


	return cfg
}