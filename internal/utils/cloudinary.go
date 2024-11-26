package utils

import (
	"context"
	"errors"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/bjorndonald/golang-backend-template/constants"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GenerateUid() string {
	return uuid.New().String()
}

func UploadFile(ctx *gin.Context, file *multipart.FileHeader) (string, error) {
	constant := constants.New()

	if file == nil {
		return "", errors.New("Image is required")
	}

	dst := "assets/uploads/" + GenerateUid() + file.Filename
	defer func() {
		os.Remove(dst)
	}()

	err := ctx.SaveUploadedFile(file, dst)
	if err != nil {
		return "", err
	}

	cld, err := cloudinary.NewFromParams(constant.CloudinaryName, constant.CloudinaryAPIKey, constant.CloudinaryApiSecret)
	if err != nil {
		return "", err
	}
	resp, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{PublicID: strings.ReplaceAll(GenerateUid(), "-", "")})

	if err != nil {
		return "", err
	}
	if resp.Error.Message != "" {
		log.Println(resp.Error.Message)
		return "", errors.New("Cloudinary error")
	}

	return resp.SecureURL, nil
}

func DeleteFile(ctx context.Context, dst string) error {
	constant := constants.New()
	cld, _ := cloudinary.NewFromParams(constant.CloudinaryName, constant.CloudinaryAPIKey, constant.CloudinaryApiSecret)

	_, err := cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: dst})

	if err != nil {
		return err
	}

	return nil
}
