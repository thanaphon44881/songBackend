package main

import (
	"log"
	"music/adapter"
	"music/repository"
	"music/service"

	"music/utils"

	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("โหลด .env ไม่ได้")
	}
	dsn := os.Getenv("HDAtABACE")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}
	err = db.AutoMigrate(&repository.User{}, &repository.Artist{}, &repository.Song{}, &repository.SongView{}, &repository.UserLibrary{})
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:4200",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept",
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

	port := os.Getenv("PORT")
	log.Fatal(app.Listen(":" + port))
}
