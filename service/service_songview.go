package service

import (
	"music/repository"
	"time"
)

type songViewImpl struct {
	repo repository.RepositorySongView
}

func NewSongViewService(r repository.RepositorySongView) ServiceSongView {
	return &songViewImpl{repo: r}
}

func (s *songViewImpl) Creat(userID, songID uint) error {
	return s.repo.Save(repository.SongView{
		UserID:   userID,
		SongID:   songID,
		PlayedAt: time.Now(),
	})
}

func (s *songViewImpl) GetByUserAll(userID uint) ([]SongViewResponse, error) {
	data, err := s.repo.ShowAll(userID)
	if err != nil {
		return nil, err
	}

	var result []SongViewResponse

	for _, v := range data {
		item := SongViewResponse{
			ID:       v.ID,
			PlayedAt: v.PlayedAt,
			UserId:   v.UserID,
			Song: SongMini{
				ID:       v.Song.ID,
				Title:    v.Song.Title,
				Cover:    v.Song.Cover,
				Duration: v.Song.Duration,
				Views:    v.Song.Views,
				Artist: ArtistMini{
					ID:    v.Song.Artist.ID,
					Name:  v.Song.Artist.Name,
					Image: v.Song.Artist.Image,
				},
			},
		}

		result = append(result, item)
	}

	return result, nil
}

func (s *songViewImpl) GetByUser(userID uint) ([]SongViewResponse, error) {
	data, err := s.repo.ShowID(userID)
	if err != nil {
		return nil, err
	}

	var result []SongViewResponse

	for _, v := range data {
		item := SongViewResponse{
			ID:       v.ID,
			PlayedAt: v.PlayedAt,
			UserId:   v.UserID,
			Song: SongMini{
				ID:       v.Song.ID,
				Title:    v.Song.Title,
				Cover:    v.Song.Cover,
				Duration: v.Song.Duration,
				Views:    v.Song.Views,
				Artist: ArtistMini{
					ID:    v.Song.Artist.ID,
					Name:  v.Song.Artist.Name,
					Image: v.Song.Artist.Image,
				},
			},
		}

		result = append(result, item)
	}

	return result, nil
}

func (s *songViewImpl) Delete(id int) error {
	err := s.repo.Delete(uint(id))
	if err != nil {
		return err
	}
	return nil
}
