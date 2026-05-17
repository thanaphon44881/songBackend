package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"regexp"

	supa "github.com/supabase-community/supabase-go"
)

func UploadArtist(file *multipart.FileHeader) (string, error) {
	client, err := supa.NewClient(
		os.Getenv("SUPABASE_URL"),
		os.Getenv("SUPABASE_KEY"),
		nil,
	)
	if err != nil {
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)

	_, err = client.Storage.UploadFile("artist", filename, src)
	if err != nil {
		return "", err
	}

	url := os.Getenv("SUPABASE_URL") + "/storage/v1/object/public/artist/" + filename
	return url, nil
}

func UploadSong(filecover *multipart.FileHeader, filesongs *multipart.FileHeader) (string, string, error) {
	client, err := supa.NewClient(
		os.Getenv("SUPABASE_URL"),
		os.Getenv("SUPABASE_KEY"),
		nil,
	)
	if err != nil {
		return "", "", err
	}

	songSrc, err := filesongs.Open()
	if err != nil {
		return "", "", err
	}
	defer songSrc.Close()

	coverSrc, err := filecover.Open()
	if err != nil {
		return "", "", err
	}
	defer coverSrc.Close()

	filenamec := fmt.Sprintf("%d_%s", time.Now().Unix(), filecover.Filename)
	_, err = client.Storage.UploadFile("cover", filenamec, coverSrc)

	safe := regexp.MustCompile(`[^a-zA-Z0-9._-]`).ReplaceAllString(filesongs.Filename, "")
	filenames := fmt.Sprintf("%d_%s", time.Now().Unix(), safe)

	_, err = client.Storage.UploadFile("Songs", filenames, songSrc)
	if err != nil {
		return "", "", err
	}

	urls := os.Getenv("SUPABASE_URL") + "/storage/v1/object/public/Songs/" + filenames
	urlc := os.Getenv("SUPABASE_URL") + "/storage/v1/object/public/cover/" + filenames
	return string(urlc), string(urls), nil
}
