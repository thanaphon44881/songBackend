package service

import (
	"music/repository"
	"time"
)

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Biaoname string `json:"biaoname"`
	Image    string `json:"image"`
}

type UserToken struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `json:"username"`
	// Email    string `json:"email"`
	Password string `json:"password"`
}

type Song struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Title    string `json:"title"`
	File     string `json:"file"`
	Cover    string `json:"cover"`
	Duration int    `json:"songtime"`
	Category string `json:"category"`
	Views    int    `json:"views"`
	ArtistID uint   `json:"artistid"`
	Up       time.Time
	Country  string            `json:"country"`
	Sub      string            `json:"sub"`
	Slug     string            `json:"slug"`
	Artist   repository.Artist `json:"artist"`
	IsLiked  bool              `json:"is_liked"`
}

type ArtistResponse struct {
	ID    uint           `json:"id"`
	Name  string         `json:"name"`
	Image string         `json:"image"`
	Bio   string         `json:"bio"`
	Songs []SongResponse `json:"songs" gorm:"foreignKey:ArtistID"`
}

type Artist struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Bio   string `json:"bio"`
}

type SongResponse struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Cover    string `json:"cover"`
	Duration int    `json:"songtime"`
	Views    int    `json:"views"`
	ArtistID uint
}

type SongViewResponse struct {
	ID       uint      `json:"id"`
	PlayedAt time.Time `json:"played_at"`
	Song     SongMini  `json:"song"`
	UserId   uint      `json:"userid"`
}

type SongMini struct {
	ID       uint       `json:"id"`
	Title    string     `json:"title"`
	Cover    string     `json:"cover"`
	Duration int        `json:"duration"`
	Views    int        `json:"views"`
	Artist   ArtistMini `json:"artist"`
}

type ArtistMini struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Bio   string `json:"bio"`
}

type PlayResponse struct {
	Song Song   `json:"song"`
	Next []Song `json:"next"`
}

type Library struct {
	Users User `json:"user"`
}

type ServiceUser interface {
	Creat(user User) error
	FindByUsername(user UserToken) (*User, error)
	First(id int) (*User, error)
}

type ServiceArtist interface {
	Create(artist Artist) error
	First(id int) (*ArtistResponse, error)
	Find() ([]Artist, error)
}

type ServiceSong interface {
	Creat(song Song) error
	Find() ([]Song, error)
	FindNew() ([]Song, []Song, error)
	First(id int) (*Song, error)
	FirstLibrary(id int, userID uint) (*Song, error)
	UpdateId(id int) error
	PlayQueue(id int, userID uint) (*PlayResponse, error)
}

type ServiceSongView interface {
	Creat(userID uint, songID uint) error
	GetByUser(userID uint) ([]SongViewResponse, error)
	Delete(id int) error
	GetByUserAll(userID uint) ([]SongViewResponse, error)
}

type ServiceuserLibrary interface {
	Creat(userID uint, SongID uint) (bool, error)
	DeleteUser(song_id int, user_id int) error
	FingUser(song_id int, user_id int) ([]repository.UserLibrary, error)
}
