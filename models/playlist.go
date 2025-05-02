package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Playlist struct {
	gorm.Model
	Username string         `json:"username"`
	Music    pq.StringArray `json:"music" gorm:"type:text[]"`
	Title    string         `json:"title"`
	Author   string         `json:"author"`
}
