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

func GetPlaylists(c echo.Context) error {
	username := c.Param("username")
	playlists := new([]models.Playlist)
	var user struct {
		Playlists string `json:"playlists"`
	}
	if err := database.DB.Model(&models.User{}).
		Select("playlists").
		Where("username = ?", username).
		First(&user).Error; err != nil {
		return c.JSON(400, echo.Map{"error": "User not found"})
	}

	playlistsArr := pq.StringArray{}
	if user.Playlists != "" {
		cleaned := strings.Trim(user.Playlists, "{}")
		if cleaned != "" {
			playlistsArr = strings.Split(cleaned, ",")
		}
	}

	var ids []uint
	for _, idStr := range playlistsArr {
		id, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			continue
		}
		ids = append(ids, uint(id))
	}

	if len(ids) > 0 {
		if err := database.DB.
			Where("id IN (?)", ids).
			Find(&playlists).Error; err != nil {
			return c.JSON(400, echo.Map{
				"error":   "Failed to fetch playlists",
				"details": err.Error(),
			})
		}

	}

	return c.JSON(200, echo.Map{
		"playlists": playlists,
	})
}

func CreatePlaylist(c echo.Context) error {
	username := c.Param("username")
	playlist := new(models.Playlist)
	if err := c.Bind(&playlist); err != nil {
		fmt.Println("alo")
		fmt.Println(err.Error())
		return c.JSON(400, echo.Map{"error": "invalid json"})
	}
	if err := database.DB.Create(&playlist).Error; err != nil {
		return c.JSON(400, echo.Map{"error": "failed to create playlist"})
	}
	id := fmt.Sprintf("%d", playlist.ID)
	if err := database.DB.Model(&models.User{}).Where("username = ?", username).Update("playlists", gorm.Expr("array_append(playlists, ?)", id)).Error; err != nil {
		return c.JSON(400, echo.Map{"error": "failed to add playlist to user"})
	}
	return c.JSON(201, echo.Map{"playlist": playlist})
}

func DeletePlaylist(c echo.Context) error {
	id := c.Param("id")
	username := c.Param("username")
	playlist := new(models.Playlist)
	if err := database.DB.Model(&models.Playlist{}).Where("id = ?", id).First(&playlist).Error; err != nil {
		return c.JSON(400, echo.Map{"error": "playlist not found"})
	}
	if err := database.DB.Model(&models.User{}).Where("username = ?", username).Update("playlists", gorm.Expr("array_remove(playlists, ?)", id)).Error; err != nil {
		return c.JSON(400, echo.Map{"error": "failed to remove playlist from user"})
	}
	if err := database.DB.Delete(&playlist, "id = ?", id).Error; err != nil {
		return c.JSON(500, echo.Map{"error": "failed to delete comment"})
	}
	return c.NoContent(200)
}

func GetPlaylist(c echo.Context) error {
	id := c.Param("id")
	playlist := new(models.Playlist)
	if err := database.DB.Model(&models.Playlist{}).Where("id = ?", id).First(&playlist).Error; err != nil {
		return c.JSON(400, echo.Map{"error": "playlist not found"})
	}
	return c.JSON(200, echo.Map{"playlist": playlist})
}

func GetAllPlaylistMusic(c echo.Context) error {
	id := c.Param("id")
	music := new([]models.MusicTrack)
	var playlist struct {
		Music string `json:"music"`
	}
	if err := database.DB.Model(&models.Playlist{}).
		Select("music").
		Where("id = ?", id).
		First(&playlist).Error; err != nil {
		return c.JSON(400, echo.Map{"error": "playlist not found"})
	}

	musicArr := pq.StringArray{}
	if playlist.Music != "" {
		cleaned := strings.Trim(playlist.Music, "{}")
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

	if len(ids) > 0 {
		if err := database.DB.
			Where("id IN (?)", ids).
			Find(&music).Error; err != nil {
			return c.JSON(400, echo.Map{
				"error":   "Failed to fetch music",
				"details": err.Error(),
			})
		}

	}

	return c.JSON(200, echo.Map{
		"music": music,
	})
}
