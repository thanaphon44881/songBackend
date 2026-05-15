package service

import (
	"music/repository"
)

type artistServiceImpl struct {
	repo repository.RepositoryArtist
}

func NewArtistService(r repository.RepositoryArtist) ServiceArtist {
	return &artistServiceImpl{repo: r}
}

func (a *artistServiceImpl) Create(artist Artist) error {
	artists := repository.Artist{
		Name:  artist.Name,
		Image: artist.Image,
		Bio:   artist.Bio,
	}
	err := a.repo.Save(artists)
	if err != nil {
		return err
	}
	return nil
}

func (a *artistServiceImpl) First(id int) (*ArtistResponse, error) {
	artists, err := a.repo.ShowID(uint(id))
	if err != nil {
		return nil, err
	}
	artist := ArtistResponse{
		ID:    artists.ID,
		Name:  artists.Name,
		Image: artists.Image,
		Bio:   artists.Bio,
	}

	for _, s := range artists.Songs {
		artist.Songs = append(artist.Songs, SongResponse{
			ID:       s.ID,
			Title:    s.Title,
			Cover:    s.Cover,
			Duration: s.Duration,
			Views:    s.Views,
		})
	}

	return &artist, nil
}

func (a *artistServiceImpl) Find() ([]Artist, error) {
	artist, err := a.repo.Show()
	artists := []Artist{}
	for _, a := range artist {
		artists = append(artists, Artist{
			ID:    a.ID,
			Name:  a.Name,
			Image: a.Image,
		})
	}
	if err != nil {
		return nil, err
	}
	return artists, err
}
