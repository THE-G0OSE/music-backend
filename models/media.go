package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	Username string `json:"username"`
	Filename string `json:"filename"`
	Path     string `json:"path" gorm:"unique"`
	Size     int64  `json:"size"`
}

type MusicTrack struct {
	gorm.Model
	Username   string         `json:"username"`
	Filename   string         `json:"filename"`
	Path       string         `json:"path" gorm:"unique"`
	CoverImage string         `json:"cover_image"`
	Size       int64          `json:"size"`
	Author     string         `json:"author"`
	Title      string         `json:"title"`
	Likes      int64          `json:"likes"`
	Genre      string         `json:"genre"`
	Comments   pq.StringArray `json:"comments" gorm:"type:text[]"`
}
