package main

import (
	"gym-admin/internal/config"
	"gym-admin/internal/router"
	"gym-admin/pkg/cache"
	"gym-admin/pkg/database"
	"gym-admin/pkg/jwt"
	"gym-admin/pkg/logger"
	"log"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger.InitLogger(cfg.Log)

	// Initialize JWT
	jwt.InitJWT(cfg.JWT.Secret)

	// Initialize database
	if err := database.InitDB(cfg.Database); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Redis cache
	if err := cache.InitRedis(cfg.Redis); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	// Setup router
	r := router.SetupRouter()

	// Start server
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
