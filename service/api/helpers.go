package api

import (
	"WasaText/service/database"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId uint `json:"id"`
}

func GenerateSessionToken(user *database.User) (string, error) {
	ttl, err := strconv.Atoi(os.Getenv("TTL_HOUR"))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(ttl) * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
	})

	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

func ParseUserToken(accessToken string) (uint, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenCliams")
	}

	return claims.UserId, nil
}

func SaveUploadedFile(file multipart.File, header *multipart.FileHeader, uploadDir string, userId uint) (string, error) {
	filename := fmt.Sprintf("%d_%s", userId, filepath.Base(header.Filename))

	err := os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	// Define file path
	filePath := filepath.Join(uploadDir, filename)

	// Create  file
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Copy the uploaded file's content to  destination file
	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
