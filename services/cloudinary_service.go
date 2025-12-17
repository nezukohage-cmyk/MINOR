package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"strings"

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

// UploadFile uploads PDFs correctly as RAW assets
func UploadFile(fileHeader *multipart.FileHeader, folder string) (string, error) {
	fmt.Println("🔥 UPLOAD USING RAW RESOURCE TYPE")

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	if !strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".pdf") {
		return "", fmt.Errorf("only PDF files are allowed")
	}

	//useFilename := true
	//uniqueFilename := false

	uploadResult, err := cld.Upload.Upload(
		context.Background(),
		file,
		uploader.UploadParams{
			Folder:       folder,
			ResourceType: "raw",
			//UseFilename:    &useFilename,
			//UniqueFilename: &uniqueFilename,
		},
	)

	if err != nil {
		fmt.Println("Cloudinary upload error:", err)
		return "", err
	}

	return uploadResult.SecureURL, nil
}
func DeleteFromCloudinary(publicID string) error {
	_, err := cld.Upload.Destroy(
		context.Background(),
		uploader.DestroyParams{
			PublicID:     publicID,
			ResourceType: "raw", // REQUIRED for PDFs
		},
	)
	return err
}
