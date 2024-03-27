package service

import (
	"errors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/net/context"
	"mime/multipart"
	"net/http"
	"os"
)

var FileNotFound = errors.New("file not found in s3")

func getClientS3() (*minio.Client, error) {
	return minio.New(
		os.Getenv("BUCKET_ENDPOINT"),
		&minio.Options{
			Creds: credentials.NewStaticV4(
				os.Getenv("ACCESS_ID"), os.Getenv("SECRET_KEY"), "",
			),
			Region: os.Getenv("REGION"),
			Secure: false,
		},
	)
}

func checkObjectExist(minioClient *minio.Client, key string) (bool, error) {
	_, err := minioClient.StatObject(
		context.TODO(),
		os.Getenv("BUCKET_NAME"),
		key,
		minio.StatObjectOptions{},
	)

	if err != nil {
		if minio.ToErrorResponse(err).StatusCode == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func UploadFile(file *multipart.FileHeader, keyName string) error {
	minioClient, err := getClientS3()
	if err != nil {
		return err
	}

	f, err := file.Open()
	if err != nil {
		return err
	}

	_, err = minioClient.PutObject(
		context.TODO(),
		os.Getenv("BUCKET_NAME"),
		keyName,
		f,
		file.Size,
		minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")},
	)

	if err != nil {
		return err
	}

	return nil
}

func DownloadFile(key string) (*minio.Object, error) {
	minioClient, err := getClientS3()
	if err != nil {
		return nil, err
	}

	status, err := checkObjectExist(minioClient, key)
	if err != nil {
		return nil, err
	}

	if !status {
		return nil, FileNotFound
	}

	object, err := minioClient.GetObject(
		context.TODO(),
		os.Getenv("BUCKET_NAME"),
		key,
		minio.GetObjectOptions{},
	)

	if err != nil {
		return nil, err
	}

	return object, nil
}

func DeleteFile(key string) error {
	minioClient, err := getClientS3()
	if err != nil {
		return err
	}

	status, err := checkObjectExist(minioClient, key)
	if err != nil {
		return err
	}

	if !status {
		return FileNotFound
	}

	err = minioClient.RemoveObject(
		context.TODO(),
		os.Getenv("BUCKET_NAME"),
		key,
		minio.RemoveObjectOptions{},
	)

	if err != nil {
		return err
	}

	return nil
}
