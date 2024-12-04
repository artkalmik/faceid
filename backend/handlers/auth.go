package handlers

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"math"

	"faceid/models"

	"github.com/gin-gonic/gin"
)

// Функция для вычисления среднего значения пикселей
func calculateAveragePixels(img image.Image) []float64 {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Разделим изображение на 16 частей (4x4 сетка)
	cellWidth := width / 4
	cellHeight := height / 4
	averages := make([]float64, 16)

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			var sum float64
			var count float64

			// Вычисляем среднее значение для каждой ячейки
			for x := i * cellWidth; x < (i+1)*cellWidth; x++ {
				for y := j * cellHeight; y < (j+1)*cellHeight; y++ {
					r, g, b, _ := img.At(x, y).RGBA()
					// Преобразуем в оттенки серого
					gray := (float64(r) + float64(g) + float64(b)) / (3 * 65535)
					sum += gray
					count++
				}
			}
			averages[i*4+j] = sum / count
		}
	}
	return averages
}

// Функция для вычисления сходства между двумя наборами средних значений
func calculateSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}

	var sum float64
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}

	// Евклидово расстояние
	distance := math.Sqrt(sum)
	// Преобразуем в показатель сходства (0-1)
	similarity := 1.0 / (1.0 + distance)
	return similarity
}

func Register(c *gin.Context) {
	var data struct {
		Username string `json:"username"`
		FaceData string `json:"faceData"` // base64 строка
	}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": "Invalid data"})
		return
	}

	// Декодируем base64 в []byte
	faceData, err := base64.StdEncoding.DecodeString(data.FaceData)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid face data"})
		return
	}

	// Сохраняем пользователя
	models.Users[data.Username] = models.User{
		Username: data.Username,
		FaceData: faceData,
	}

	c.JSON(200, gin.H{"message": "User registered successfully"})
}

func Verify(c *gin.Context) {
	var data struct {
		Username string `json:"username"`
		FaceData string `json:"faceData"`
	}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": "Invalid data"})
		return
	}

	user, exists := models.Users[data.Username]
	if !exists {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Декодируем полученные данные
	newFaceData, err := base64.StdEncoding.DecodeString(data.FaceData)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid face data"})
		return
	}

	// Преобразуем []byte в image.Image
	storedImg, err := jpeg.Decode(bytes.NewReader(user.FaceData))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to process stored image"})
		return
	}

	newImg, err := jpeg.Decode(bytes.NewReader(newFaceData))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to process new image"})
		return
	}

	// Вычисляем средние значения для обоих изображений
	storedAverages := calculateAveragePixels(storedImg)
	newAverages := calculateAveragePixels(newImg)

	// Вычисляем сходство
	similarity := calculateSimilarity(storedAverages, newAverages)

	// Порог сходства (можно настроить)
	const similarityThreshold = 0.85

	if similarity >= similarityThreshold {
		c.JSON(200, gin.H{"verified": true, "similarity": similarity})
	} else {
		c.JSON(200, gin.H{"verified": false, "similarity": similarity})
	}
}

func GetUsers(c *gin.Context) {
	usernames := make([]string, 0)
	for username := range models.Users {
		usernames = append(usernames, username)
	}
	c.JSON(200, usernames)
}
