package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func InitMinio(cfg *config.MinioConfig) error {
	var err error
	MinioClient, err = minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("init minio client failed: %w", err)
	}

	// 检查并创建存储桶
	exists, err := MinioClient.BucketExists(context.Background(), cfg.Bucket)
	if err != nil {
		return fmt.Errorf("check bucket exists failed: %w", err)
	}
	if !exists {
		err = MinioClient.MakeBucket(context.Background(), cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("create bucket failed: %w", err)
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
		return "", fmt.Errorf("upload file failed: %w", err)
	}

	presignedURL, err := MinioClient.PresignedGetObject(ctx, bucketName, objectName, 7*24*time.Hour, nil)
	if err != nil {
		return "", fmt.Errorf("generate presigned url failed: %w", err)
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

// GetPresignedURL 获取文件的预签名URL
func GetPresignedURL(bucketName, objectName string, expires time.Duration) (string, error) {
	ctx := context.Background()
	presignedURL, err := MinioClient.PresignedGetObject(ctx, bucketName, objectName, expires, nil)
	if err != nil {
		return "", fmt.Errorf("generate presigned url failed: %w", err)
	}
	return presignedURL.String(), nil
}
