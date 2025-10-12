package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mathiuskitchens/whoop-xp/internal/database"
	"github.com/mathiuskitchens/whoop-xp/internal/spiritual"
	"github.com/mathiuskitchens/whoop-xp/internal/whoop"
	"net/http"
)

func main() {
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

	r.Run(":8080")

}
