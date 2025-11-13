package bootstrap

import (
	"akeneo-event-config/internal/akeneo"
	"akeneo-event-config/internal/config"
	"akeneo-event-config/internal/handlers"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	gin.SetMode(cfg.Server.GinMode)

	client := akeneo.NewClient(
		cfg.Akeneo.BaseURL,
		cfg.Akeneo.EventPlatformURL,
		cfg.Akeneo.ClientID,
		cfg.Akeneo.ClientSecret,
		cfg.Akeneo.Username,
		cfg.Akeneo.Password,
	)

	handler := handlers.NewHandler(client)

	r := gin.Default()

	// Configure CORS
	log.Printf("CORS allowed origins: %s", cfg.CORS.AllowedOrigins)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.CORS.AllowedOrigins},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}))

	// Request logging middleware
	r.Use(func(c *gin.Context) {
		log.Printf("[%s] %s - Origin: %s", c.Request.Method, c.Request.URL.Path, c.Request.Header.Get("Origin"))
		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		api.GET("/subscriber", handler.GetSubscriber)
		api.POST("/subscriber", handler.CreateSubscriber)
		api.PATCH("/subscriber/:id", handler.UpdateSubscriber)
		api.DELETE("/subscriber/:id", handler.DeleteSubscriber)

		api.GET("/subscriptions", handler.GetSubscriptions)
		api.POST("/subscriptions", handler.CreateSubscription)
		api.PATCH("/subscriptions/:id", handler.UpdateSubscription)
		api.DELETE("/subscriptions/:id", handler.DeleteSubscription)

		api.GET("/event-types", handler.GetEventTypes)
	}

	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Printf("PIM URL: %s", cfg.Akeneo.BaseURL)
	log.Printf("Event Platform URL: %s", cfg.Akeneo.EventPlatformURL)
	return r.Run(":" + cfg.Server.Port)
}
