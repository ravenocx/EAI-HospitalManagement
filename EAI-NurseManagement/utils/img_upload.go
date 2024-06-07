package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

const CLOUDINARY_URL = "cloudinary://344274639165551:XCldWCJ8c32LVUGC5HudzNShmAE@dmueiy8a3"

func UploadImage(image *multipart.FileHeader) (urlImage string, err error) {
	if image == nil {
		return "", nil
	}
	file, _ := image.Open()

	//format file name of the image
	imgName := image.Filename
	imgFormatIndex := strings.LastIndex(imgName, ".")

	ext := imgName[imgFormatIndex:]

	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return "", errors.New("only JPG and PNG files are allowed")
	}

	formatedName := fmt.Sprintf("%s-%s", time.Now().Format("20060102"), imgName[:imgFormatIndex])

	cld, _ := cloudinary.NewFromURL(CLOUDINARY_URL)
	cld.Config.URL.Secure = true

	context := context.Background()

	resp, err := cld.Upload.Upload(context, file, uploader.UploadParams{
		PublicID:       formatedName,
		UniqueFilename: api.Bool(true),
	})

	if err != nil {
		return "", err
	}

	urlImage = resp.SecureURL
	log.Printf("urlImage : %+s", urlImage)

	return urlImage, nil
}
