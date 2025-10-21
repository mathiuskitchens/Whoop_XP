package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mathiuskitchens/whoop-xp/internal/database"
	"github.com/mathiuskitchens/whoop-xp/internal/models"
	"github.com/mathiuskitchens/whoop-xp/internal/spiritual"
	"github.com/mathiuskitchens/whoop-xp/internal/whoop"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var jwtSecret = []byte("your-super-secret-key")

func main() {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	clientID := os.Getenv("WHOOP_CLIENT_ID")
	log.Println("Client ID:", clientID)

	db := database.Connect()
	defer db.Close()

	r := gin.Default()

	// Get /metrics
	r.GET("/metrics", func(c *gin.Context) {
		metrics, err := whoop.GetAllMetrics(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, metrics)
	})

	// POST /metrics
	r.POST("/metrics", func(c *gin.Context) {
		var m whoop.Metric
		if err := c.ShouldBindJSON(&m); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := whoop.InsertMetric(db, m); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "metric added successfully"})
	})

	r.GET("/health", func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(500, gin.H{"status": "db down"})
			return
		}
		c.JSON(200, gin.H{"status": "pong"})
	})

	// Disciplines

	// GET disciplines
	r.GET("/disciplines", func(c *gin.Context) {
		data, err := spiritual.GetAllDisciplines(db)
		if err != nil {
			c.JSON(500, gin.H{"error on get all disciplines": err.Error()})
			return
		}
		c.JSON(200, data)
	})

	// POST disciplines
	r.POST("/disciplines", func(c *gin.Context) {
		var d spiritual.Discipline
		if err := c.ShouldBindJSON(&d); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := spiritual.InsertDiscipline(db, d); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"message": "discipline added successfully"})

	})

	r.POST("/register", func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		if req.Username == "" || req.Password == "" {
			c.JSON(400, gin.H{"error": "username and password required"})
			return
		}

		err := models.CreateUser(db, req.Username, req.Password)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}

		c.JSON(200, gin.H{"message": "User created successfully"})
	})

	r.POST("/login", func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		userID, err := models.AuthenticateUser(db, req.Username, req.Password)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid username or password"})
			return
		}

		// Create a JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(time.Hour * 72).Unix(),
		})

		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			c.JSON(500, gin.H{"error": "unable to create token"})
			return
		}

		c.JSON(200, gin.H{"token": tokenString})
	})

	whoopGroup := r.Group("/whoop")
	{
		whoopGroup.GET("/auth", whoop.StartAuthHandler(db))
	}

	// Start server...
	r.Run(":8080")

}
