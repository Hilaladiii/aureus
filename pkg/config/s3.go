package config

import (
	"context"
	"fmt"
	"io"

	"github.com/Hilaladiii/aureus/pkg/exception"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type SeaweedFSStorageItf interface {
	UploadFile(ctx context.Context, bucketName string, objectName string, file io.Reader, objectSize int64, contentType string) (string, error)
	MakeBucket(ctx context.Context, bucketName string, region string) error
}

type SeaweedFSStorage struct {
	client *minio.Client
}

func NewSeaweedFSStorage(env Env) *SeaweedFSStorage {
	client, err := minio.New(env.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(env.S3AccessKey, env.S3SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		panic(err)
	}

	return &SeaweedFSStorage{client}
}

func (s *SeaweedFSStorage) UploadFile(ctx context.Context, bucketName string, objectName string, file io.Reader, objectSize int64, contentType string) (string, error) {
	info, err := s.client.PutObject(ctx, bucketName, objectName, file, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	protocol := "http"
	if s.client.EndpointURL().Scheme == "https" {
		protocol = "https"
	}
	fileUrl := fmt.Sprintf("%s://%s/%s/%s", protocol, s.client.EndpointURL().Host, bucketName, info.Key)
	return fileUrl, nil
}

func (s *SeaweedFSStorage) MakeBucket(ctx context.Context, bucketName string, region string) error {
	err := s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
		Region: region,
	})
	if err != nil {
		exists, errBucketExists := s.client.BucketExists(ctx, bucketName)
		fmt.Println(exists)
		if errBucketExists == nil && exists {
			return exception.NewBadRequestError("bucket already exists")
		}

		return err
	}
	return nil
}
