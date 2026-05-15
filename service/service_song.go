package service

import (
	"music/repository"
	"music/slugs"

	"strconv"
)

type songServiceImpl struct {
	repo repository.RepositorySong
}

func NewSongService(r repository.RepositorySong) ServiceSong {
	return &songServiceImpl{repo: r}
}

func (s *songServiceImpl) Creat(song Song) error {
	slug, err := s.generateUniqueSlug(song.Title)
	if err != nil {
		return err
	}
	song.Slug = slug
	songs := repository.Song{
		Title:    song.Title,
		File:     song.File,
		Cover:    song.Cover,
		Duration: song.Duration,
		Category: song.Category,
		Views:    0,
		ArtistID: song.ArtistID,
		Sub:      song.Sub,
		Slug:     song.Slug,
		Country:  song.Country,
	}

	err = s.repo.Save(songs)
	if err != nil {
		return err
	}
	return nil
}

func (s *songServiceImpl) Find() ([]Song, error) {
	song, err := s.repo.Show()
	if err != nil {
		return nil, nil
	}

	songsRepo := []Song{}

	for _, songs := range song {
		songsRepo = append(songsRepo, mapToSongResponse(songs, true))
	}
	return songsRepo, nil
}

func (s *songServiceImpl) FindNew() ([]Song, []Song, error) {
	songnew, songtop, err := s.repo.ShowNewTop()
	if err != nil {
		return nil, nil, err
	}
	songsRepoN := []Song{}
	songsRepoT := []Song{}

	for _, songsn := range songnew {
		songsRepoN = append(songsRepoN, mapToSongResponse(songsn, true))
	}
	for _, songtop := range songtop {
		songsRepoT = append(songsRepoT, mapToSongResponse(songtop, true))
	}
	return songsRepoN, songsRepoT, nil
}

func (s *songServiceImpl) First(id int) (*Song, error) {
	song, err := s.repo.ShowID(uint(id))
	if err != nil {
		return nil, err
	}
	songs := mapToSongResponse(*song, true)
	return &songs, err
}

func (s *songServiceImpl) FirstLibrary(id int, userID uint) (*Song, error) {
	song, err := s.repo.ShowID(uint(id))
	if err != nil {
		return nil, err
	}
	liked, err := s.repo.LibrarySong(userID, uint(id))
	if err != nil {
		return nil, err
	}
	songs := mapToSongResponse(*song, false)
	songs.IsLiked = liked
	return &songs, err
}

func (s *songServiceImpl) UpdateId(id int) error {
	err := s.repo.IncreaseView(uint(id))
	if err != nil {
		return err
	}
	return nil
}

func (s *songServiceImpl) PlayQueue(id int, userID uint) (*PlayResponse, error) {
	current, err := s.repo.ShowID(uint(id))
	if err != nil {
		return nil, err
	}

	nextSongs, err := s.repo.GetNextSongs(uint(id))
	if err != nil {
		return nil, err
	}

	liked, err := s.repo.LibrarySong(userID, uint(id))
	if err != nil {
		return nil, err
	}

	resp := &PlayResponse{
		Song: mapToSongResponse(*current, false),
		Next: []Song{},
	}

	resp.Song.IsLiked = liked

	for _, item := range nextSongs {
		resp.Next = append(resp.Next, mapToSongResponse(item, false))
	}

	return resp, nil
}

func (s *songServiceImpl) generateUniqueSlug(title string) (string, error) {
	base := slugs.GenerateSlug(title)
	slug := base
	i := 1

	for {
		exists, err := s.repo.CheckSlugExists(slug)
		if err != nil {
			return "", err
		}
		if !exists {
			break
		}
		slug = base + "-" + strconv.Itoa(i)
		i++
	}

	return slug, nil
}

func mapToSongResponse(song repository.Song, includeSub bool) Song {
	res := Song{
		ID:       song.ID,
		Title:    song.Title,
		File:     song.File,
		Cover:    song.Cover,
		Duration: song.Duration,
		Views:    song.Views,
		Up:       song.Up,
		Slug:     song.Slug,
		Sub:      song.Sub,
		Category: song.Category,
		Country:  song.Country,
		Artist: repository.Artist{
			ID:    song.Artist.ID,
			Name:  song.Artist.Name,
			Image: song.Artist.Image,
			Bio:   song.Artist.Bio,
		},
	}

	if includeSub {
		song.Sub = "-"
		song.Artist.Bio = "-"
		res.Sub = song.Sub
		res.Artist.Bio = song.Artist.Bio
	}

	return res
}
