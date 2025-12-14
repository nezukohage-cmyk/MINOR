package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cld *cloudinary.Cloudinary

// Initialize Cloudinary once at startup
func InitCloudinary() {
	var err error

	cld, err = cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)

	if err != nil {
		panic("Failed to initialize Cloudinary: " + err.Error())
	}

	fmt.Println("Cloudinary is connected")
}

// FINAL VERSION: UploadFile without gin.Context
func UploadFile(fileHeader *multipart.FileHeader, folder string) (string, error) {

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	uploadResult, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
		Folder:       folder, // dynamic cluster folder
		ResourceType: "raw",  // required for PDFs
	})

	if err != nil {
		fmt.Println("Cloudinary upload error:", err)
		return "", err
	}

	return uploadResult.SecureURL, nil
}
