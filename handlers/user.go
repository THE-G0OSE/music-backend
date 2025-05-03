package handlers

import (
	"fmt"
	"music-backend/database"
	"music-backend/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type MusicRequest struct {
	Username string `json:"username"`
}

func GetUser(c echo.Context) error {
	username := c.Param("username")
	var user models.User
	database.DB.First(&user, "username = ?", username)
	return c.JSON(http.StatusOK, user)
}

func Login(c echo.Context) error {
	req := new(AuthRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(400, echo.Map{"error": "Invalid request"})
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).Select("username", "password", "music", "playlists", "likes", "image", "genre").First(&user).Error; err != nil {
		return echo.ErrUnauthorized
	}

	if user.Password != req.Password {
		return echo.ErrUnauthorized
	}

	if user.Likes == nil {
		user.Likes = []string{}
	}

	if user.Playlists == nil {
		user.Playlists = []string{}
	}

	if user.Music == nil {
		user.Music = []string{}
	}

	return c.JSON(200, models.ToUserResponse(user))
}

func CreateUser(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	if user.Username == "" || user.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username and password are required"})
	}

	if user.Likes == nil {
		user.Likes = []string{}
	}
	if user.Music == nil {
		user.Music = []string{}
	}
	if user.Playlists == nil {
		user.Playlists = []string{}
	}
	user.Genre = "rock"
	if err := database.DB.Create(user).Error; err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	return c.JSON(http.StatusCreated, user)
}

func UpdateUserGenre(c echo.Context) error {
	username := c.Param("username")
	genre := c.FormValue("genre")
	if err := database.DB.Model(&models.User{}).Debug().Where("username = ?", username).Update("genre", genre).Error; err != nil {
		fmt.Println(err.Error())
		return c.JSON(404, "failed to update genre")
	}
	return c.NoContent(http.StatusOK)
}

func LikeMusic(c echo.Context) error {
	username := c.Param("username")
	musicId := c.Param("id")
	music := new(models.MusicTrack)
	if err := database.DB.Model(&models.User{}).Where("username = ?", username).Update("likes", gorm.Expr("array_append(likes, ?)", musicId)).Error; err != nil {
		return c.JSON(400, echo.Map{"error": "failed to append like"})
	}
	if err := database.DB.Model(&models.MusicTrack{}).Where("id = ?", musicId).First(&music).Error; err != nil {
		return c.JSON(400, echo.Map{"error": "music not found"})
	}
	if err := database.DB.Model(&models.MusicTrack{}).Where("id = ?", musicId).Update("likes", music.Likes+1).Error; err != nil {
		return c.JSON(500, echo.Map{"error": "failed to update likes"})
	}
	return c.NoContent(200)
}
func UnlikeMusic(c echo.Context) error {
	username := c.Param("username")
	musicId := c.Param("id")
	music := new(models.MusicTrack)
	if err := database.DB.Model(&models.User{}).Where("username = ?", username).Update("likes", gorm.Expr("array_remove(likes, ?)", musicId)).Error; err != nil {
		return c.JSON(400, echo.Map{"error": "failed to remove like"})
	}
	if err := database.DB.Model(&models.MusicTrack{}).Where("id = ?", musicId).First(&music).Error; err != nil {
		return c.JSON(400, echo.Map{"error": "music not found"})
	}
	if err := database.DB.Model(&models.MusicTrack{}).Where("id = ?", musicId).Update("likes", music.Likes-1).Error; err != nil {
		return c.JSON(500, echo.Map{"error": "failed to update likes"})
	}
	return c.NoContent(200)
}

func DeleteUser(c echo.Context) error {
	username := c.Param("username")
	database.DB.Delete(&models.User{}, "username = ?", username)
	return c.NoContent(http.StatusNoContent)
}

func GetAllUserMusic(c echo.Context) error {
	req := new(MusicRequest)
	if err := c.Bind(req); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid json data"})
	}
	var user struct {
		Music string `json:"music"`
	}
	if err := database.DB.Model(&models.User{}).
		Select("music").
		Where("username = ?", req.Username).
		First(&user).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
	}

	musicArr := pq.StringArray{}
	if user.Music != "" {
		cleaned := strings.Trim(user.Music, "{}")
		if cleaned != "" {
			musicArr = strings.Split(cleaned, ",")
		}
	}

	var ids []uint
	for _, idStr := range musicArr {
		id, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			continue
		}
		ids = append(ids, uint(id))
	}

	music := new([]models.MusicTrack)

	if len(ids) > 0 {
		if err := database.DB.
			Where("id IN (?)", ids).
			Find(&music).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error":   "Failed to fetch music",
				"details": err.Error(),
			})
		}

	}

	return c.JSON(200, echo.Map{
		"music": music,
	})
}

func GetLikedMusic(c echo.Context) error {
	username := c.Param("username")
	music := new([]models.MusicTrack)
	var user struct {
		Likes string `json:"likes"`
	}
	if err := database.DB.Select("likes").Where("username = ?", username).Model(&models.User{}).First(&user).Error; err != nil {
		return c.JSON(404, echo.Map{"error": "user not found"})
	}
	musicArr := pq.StringArray{}
	if user.Likes != "" {
		cleared := strings.Trim(user.Likes, "{}")
		if cleared != "" {
			musicArr = strings.Split(cleared, ",")
		}
	}
	var ids []uint
	for _, idStr := range musicArr {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, uint(id))
	}
	if len(ids) > 0 {
		if err := database.DB.Where("id IN (?)", ids).Find(&music).Error; err != nil {
			return c.JSON(400, echo.Map{"error": "music not found"})
		}
	}

	return c.JSON(200, echo.Map{"music": music})
}
