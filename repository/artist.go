package repository

type Artist struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Image string
	Bio   string

	Songs []Song `gorm:"foreignKey:ArtistID"`
}

type RepositoryArtist interface {
	Save(artist Artist) error
	ShowID(id uint) (*Artist, error)
	Show() ([]Artist, error)
}
