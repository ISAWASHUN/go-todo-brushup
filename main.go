package main

import (
	"fmt"
	"go-todo-brushup/models"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// .envファイルから環境変数を読み込む
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 環境変数から接続情報を取得
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := "localhost" // または環境変数から取得
	dbPort := "5432"      // または環境変数から取得

	// DSNを構築
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo", dbHost, dbUser, dbPassword, dbName, dbPort)

	// GORMでデータベースに接続
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// データベースにテーブルを作成
	db.AutoMigrate(&models.Task{})

	// Ginエンジンのインスタンスを作成
	r := gin.Default()

	// タスクを取得するエンドポイント
	r.GET("/tasks", func(c *gin.Context) {
		var tasks []models.Task
		db.Find(&tasks)
		c.JSON(http.StatusOK, tasks)
	})

	// 新しいタスクを作成するエンドポイント
	r.POST("/tasks", func(c *gin.Context) {
		var task models.Task
		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Create(&task)
		c.JSON(http.StatusOK, task)
	})

	// タスクを更新するエンドポイント
	r.PUT("/tasks/:id", func(c *gin.Context) {
		var task models.Task
		id := c.Param("id")

		if err := db.First(&task, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}

		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		db.Save(&task)
		c.JSON(http.StatusOK, task)
	})

	// タスクを削除するエンドポイント
	r.DELETE("/tasks/:id", func(c *gin.Context) {
		var task models.Task
		id := c.Param("id")

		if err := db.First(&task, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}

		db.Delete(&task)
		c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
	})

	// 8080ポートでサーバーを起動
	r.Run(":8080")
}
