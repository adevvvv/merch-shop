package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"shop/internal/handlers"
	"shop/internal/middleware"
	"shop/internal/models"
	"shop/internal/repositories"
	"shop/internal/services"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db := initDB()
	defer db.Close()

	// Инициализация репозиториев
	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	inventoryRepo := repositories.NewInventoryRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	// Сервисы
	authService := services.NewAuthService(
		userRepo,
		sessionRepo,
		[]byte(os.Getenv("JWT_SECRET")),
		os.Getenv("USE_SESSION") == "true",
	)
	
	userService := services.NewUserService(userRepo, inventoryRepo, transactionRepo)
	transactionService := services.NewTransactionService(userRepo, transactionRepo)

	// Обработчики
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Роутер
	router := gin.Default()
	
	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(
		sessionRepo,
		[]byte(os.Getenv("JWT_SECRET")),
		os.Getenv("USE_SESSION") == "true",
	)

	// Маршруты
	api := router.Group("/api")
	{
		api.POST("/auth", authHandler.AuthHandler)
		api.GET("/info", authMiddleware.Handler(), userHandler.GetInfo)
		api.POST("/send", authMiddleware.Handler(), transactionHandler.SendCoin)
		api.POST("/buy/:item", authMiddleware.Handler(), userHandler.BuyItem)
	}

	// Сервер
	srv := &http.Server{
		Addr:    ":" + getEnv("PORT", "8080"),
		Handler: router,
	}

	go gracefulShutdown(srv)
	log.Fatal(srv.ListenAndServe())
}

func initDB() *sql.DB {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "user"),
		getEnv("DB_PASSWORD", "password"),
		getEnv("DB_NAME", "shop"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown error:", err)
	}
	log.Println("Server stopped")
}