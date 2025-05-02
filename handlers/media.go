package handlers

import (
	"fmt"
	"io"
	"music-backend/database"
	"music-backend/models"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

const (
	UploadDir     = "uploads"
	MaxImageSize  = 5 << 20
	MaxMusicSize  = 50 << 20
	AllowedImages = "image/jpeg, image/png"
)

func UpdateUserImage(c echo.Context) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.JSON(400, echo.Map{"error": "No image provided"})
	}
	username := c.FormValue("username")
	if !strings.Contains(AllowedImages, file.Header.Get("Content-Type")) {
		return c.JSON(400, echo.Map{"error": "Invalid image type"})
	}
	if file.Size > MaxImageSize {
		return c.JSON(400, echo.Map{"error": "Image too large"})
	}

	ext := filepath.Ext(file.Filename)
	newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(UploadDir, "images", newFilename)

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to open uploaded file"})
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to create file on server",
		})
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to save file content",
		})
	}

	image := models.Image{
		Username: username,
		Filename: file.Filename,
		Path:     filePath,
		Size:     file.Size,
	}

	if err := database.DB.Create(&image).Error; err != nil {
		fmt.Println(err)
		os.Remove(filePath)
		return err
	}

	if err := database.DB.Model(&models.User{}).Where("username = ?", username).Update("image", filePath).Error; err != nil {
		fmt.Println(err)
		os.Remove(filePath)
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "File uploaded successfully",
		"path":    filePath,
	})

}

func UploadMusicTrack(c echo.Context) error {
	musicFile, err := c.FormFile("music")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "music file not supported"})
	}
	imageFile, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "image file not supported"})
	}
	title := c.FormValue("title")
	author := c.FormValue("author")
	username := c.FormValue("username")
	genre := c.FormValue("genre")

	if musicFile.Size > MaxMusicSize {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "music file size is too large"})
	}
	if imageFile.Size > MaxImageSize {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "image file size is too large"})
	}

	musicExt := filepath.Ext(musicFile.Filename)
	newMusicFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), musicExt)
	musicFilePath := filepath.Join(UploadDir, "music", newMusicFileName)

	imageExt := filepath.Ext(imageFile.Filename)
	newImageFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), imageExt)
	imageFilePath := filepath.Join(UploadDir, "images", newImageFileName)

	musicSrc, err := musicFile.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to open music file"})
	}
	defer musicSrc.Close()

	imageSrc, err := imageFile.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to open image file"})
	}
	defer imageSrc.Close()

	musicDst, err := os.Create(musicFilePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create music file on server"})
	}
	defer musicDst.Close()

	imageDst, err := os.Create(imageFilePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create image file on server"})
	}
	defer imageDst.Close()

	if _, err = io.Copy(musicDst, musicSrc); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create music file on server"})
	}

	if _, err = io.Copy(imageDst, imageSrc); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create image file on server"})
	}
	comments := pq.StringArray{}

	music := models.MusicTrack{
		Username:   username,
		Filename:   musicFile.Filename,
		Path:       musicFilePath,
		Size:       musicFile.Size,
		CoverImage: imageFilePath,
		Title:      title,
		Author:     author,
		Likes:      0,
		Genre:      genre,
		Comments:   comments,
	}

	if err := database.DB.Create(&music).Error; err != nil {
		fmt.Println(err)
		os.Remove(imageFilePath)
		os.Remove(musicFilePath)
		return err
	}

	musicIDStr := fmt.Sprintf("%d", music.ID)

	if err := database.DB.Model(&models.User{}).Where("username = ?", username).Update("music", gorm.Expr("array_append(music, ?)", musicIDStr)).Error; err != nil {
		fmt.Println(err)
		os.Remove(musicFilePath)
		os.Remove(imageFilePath)
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":          "File uploaded successfully",
		"cover image path": imageFilePath,
		"music path":       musicFilePath,
		"music":            music.ID,
	})

}

func GetAllMusic(c echo.Context) error {
	music := new([]models.MusicTrack)
	if err := database.DB.Model(&models.MusicTrack{}).Find(&music).Error; err != nil {
		return c.JSON(404, echo.Map{"error": "failed to searce music"})
	}
	return c.JSON(200, echo.Map{
		"music": music,
	})
}

func GetLatestMusic(c echo.Context) error {
	music := new([]models.MusicTrack)
	if err := database.DB.Order("created_at DESC").Limit(10).Find(&music).Error; err != nil {
		fmt.Println(err.Error())
		return c.JSON(404, echo.Map{"error": "can not found latest tracks"})
	}
	return c.JSON(200, echo.Map{
		"music": music,
	})
}

func GetRecomendedMusic(c echo.Context) error {
	username := c.Param("username")
	user := new(models.User)
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return c.JSON(400, echo.Map{"error": "user not found"})
	}
	music := new([]models.MusicTrack)
	if err := database.DB.Where("genre = ?", user.Genre).Find(&music).Limit(10).Error; err != nil {
		return c.JSON(400, echo.Map{"error": "music not found"})
	}
	return c.JSON(200, echo.Map{"music": music})
}

func GetPopularMusic(c echo.Context) error {
	music := new([]models.MusicTrack)
	if err := database.DB.Model(&models.MusicTrack{}).Order("likes DESC").Limit(10).Find(&music).Error; err != nil {
		return c.JSON(404, echo.Map{"error": "music not found"})
	}
	return c.JSON(200, echo.Map{"music": music})
}

func GetMusic(c echo.Context) error {
	id := c.Param("id")
	music := new(models.MusicTrack)
	if err := database.DB.Where("id = ?", id).First(&music).Error; err != nil {
		return c.JSON(404, echo.Map{"error": "music not found"})
	}
	return c.JSON(200, echo.Map{"music": music})
}

func DeleteMusic(c echo.Context) error {
	id := c.Param("id")
	music := new(models.MusicTrack)
	if err := database.DB.Where("id = ?", id).First(&music).Error; err != nil {
		return c.JSON(404, echo.Map{"error": "music not found"})
	}
	if err := database.DB.Delete(&models.MusicTrack{}, "id = ?", id).Error; err != nil {
		return c.JSON(404, echo.Map{"error": "failed to delete music"})
	}
	database.DB.Model(&models.User{}).Where("username = ?", music.Username).Update("music", gorm.Expr("array_remove(music, ?)", id))
	os.Remove(music.Path)
	os.Remove(music.CoverImage)
	return c.NoContent(200)
}
