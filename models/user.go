package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string         `json:"username" gorm:"unique;not null" validate:"required,min=6"`
	Password  string         `json:"password" gorm:"not null" validate:"required,min=8"`
	Image     string         `json:"image"`
	Likes     pq.StringArray `json:"likes" gorm:"type:text[]"`
	Music     pq.StringArray `json:"music" gorm:"type:text[]"`
	Playlists pq.StringArray `json:"playlists" gorm:"type:text[]"`
	Genre     string         `json:"genre"`
	Comments  pq.StringArray `json:"comments" gorm:"type:text[]"`
}

type UserResponse struct {
	Username  string         `json:"username"`
	Image     string         `json:"image"`
	Likes     pq.StringArray `json:"likes"`
	Music     pq.StringArray `json:"music"`
	Playlists pq.StringArray `json:"playlists"`
	Genre     string         `json:"genre"`
	Comments  pq.StringArray `json:"comments"`
}

func ToUserResponse(user User) UserResponse {
	return UserResponse{
		Username:  user.Username,
		Image:     user.Image,
		Likes:     user.Likes,
		Music:     user.Music,
		Playlists: user.Playlists,
		Genre:     user.Genre,
		Comments:  user.Comments,
	}
}
