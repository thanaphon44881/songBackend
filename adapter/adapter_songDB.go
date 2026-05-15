package adapter

import (
	"music/repository"

	"gorm.io/gorm"
)

type gormSongRepository struct {
	db *gorm.DB
}

func NewSongGormDB(db *gorm.DB) repository.RepositorySong {
	return &gormSongRepository{db: db}
}

func (g *gormSongRepository) Save(song repository.Song) error {
	results := g.db.Create(&song)
	if results.Error != nil {
		return results.Error
	}
	return nil
}

func (g *gormSongRepository) Show() ([]repository.Song, error) {
	var songs []repository.Song

	resultes := g.db.Preload("Artist").Find(&songs)
	if resultes.Error != nil {
		return nil, resultes.Error
	}

	return songs, nil
}

func (g *gormSongRepository) ShowNewTop() ([]repository.Song, []repository.Song, error) {
	var songsnew []repository.Song
	var songstop []repository.Song
	results := g.db.Preload("Artist").Order("created_at DESC").Limit(10).Find(&songsnew)
	if results.Error != nil {
		return nil, nil, results.Error
	}
	results = g.db.Preload("Artist").Order("views DESC").Limit(10).Find(&songstop)
	if results.Error != nil {
		return nil, nil, results.Error
	}
	return songsnew, songstop, nil
}

func (g *gormSongRepository) ShowID(id uint) (*repository.Song, error) {
	var song repository.Song
	results := g.db.Preload("Artist").First(&song, id)
	if results.Error != nil {
		return nil, results.Error
	}
	return &song, nil
}

func (r *gormSongRepository) IncreaseView(id uint) error {
	return r.db.Model(&repository.Song{}).
		Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + ?", 1)).Error
}

func (r *gormSongRepository) CheckSlugExists(slug string) (bool, error) {
	var count int64
	err := r.db.Model(&repository.Song{}).
		Where("slug = ?", slug).
		Count(&count).Error

	return count > 0, err
}

func (r *gormSongRepository) LibrarySong(userId uint, id uint) (bool, error) {
	var count int64
	r.db.Model(&repository.UserLibrary{}).
		Where("user_id = ? AND song_id = ?", userId, id).
		Count(&count)

	isLiked := count > 0
	return isLiked, nil
}

func (r *gormSongRepository) GetNextSongs(currentID uint) ([]repository.Song, error) {
	var songs []repository.Song

	err := r.db.Preload("Artist").Where("id != ?", currentID).Order("id desc").Limit(10).Find(&songs).Error
	if err != nil {
		return nil, err
	}

	return songs, err
}
