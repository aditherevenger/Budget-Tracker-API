package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/aditherevenger/Budget-Tracker-API/controllers"
	"github.com/aditherevenger/Budget-Tracker-API/database"
	"github.com/aditherevenger/Budget-Tracker-API/middleware"
	"github.com/aditherevenger/Budget-Tracker-API/repository"
	"github.com/aditherevenger/Budget-Tracker-API/services"
)

func main() {

	// Connect to the database
	database.Connect()
	database.Migrate()

	// Initialize repositories
	userRepo := repository.NewUserRepository()
	categoryRepo := repository.NewCategoryRepository()
	transactionRepo := repository.NewTransactionRepository()

	// Initialize services
	authService := services.NewAuthService(userRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	transactionService := services.NewTransactionService(transactionRepo, categoryRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	categoryController := controllers.NewCategoryController(categoryService)
	transactionController := controllers.NewTransactionController(transactionService)

	// Set up routes
	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	// Public routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	// Protected routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// User Profile
		api.GET("/profile", authController.GetProfile)

		//Categories
		categories := api.Group("/categories")
		{
			categories.GET("", categoryController.GetCategories)
			categories.POST("", categoryController.CreateCategories)
			categories.PUT("/:id", categoryController.UpdateCategory)
			categories.DELETE("/:id", categoryController.DeleteCategory)
		}

		//Transactions
		transactions := api.Group("/transactions")
		{
			transactions.GET("/", transactionController.GetTransactions)
			transactions.POST("/", transactionController.CreateTransaction)
			transactions.GET("/:id", transactionController.GetTransaction)
			transactions.PUT("/:id", transactionController.UpdateTransaction)
			transactions.DELETE("/:id", transactionController.DeleteTransaction)
			transactions.GET("/summary", transactionController.GetSummary)
		}
	}

	// Health Check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(router.Run(":" + port))
}
