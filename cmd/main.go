package main

import (
	"log"
	"net/http"

	"time"

	"github.com/isurucuma/rate-limiting/internal/domain"
	"github.com/isurucuma/rate-limiting/internal/infrastructure"
	handler "github.com/isurucuma/rate-limiting/internal/server"
	"github.com/isurucuma/rate-limiting/internal/usecase"
)

func main() {
	// Redis configuration
	redisAddr := "localhost:6379"
	redisPassword := "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81"
	redisDB := 0
	redisClient := infrastructure.NewRedisClient(redisAddr, redisPassword, redisDB)

	// Rate limiter configuration
	timeWindow := 1 * time.Minute
	limit := 10
	rateLimiter := domain.NewRateLimiter(redisClient, "rate_limiter_key", timeWindow, limit)

	// Initialize usecase and handler
	rateLimiterUsecase := usecase.NewRateLimiterUsecase(rateLimiter)
	httpHandler := handler.NewHTTPHandler(rateLimiterUsecase)

	// Setup router
	router := infrastructure.NewRouter(httpHandler)

	// Start HTTP server
	log.Println("Starting server on :8084...")
	if err := http.ListenAndServe(":8084", router); err != nil {
		log.Fatal(err)
	}
}
