package main

import (
	"flag"
	"log"
	apihttp "majoo-case1-rest-api/api/http"
	"majoo-case1-rest-api/config"
	"majoo-case1-rest-api/internal/comment"
	"majoo-case1-rest-api/internal/database"
	"majoo-case1-rest-api/internal/http/middleware"
	"majoo-case1-rest-api/internal/post"
	"majoo-case1-rest-api/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	envPath := flag.String("env-path", "", "Path to .env file (default: config/.env or .env)")
	flag.Parse()

	// Load env file if path is provided
	if *envPath != "" {
		if err := godotenv.Load(*envPath); err != nil {
			log.Printf("Warning: Failed to load env file from %s: %v", *envPath, err)
		} else {
			log.Printf("Loaded env from %s", *envPath)
		}
	}

	cfg := config.Load()

	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to init DB:", err)
	}
	defer db.Close()
    // Migrations are managed via golang-migrate and the Makefile targets

	r := gin.Default()

	// CORS minimal setup (adjust origin in production)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Wiring usecases
	userRepo := user.NewRepository(db)
	userUC := user.NewUsecase(userRepo, cfg)
	postRepo := post.NewRepository(db)
	postUC := post.NewUsecase(db, postRepo)
	commentRepo := comment.NewRepository(db)
	commentUC := comment.NewUsecase(db, commentRepo)

	api := r.Group("/api/v1")
	apihttp.RegisterAuthRoutes(api, userUC, cfg)

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg))
	apihttp.RegisterPostRoutes(protected, postUC)
	apihttp.RegisterCommentRoutes(protected, commentUC)

	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
