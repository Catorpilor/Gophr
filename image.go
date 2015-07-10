package main

import (
	"image"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/disintegration/imaging"
)

const (
	imageIDLength  = 10
	widthThumbnail = 400
	widthPreview   = 800
)

type Image struct {
	ID          string
	UserID      string
	Name        string
	Location    string
	Size        int64
	CreateAt    time.Time
	Description string
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func NewImage(user *User) *Image {
	return &Image{
		ID:       GenerateID("img", imageIDLength),
		UserID:   user.ID,
		CreateAt: time.Now(),
	}
}

//A map of accepted mime types and their file extension

var mimeExtensions = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpeg",
	"image/gif":  ".gif",
}

func (image *Image) CreateFromURL(imageURL string) error {
	//get the response of the url
	res, err := http.Get(imageURL)
	if err != nil {
		return err
	}

	//get the response
	if res.StatusCode != http.StatusOK {
		return errInvalidImageURL
	}

	defer res.Body.Close()

	//Ascertain the type of file downloaded
	mimeType, _, err := mime.ParseMediaType(res.Header.Get("Content-Type"))
	if err != nil {
		return errInvalidImageType
	}

	//Get an extension of the file
	ext, valid := mimeExtensions[mimeType]
	if !valid {
		return errInvalidImageType
	}

	//Get name from url
	image.Name = filepath.Base(imageURL)
	image.Location = image.ID + ext

	//Open file at target location
	savedFile, err := os.Create("./data/images/" + image.Location)
	if err != nil {
		return err
	}
	defer savedFile.Close()

	//Copy the entire res to the output file
	size, err := io.Copy(savedFile, res.Body)
	if err != nil {
		return err
	}
	image.Size = size

	err = image.CreateResizedImages()
	if err != nil {
		return err
	}

	//save to db
	return globalImageStore.Save(image)
}

func (image *Image) CreateFromFile(file multipart.File, headers *multipart.FileHeader) error {
	image.Name = headers.Filename
	image.Location = image.ID + filepath.Ext(image.Name)

	//open file
	savedFile, err := os.Create("./data/images/" + image.Location)
	if err != nil {
		return err
	}

	defer savedFile.Close()

	//copy the uploaded file
	size, err := io.Copy(savedFile, file)
	if err != nil {
		return err
	}

	image.Size = size

	//create resized image
	err = image.CreateResizedImages()
	if err != nil {
		return err
	}

	//save image
	return globalImageStore.Save(image)
}

func (image *Image) ShowRoute() string {
	return "/image/" + image.ID
}

func (image *Image) StaticRoute() string {
	return "/im/" + image.Location
}

func (image *Image) CreateResizedImages() error {
	//generate an image from file
	srcImage, err := imaging.Open("./data/images/" + image.Location)
	if err != nil {
		return err
	}

	//create a channel to receive errors on
	errorChan := make(chan error)

	//process each resize
	go image.resizePreview(errorChan, srcImage)
	go image.resizeThumbnail(errorChan, srcImage)

	//wait for images to finish resizing
	for i := 0; i < 2; i++ {
		err = <-errorChan
		if err != nil {
			return err
		}
	}
	return nil
}

func (image *Image) resizeThumbnail(errorChan chan error, srcImage image.Image) {
	dstImage := imaging.Thumbnail(srcImage, widthThumbnail, widthThumbnail, imaging.Lanczos)
	dest := "./data/images/thumbnail/" + image.Location
	errorChan <- imaging.Save(dstImage, dest)
}

func (image *Image) resizePreview(errorChan chan error, srcImage image.Image) {
	size := srcImage.Bounds().Size()
	ratio := float64(size.Y) / float64(size.X)
	targetHeight := int(float64(widthPreview) * ratio)

	dstImage := imaging.Resize(srcImage, widthPreview, targetHeight, imaging.Lanczos)
	dest := "./data/images/preview/" + image.Location
	errorChan <- imaging.Save(dstImage, dest)
}

func (image *Image) StaticThumbnailRoute() string {
	return "/im/thumbnail/" + image.Location
}

func (image *Image) StaticPreviewRoute() string {
	return "/im/preview/" + image.Location
}
