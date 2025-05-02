package handlers

import (
	"fmt"
	"music-backend/database"
	"music-backend/models"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func CreateComment(c echo.Context) error {
	comment := new(models.Comment)
	if err := c.Bind(comment); err != nil {
		fmt.Println(err.Error())
		return c.JSON(400, echo.Map{"error": "json is invalid"})
	}
	if err := database.DB.Create(&comment).Error; err != nil {
		return c.JSON(500, echo.Map{"error": "failed to create comment"})
	}
	id := fmt.Sprintf("%d", comment.ID)
	if err := database.DB.Model(&models.User{}).Where("username = ?", comment.Username).Update("comments", gorm.Expr("array_append(comments, ?)", id)).Error; err != nil {
		return c.JSON(400, echo.Map{"error": "failed to add comment to user"})
	}
	if err := database.DB.Model(&models.MusicTrack{}).Where("id = ?", comment.MusicId).Update("comments", gorm.Expr("array_append(comments, ?)", id)).Error; err != nil {
		return c.JSON(400, echo.Map{"error": "failed to add comment to music"})
	}
	return c.JSON(201, echo.Map{"comment": comment})
}

func DeleteComment(c echo.Context) error {
	id := c.Param("id")
	comment := new(models.Comment)
	if err := database.DB.Model(&models.Comment{}).Where("id = ?", id).First(&comment).Error; err != nil {
		return c.JSON(400, echo.Map{"error": "comment not found"})
	}
	if err := database.DB.Model(&models.User{}).Where("username = ?", comment.Username).Update("comments", gorm.Expr("array_remove(comments, ?)", comment.ID)).Error; err != nil {
		return c.JSON(400, echo.Map{"error": "failed to remove comment from user"})
	}
	if err := database.DB.Model(&models.MusicTrack{}).Where("id = ?", comment.MusicId).Update("comments", gorm.Expr("array_remove(comments, ?)", comment.ID)).Error; err != nil {
		return c.JSON(400, echo.Map{"error": "failed to remove comment from music"})
	}
	if err := database.DB.Delete(&comment, "id = ?", comment.ID).Error; err != nil {
		return c.JSON(400, echo.Map{"error": "failed to delete comment"})
	}
	return c.NoContent(200)
}

func GetAllMusicComments(c echo.Context) error {
	id := c.Param("id")
	comments := new([]models.Comment)
	var music struct {
		Comments string `json:"comments"`
	}
	if err := database.DB.Model(&models.MusicTrack{}).
		Select("comments").
		Where("id = ?", id).
		First(&music).Error; err != nil {
		return c.JSON(400, echo.Map{"error": "Music not found"})
	}

	commentsArr := pq.StringArray{}
	if music.Comments != "" {
		cleaned := strings.Trim(music.Comments, "{}")
		if cleaned != "" {
			commentsArr = strings.Split(cleaned, ",")
		}
	}

	var ids []uint
	for _, idStr := range commentsArr {
		id, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			continue
		}
		ids = append(ids, uint(id))
	}

	if len(ids) > 0 {
		if err := database.DB.
			Where("id IN (?)", ids).
			Find(&comments).Error; err != nil {
			return c.JSON(400, echo.Map{
				"error":   "Failed to fetch comments",
				"details": err.Error(),
			})
		}

	}

	return c.JSON(200, echo.Map{
		"comments": comments,
	})
}

func GetAllUserComments(c echo.Context) error {
	username := c.Param("username")
	comments := new([]models.Comment)
	var user struct {
		Comments string `json:"comments"`
	}
	if err := database.DB.Model(&models.User{}).
		Select("comments").
		Where("username = ?", username).
		First(&user).Error; err != nil {
		return c.JSON(400, echo.Map{"error": "User not found"})
	}

	commentsArr := pq.StringArray{}
	if user.Comments != "" {
		cleaned := strings.Trim(user.Comments, "{}")
		if cleaned != "" {
			commentsArr = strings.Split(cleaned, ",")
		}
	}

	var ids []uint
	for _, idStr := range commentsArr {
		id, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			continue
		}
		ids = append(ids, uint(id))
	}

	if len(ids) > 0 {
		if err := database.DB.
			Where("id IN (?)", ids).
			Find(&comments).Error; err != nil {
			return c.JSON(400, echo.Map{
				"error":   "Failed to fetch comments",
				"details": err.Error(),
			})
		}

	}

	return c.JSON(200, echo.Map{
		"comments": comments,
	})
}
