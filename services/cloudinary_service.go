package services

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cld *cloudinary.Cloudinary

// Initialize Cloudinary connection once at startup
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

	fmt.Println("Cloudinary jootege sambanda sthapane aagide")
}

func UploadFile(filePath string, subjectIDs []string, topicIDs []string, userID string) (string, string, error) {

	// Open file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	// Pick the first subject/topic safely
	subject := "unknown_subject"
	topic := "unknown_topic"

	if len(subjectIDs) > 0 {
		subject = subjectIDs[0]
	}
	if len(topicIDs) > 0 {
		topic = topicIDs[0]
	}

	// Clean folder structure
	uploadFolder := fmt.Sprintf("notes/%s/%s/%s", subject, topic, userID)

	// Upload using Cloudinary uploader
	uploadResult, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
		Folder: uploadFolder,
	})

	if err != nil {
		fmt.Println("Cloudinary Upload Error:", err)
		return "", "", err
	}

	fmt.Printf("Cloudinary Upload Response: %+v\n", uploadResult)

	return uploadResult.SecureURL, uploadResult.PublicID, nil
}
