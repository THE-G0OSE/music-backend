package main

import (
	"music-backend/database"
	"music-backend/handlers"
	"music-backend/models"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	//db connect
	database.Connect()
	database.DB.Debug().AutoMigrate(&models.User{}, &models.Image{}, &models.MusicTrack{}, &models.Comment{}, &models.Playlist{})

	e := echo.New()

	//cors off
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	//static
	e.Static("/uploads", "uploads")

	// users
	e.GET("/users/:username", handlers.GetUser)
	e.GET("/users/getLikes/:username", handlers.GetLikedMusic)

	e.POST("/users/getMusic", handlers.GetAllUserMusic)
	e.POST("/users", handlers.CreateUser)
	e.POST("/users/login", handlers.Login)

	e.PUT("/users/:username", handlers.UpdateUserGenre)

	e.DELETE("/users/:username", handlers.DeleteUser)

	//media
	e.GET("/media/get/:id", handlers.GetMusic)
	e.GET("/media/allMusic", handlers.GetAllMusic)
	e.GET("/media/getNew", handlers.GetLatestMusic)
	e.GET("/media/getRec/:username", handlers.GetRecomendedMusic)
	e.GET("/media/getPop", handlers.GetPopularMusic)

	e.POST("/media/uploadimage", handlers.UpdateUserImage)
	e.POST("/media/uploadtrack", handlers.UploadMusicTrack)

	e.PUT("/media/like/:id/:username", handlers.LikeMusic)
	e.PUT("/media/unlike/:id/:username", handlers.UnlikeMusic)

	e.DELETE("/media/delete/:id", handlers.DeleteMusic)

	//comments
	e.POST("/comments", handlers.CreateComment)

	e.GET("/comments/music/:id", handlers.GetAllMusicComments)
	e.GET("/comments/user/:username", handlers.GetAllUserComments)

	e.DELETE("/comments/:id", handlers.DeleteComment)

	//playlists
	e.GET("/playlists/:username", handlers.GetPlaylists)
	e.GET("/playlists/music/:id", handlers.GetAllPlaylistMusic)
	e.GET("/playlists/getOne/:id", handlers.GetPlaylist)

	e.POST("/playlists/:username", handlers.CreatePlaylist)

	e.DELETE("/playlists/:id/:username", handlers.DeletePlaylist)

	//start
	e.Logger.Fatal(e.Start(":3200"))
}
