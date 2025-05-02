package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content  string `json:"content"`
	Username string `json:"username"`
	MusicId  string `json:"music_id"`
}
