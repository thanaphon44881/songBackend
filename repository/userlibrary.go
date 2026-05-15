package repository

type UserLibrary struct {
	UserID uint `gorm:"primaryKey"`
	SongID uint `gorm:"primaryKey"`

	User User `gorm:"constraint:OnDelete:CASCADE;"`
	Song Song `gorm:"constraint:OnDelete:CASCADE;"`
}

type UserLibraryRepossitory interface {
	Save(library UserLibrary) error
	Delete(song_id uint, user_id uint) error
	ShowUserID(song_id uint, user_id uint) ([]UserLibrary, error)
}
