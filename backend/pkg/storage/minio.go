package storage

import (
	"context"
	"fmt"
	"io"
	"path"
	"strings"
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

func EnsureBucketPublicRead(ctx context.Context, client *minio.Client, bucketName string) error {
	if client == nil {
		return errs.Detail(errs.ErrMinioInitClient, "minio client is nil")
	}
	cleanBucketName := strings.TrimSpace(bucketName)
	if cleanBucketName == "" {
		return errs.Detail(errs.ErrMinioCreateBucket, "bucket name is empty")
	}

	policy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Sid":"PublicReadObjects","Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}`, cleanBucketName)
	if err := client.SetBucketPolicy(ctx, cleanBucketName, policy); err != nil {
		return errs.Wrap(errs.ErrMinioCreateBucket, err)
	}
	return nil
}

func UploadFile(bucketName, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error) {
	if err := PutFile(bucketName, objectName, reader, objectSize, contentType); err != nil {
		return "", errs.Wrap(errs.ErrMinioUploadFile, err)
	}

	ctx := context.Background()
	presignedURL, err := MinioClient.PresignedGetObject(ctx, bucketName, objectName, 7*24*time.Hour, nil)
	if err != nil {
		return "", errs.Wrap(errs.ErrMinioGeneratePresigned, err)
	}

	return presignedURL.String(), nil
}

func PutFile(bucketName, objectName string, reader io.Reader, objectSize int64, contentType string) error {
	ctx := context.Background()
	_, err := MinioClient.PutObject(ctx, bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
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

func BuildPublicObjectPath(objectName string) string {
	cleanObjectName := strings.TrimSpace(objectName)
	if cleanObjectName == "" {
		return ""
	}
	cleanObjectName = strings.TrimLeft(cleanObjectName, "/")
	if strings.HasPrefix(cleanObjectName, "storage/") {
		cleanObjectName = strings.TrimPrefix(cleanObjectName, "storage/")
	}
	if cleanObjectName == "" {
		return ""
	}
	return path.Join("/storage", cleanObjectName)
}
