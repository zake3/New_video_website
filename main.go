package main

import (
	"log"
	"net/http"
	// "os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"video-website/models"
	// "fmt"
	// "os/exec"
	// "strings"

)

func setupDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("videos.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.Video{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	return db
}

func setupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Static file serving
	r.Static("/static", "./static")
	r.Static("/uploads", "./static/uploads")

	// HTML templates
	r.LoadHTMLGlob("templates/*")

	// Routes
	r.GET("/", func(c *gin.Context) {
		videos := []models.Video{}
		db.Find(&videos)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"videos": videos,
		})
	})

	r.GET("/upload", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", nil)
	})

	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("video")
		if err != nil {
			c.String(http.StatusBadRequest, "Failed to upload video")
			return
		}

		// Create uploads directory if not exists
		// uploadsDir := "./static/uploads"
		// os.MkdirAll(uploadsDir, os.ModePerm)

		// // Generate unique filename
		// filename := filepath.Join(uploadsDir, file.Filename)
		// if err := c.SaveUploadedFile(file, filename); err != nil {
		// 	c.String(http.StatusInternalServerError, "Failed to save video")
		// 	return
		// }

		// // Save video metadata to database
		// video := models.Video{
		// 	Title:       c.PostForm("title"),
		// 	Description: c.PostForm("description"),
		// 	FilePath:    filename,
		// 	UploadDate:  time.Now(),
		// }
		// db.Create(&video)
		filename := filepath.Join("uploads", file.Filename)
fullPath := filepath.Join("./static", filename)

// Save video metadata to database
video := models.Video{
    Title:       c.PostForm("title"),
    Description: c.PostForm("description"),
    FilePath:    filename,  // Save relative path
    UploadDate:  time.Now(),
}
db.Create(&video)

// Save the actual file
if err := c.SaveUploadedFile(file, fullPath); err != nil {
    c.String(http.StatusInternalServerError, "Failed to save video")
    return
}

		c.Redirect(http.StatusFound, "/")
	})

	r.GET("/video/:id", func(c *gin.Context) {
		var video models.Video
		if err := db.First(&video, c.Param("id")).Error; err != nil {
			c.String(http.StatusNotFound, "Video not found")
			return
		}

		c.HTML(http.StatusOK, "view.html", gin.H{
			"video": video,
		})
	})

	return r
}


func main() {
	db := setupDatabase()
	r := setupRouter(db)

	// setupVideoUploadRoutes(r, db)

	log.Println("Starting server on :8080")
	r.Run(":8080")
}