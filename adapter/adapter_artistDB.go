package adapter

import (
	"music/repository"

	"gorm.io/gorm"
)

type gormArtistRepository struct {
	db *gorm.DB
}

func NewArtistGormDB(db *gorm.DB) *gormArtistRepository {
	return &gormArtistRepository{db: db}
}

func (g *gormArtistRepository) Save(artist repository.Artist) error {
	results := g.db.Create(&artist)
	if results.Error != nil {
		return results.Error
	}
	return nil
}

func (g *gormArtistRepository) ShowID(id uint) (*repository.Artist, error) {
	var artist repository.Artist
	results := g.db.Preload("Songs").First(&artist, id)
	if results.Error != nil {
		return nil, results.Error
	}

	return &artist, nil
}

func (g *gormArtistRepository) Show() ([]repository.Artist, error) {
	var artist []repository.Artist

	results := g.db.Find(&artist)
	if results.Error != nil {
		return nil, results.Error
	}
	return artist, nil
}
