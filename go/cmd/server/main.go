package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/microservices-boilerplate/go/internal/config"
	"github.com/microservices-boilerplate/go/internal/database"
	"github.com/microservices-boilerplate/go/internal/handlers"
	"github.com/microservices-boilerplate/go/internal/services"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	pool, err := database.NewPool(ctx, cfg.DBURL)
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}
	defer pool.Close()

	if err := database.CreateSchema(ctx, pool.Pool); err != nil {
		panic("failed to create schema: " + err.Error())
	}

	// Handlers
	healthHandler := handlers.NewHealthHandler(cfg)
	itemsHandler := handlers.NewItemsHandler(services.NewItemService(pool.Pool))

	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.GET("/", root)
	e.GET("/health", healthHandler.Liveness)
	e.GET("/health/ready", healthHandler.Readiness)

	api := e.Group("/items")
	api.GET("", itemsHandler.List)
	api.GET("/:id", itemsHandler.Get)
	api.POST("", itemsHandler.Create)
	api.PATCH("/:id", itemsHandler.Update)
	api.DELETE("/:id", itemsHandler.Delete)

	// Graceful shutdown
	go func() {
		if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		e.Logger.Fatal(err)
	}
}

func root(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"service": "go-microservice",
		"docs":    "OpenAPI coming soon - see README for API docs",
		"health":  "/health",
	})
}
