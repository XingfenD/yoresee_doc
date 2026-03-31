package storage

import (
	"context"
	"io"
	"time"

	"github.com/XingfenD/yoresee_doc/pkg/errs"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func NewMinioClient(cfg *MinioOptions) (*minio.Client, error) {
	if cfg == nil {
		return nil, errs.Detail(errs.ErrMinioInitClient, "minio options is nil")
	}
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, errs.Wrap(errs.ErrMinioInitClient, err)
	}
	return client, nil
}

func EnsureBucket(ctx context.Context, client *minio.Client, bucketName string) error {
	if client == nil {
		return errs.Detail(errs.ErrMinioInitClient, "minio client is nil")
	}
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return errs.Wrap(errs.ErrMinioCheckBucketExists, err)
	}
	if !exists {
		if err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return errs.Wrap(errs.ErrMinioCreateBucket, err)
		}
	}
	return nil
}

func UploadFile(bucketName, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error) {
	ctx := context.Background()
	_, err := MinioClient.PutObject(ctx, bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", errs.Wrap(errs.ErrMinioUploadFile, err)
	}

	presignedURL, err := MinioClient.PresignedGetObject(ctx, bucketName, objectName, 7*24*time.Hour, nil)
	if err != nil {
		return "", errs.Wrap(errs.ErrMinioGeneratePresigned, err)
	}

	return presignedURL.String(), nil
}

func GetFile(bucketName, objectName string) (*minio.Object, error) {
	ctx := context.Background()
	return MinioClient.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
}

func DeleteFile(bucketName, objectName string) error {
	ctx := context.Background()
	return MinioClient.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}

func GetPresignedURL(bucketName, objectName string, expires time.Duration) (string, error) {
	ctx := context.Background()
	presignedURL, err := MinioClient.PresignedGetObject(ctx, bucketName, objectName, expires, nil)
	if err != nil {
		return "", errs.Wrap(errs.ErrMinioGeneratePresigned, err)
	}
	return presignedURL.String(), nil
}
