package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Dsn         string
	JwtSecret   string
	FrontendURL string

	// Stripe Config
	StripeSecretKey     string
	StripeWebhookSecret string
	StripeSuccessURL    string
	StripeCancelURL     string
}

func LoadEnv() *Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error loading .env file")
	}

	return &Config{
		Port:        os.Getenv("PORT"),
		Dsn:         os.Getenv("DSN"),
		JwtSecret:   os.Getenv("JWT_SECRET"),
		FrontendURL: os.Getenv("FRONTEND_URL"),

		// Stripe
		StripeSecretKey:     os.Getenv("STRIPE_SECRET_KEY"),
		StripeWebhookSecret: os.Getenv("STRIPE_WEBHOOK_SECRET"),
		StripeSuccessURL:    os.Getenv("STRIPE_SUCCESS_URL"),
		StripeCancelURL:     os.Getenv("STRIPE_CANCEL_URL"),
	}
}
