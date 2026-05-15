package repository

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Email    string
	Password string
	Biaoname string
	Image    string
}

type Playlist struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint
	Name        string
	Description string
}

type PlaylistSong struct {
	PlaylistID uint `gorm:"primaryKey"`
	SongID     uint `gorm:"primaryKey"`

	Playlist Playlist `gorm:"constraint:OnDelete:CASCADE;"`
	Song     Song     `gorm:"constraint:OnDelete:CASCADE;"`
}

type UserRepossitory interface {
	Save(user User) error
	Logins(username string) (*User, error)
	GetuserID(id uint) (*User, error)
}
