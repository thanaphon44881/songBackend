package adapter

import (
	"music/repository"

	"errors"

	"gorm.io/gorm"
)

type gormOrderRepository struct {
	db *gorm.DB
}

func NewGormUserDB(db *gorm.DB) repository.UserRepossitory {
	return &gormOrderRepository{db: db}
}

func (d *gormOrderRepository) Save(user repository.User) error {
	var existing repository.User
	if err := d.db.Where("username = ?", user.Username).First(&existing).Error; err == nil {
		return errors.New("username already exists")
	}
	result := d.db.Create(&user)
	if result != nil {
		return result.Error
	}
	result = d.db.Create(&repository.Playlist{
		UserID:      user.ID,
		Name:        user.Username + "Music",
		Description: "Music",
	})
	if result != nil {
		return result.Error
	}
	return nil
}

func (d *gormOrderRepository) Logins(username string) (*repository.User, error) {
	var user repository.User
	err := d.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *gormOrderRepository) GetuserID(id uint) (*repository.User, error) {
	var user repository.User
	err := d.db.First(&user, id)
	if err.Error != nil {
		return nil, err.Error
	}
	return &user, nil
}
