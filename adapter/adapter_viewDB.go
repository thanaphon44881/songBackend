package adapter

import (
	"music/repository"

	"errors"

	"gorm.io/gorm"
)

type gormSongViweRepository struct {
	db *gorm.DB
}

func NewSongViewGormDB(db *gorm.DB) repository.RepositorySongView {
	return &gormSongViweRepository{db: db}
}

func (r *gormSongViweRepository) Save(view repository.SongView) error {
	return r.db.Create(&view).Error
}

func (r *gormSongViweRepository) ShowAll(userID uint) ([]repository.SongView, error) {
	var views []repository.SongView

	err := r.db.
		Preload("Song").
		Preload("Song.Artist").
		Where("user_id = ?", userID).
		Order("played_at DESC").
		Limit(20).
		Find(&views).Error

	return views, err
}

func (r *gormSongViweRepository) Delete(id uint) error {
	var views repository.SongView
	if err := r.db.First(&views, id); err.Error != nil {
		return errors.New("ไม่มี ip ในฐานข้อมูล")
	}
	return r.db.Delete(&views, id).Error
}

func (r *gormSongViweRepository) ShowID(userID uint) ([]repository.SongView, error) {
	var views []repository.SongView

	err := r.db.
		Preload("Song").
		Preload("Song.Artist").
		Where("user_id = ?", userID).
		Select("DISTINCT ON (song_id) *").
		Order("song_id, played_at DESC").
		Limit(10).
		Find(&views).Error

	return views, err
}
