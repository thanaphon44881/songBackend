package repository

import "time"

type SongView struct {
	ID       uint `gorm:"primaryKey"`
	UserID   uint `gorm:"index"`
	SongID   uint `gorm:"index"`
	PlayedAt time.Time

	User User
	Song Song
}

type RepositorySongView interface {
	Save(view SongView) error
	ShowID(userID uint) ([]SongView, error)
	Delete(id uint) error
	ShowAll(userID uint) ([]SongView, error)
}
