package adapter

import (
	"music/repository"

	"gorm.io/gorm"
)

type gormLibraryRepository struct {
	db *gorm.DB
}

func NewLibraryGormDB(db *gorm.DB) repository.UserLibraryRepossitory {
	return &gormLibraryRepository{db: db}
}

func (g *gormLibraryRepository) Save(library repository.UserLibrary) error {
	results := g.db.Create(&library)
	if results.Error != nil {
		return results.Error
	}

	return nil
}

func (g *gormLibraryRepository) Delete(song_id uint, user_id uint) error {
	resulte := g.db.Where("user_id = ? AND song_id = ?", user_id, song_id).Delete(&repository.UserLibrary{})
	if resulte != nil {
		return resulte.Error
	}
	return nil
}

func (g *gormLibraryRepository) ShowUserID(song_id uint, user_id uint) ([]repository.UserLibrary, error) {
	var library []repository.UserLibrary
	err := g.db.Where("user_id = ?", user_id).Find(&library).Error
	if err != nil {
		return nil, err
	}
	return library, nil
}
