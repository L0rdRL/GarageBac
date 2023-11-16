package controllers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/initializers"
	"github.com/models"
)

func AddDocument(c *gin.Context) {
	user, _ := c.Get("user")
	var body struct {
		Name   string
		Type   string
		Status string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Получение файла из запроса
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get file from request",
		})
		return
	}

	// Сохраняем файл в S3
	S3URL, err := uploadFileToS3(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload file to S3",
		})
		return
	}

	// Создаем новый документ
	document := models.Document{
		UserID:    user.(models.User).ID,
		Name:      body.Name,
		S3URL:     S3URL,
		Type:      body.Type,
		Status:    body.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Сохраняем документ в DB
	result := initializers.DB.Create(&document)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to add document",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})

}
func GetDocuments(c *gin.Context) {
	var documents []models.Document
	initializers.DB.Find(&documents)

	c.JSON(http.StatusOK, documents)
}

func uploadFileToS3(file *multipart.FileHeader) (string, error) {

	// Создание Клиента S3
	svc := initializers.SVC

	// Открытие файлов для чтения
	srcFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	// Формирование уникального имени файла в S3
	id := uuid.New()

	// Загрузка файла в S3
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(initializers.BUCKET_NAME),
		Key:    aws.String(id.String()), // Преобразуйте UUID в строку
		Body:   srcFile,
	})
	if err != nil {
		return "", err
	}
	// Формирование URL файла
	// Здесь вы можете создать URL, основываясь на имени бакета и ключа файла
	fileURL := fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", initializers.BUCKET_NAME, "kz-Ast", id.String())

	return fileURL, nil
}
