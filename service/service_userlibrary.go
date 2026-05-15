package service

import (
	"music/repository"
)

type userLibraryImpl struct {
	repo repository.UserLibraryRepossitory
}

func NewuserLibraryService(r repository.UserLibraryRepossitory) ServiceuserLibrary {
	return &userLibraryImpl{repo: r}
}

func (u *userLibraryImpl) Creat(userID uint, songID uint) (bool, error) {
	librarys := repository.UserLibrary{
		UserID: userID,
		SongID: songID,
	}
	err := u.repo.Save(librarys)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *userLibraryImpl) DeleteUser(song_id int, user_id int) error {
	err := u.repo.Delete(uint(song_id), uint(user_id))
	if err != nil {
		return err
	}
	return nil
}

func (u *userLibraryImpl) FingUser(song_id int, user_id int) ([]repository.UserLibrary, error) {
	library, err := u.repo.ShowUserID(uint(song_id), uint(user_id))

	if err != nil {
		return nil, err
	}
	return library, nil
}
