package repository

import (
	"time"

	"gorm.io/gorm"
)

type Song struct {
	gorm.Model
	Title    string
	File     string
	Cover    string
	Duration int
	Category string
	Slug     string `gorm:"uniqueIndex"`
	Sub      string
	Views    int
	ArtistID uint
	Artist   Artist `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Up       time.Time
	Country  string
}

type RepositorySong interface {
	Save(song Song) error
	Show() ([]Song, error)
	ShowNewTop() ([]Song, []Song, error)
	ShowID(id uint) (*Song, error)
	IncreaseView(id uint) error
	CheckSlugExists(slug string) (bool, error)
	LibrarySong(userId uint, id uint) (bool, error)
	GetNextSongs(currentID uint) ([]Song, error)
}
