package config


import(
	"log"
	"os"

	"github.com/joho/godotenv"
)


type Config struct {
	DatabaseURL string
	Port 		string
}


func LoadConfig() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found, reading from system env")
	}

	cfg := &Config {
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is not set in .env")
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	return cfg
}