package handler

import (
	"net/http"
	"os"

	"music/adapter"
	"music/service"
	"music/utils"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var app *fiber.App

func init() {
	dsn := os.Getenv("HDAtABACE")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	app = fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://song-forns.vercel.app",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
	}))

	//User
	adapterdb := adapter.NewGormUserDB(db)
	serveruser := service.NewUserService(adapterdb)
	adapterhttp := adapter.NewUserHttps(serveruser)

	app.Post("/user", adapterhttp.CreatUser)
	app.Post("/login", adapterhttp.Login)
	app.Get("/me", utils.AuthRequired, adapterhttp.ShowUsserID)

	//แสดงรูปภาพ
	app.Get("/images/:name", utils.AuthRequired, adapterhttp.ImageShow)
	// app.Static("/images", "./uploads/images")

	//Artist
	artistdb := adapter.NewArtistGormDB(db)
	svArtist := service.NewArtistService(artistdb)
	artisthttps := adapter.NewArtistHttps(svArtist)

	app.Use("/artist", utils.AuthRequired)
	app.Post("/artist", artisthttps.SaveArtist)
	app.Get("/artist/:id", artisthttps.ShowArtist)
	app.Get("/artist", artisthttps.ShowArtistAll)
	app.Get("/artists/:name", utils.AuthRequired, artisthttps.ArtistsImage)

	//Song
	songdb := adapter.NewSongGormDB(db)
	svSong := service.NewSongService(songdb)
	songehttps := adapter.NewSongHttps(svSong)

	app.Use("/song", utils.AuthRequired)
	app.Post("/song", songehttps.SaveSong)
	app.Get("/song", songehttps.ShowSong)
	app.Get("/songnew", songehttps.ShowNew)
	app.Get("/song/play/:id", songehttps.PlaySongNow)
	app.Get("/music/:id", songehttps.PlaySong)
	app.Get("/song/:id/:slug", songehttps.GetSongBySlug)
	app.Get("/cover/:name", utils.AuthRequired, songehttps.CoverImage)

	//SongView
	viewDB := adapter.NewSongViewGormDB(db)
	svview := service.NewSongViewService(viewDB)
	viewHandler := adapter.NewSonViewgHttps(svview)

	app.Use("/history", utils.AuthRequired)
	app.Post("/history", utils.AuthRequired, viewHandler.Create)
	app.Get("/history", utils.AuthRequired, viewHandler.GetMyHistory)
	app.Get("/songview", utils.AuthRequired, viewHandler.GetMyHistoryAll)
	app.Delete("/history/:id", utils.AuthRequired, viewHandler.DeleteSongView)

	//UserLibrary
	userlibraryDB := adapter.NewLibraryGormDB(db)
	svuserlibrary := service.NewuserLibraryService(userlibraryDB)
	userlibraryHandler := adapter.NewuserLibraryHttps(svuserlibrary)

	app.Use("/library", utils.AuthRequired)
	app.Post("/library", userlibraryHandler.SaveLibrary)
	app.Delete("/library/:sonId", userlibraryHandler.DeleteUserLibrary)
	app.Get("/library/:songId", userlibraryHandler.ShowSong)

}

func Handler(w http.ResponseWriter, r *http.Request) {
	adaptor.FiberApp(app)(w, r)
}
