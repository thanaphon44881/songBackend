package service

import (
	"errors"
	"music/repository"

	"golang.org/x/crypto/bcrypt"
)

type userServiceImpl struct {
	repo repository.UserRepossitory
}

func NewUserService(repo repository.UserRepossitory) ServiceUser {
	return &userServiceImpl{repo: repo}
}

func (r *userServiceImpl) Creat(user User) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	users := repository.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Biaoname: user.Biaoname,
		Image:    user.Image,
	}
	err = r.repo.Save(users)
	if err != nil {
		return err
	}
	return nil
}

func (r *userServiceImpl) FindByUsername(usert UserToken) (*User, error) {
	username := usert.Username
	user, err := r.repo.Logins(username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(usert.Password)); err != nil {
		return nil, errors.New("wrong password")
	}
	users := User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
	}
	return &users, err
}

func (r *userServiceImpl) First(id int) (*User, error) {
	user, err := r.repo.GetuserID(uint(id))
	if err != nil {
		return nil, err
	}
	users := User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Biaoname: user.Biaoname,
		Image:    user.Image,
	}
	return &users, nil
}
