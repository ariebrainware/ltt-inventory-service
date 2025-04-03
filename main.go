// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	_ "time/tzdata"

	"github.com/ariebrainware/ltt-inventory-service/config"
	"github.com/ariebrainware/ltt-inventory-service/endpoint"
	"github.com/ariebrainware/ltt-inventory-service/middleware"
	"github.com/ariebrainware/ltt-inventory-service/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Load the configuration
	cfg := config.LoadConfig()

	// Set the timezone to Asia/Jakarta
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatalf("Error loading timezone: %v", err)
	}
	time.Local = location
	gormConfig := &gorm.Config{}
	if cfg.AppEnv == "production" {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	} else {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}
	db, err := config.ConnectMySQL()
	if err != nil {
		log.Fatalf("Error connecting to MySQL: %v", err)
	}
	db.AutoMigrate(&model.User{}, &model.Session{}, &model.Role{}, &model.InventoryMaster{}, &model.InventoryDetails{})

	// Set Gin mode from config
	gin.SetMode(cfg.GinMode)

	// Create a Gin router with default middleware
	r := gin.Default()

	// Use custom CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Basic HTTP handler for root path
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Welcome to %s!", cfg.AppName),
		})
	})
	// Group routes that require a valid login token
	auth := r.Group("/")
	auth.Use(middleware.ValidateLoginToken())
	{
		auth.DELETE("/logout", endpoint.Logout)

		auth.GET("/inventory", endpoint.ListInventory)
		auth.POST("/inventory", endpoint.CreateInventory)
		// auth.GET("/inventory/:id", endpoint.GetInventory)
		// auth.PATCH("/inventory/:id", endpoint.UpdateInventory)
		// auth.DELETE("/inventory/:id", endpoint.DeleteInventory)
		// auth.PUT("/inventory/:id", endpoint.TherapistApproval)
	}

	r.POST("/login", endpoint.Login)
	r.POST("/signup", endpoint.Signup)
	r.GET("/token/validate", endpoint.ValidateToken)

	// Start server on specified port
	address := fmt.Sprintf(":%d", cfg.AppPort)
	if err := r.Run(address); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
