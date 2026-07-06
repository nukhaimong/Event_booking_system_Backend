package config

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	cld *cloudinary.Cloudinary
	ctx context.Context
}

func NewCloudinaryService() (*CloudinaryService, error) {
	// Cloudinary will automatically read from CLOUDINARY_URL env var
	cld, err := cloudinary.New()
	if err != nil {
		return nil, err
	}

	return &CloudinaryService{
		cld: cld,
		ctx: context.Background(),
	}, nil
}

// UploadEventImage uploads an event photo to Cloudinary
func (s *CloudinaryService) UploadEventImage(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// Generate a unique public ID with timestamp
	timestamp := time.Now().Unix()
	publicID := fmt.Sprintf("events/%d_%s", timestamp, fileHeader.Filename)

	// Remove file extension from public ID (optional)
	// publicID = strings.TrimSuffix(publicID, filepath.Ext(fileHeader.Filename))

	// Upload to Cloudinary
	uploadParams := uploader.UploadParams{
		PublicID: publicID,
		Folder:   "events", // Organize in a folder
		Tags:     []string{"event", "photo"},
	}

	resp, err := s.cld.Upload.Upload(s.ctx, file, uploadParams)
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}

// DeleteEventImage deletes an image from Cloudinary (optional)
func (s *CloudinaryService) DeleteEventImage(publicID string) error {
	_, err := s.cld.Upload.Destroy(s.ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}
