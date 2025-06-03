package minio

import (
	"context"
	"io"
	"mime/multipart"
	"path"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioClient MinIO客户端
type MinioClient struct {
	client     *minio.Client
	endpoint   string
	bucketName string
	useSSL     bool
	domain     string
}

// 全局MinIO客户端
var minioClient *MinioClient

// Init 初始化MinIO客户端
func Init() error {
	endpoint := g.Cfg().MustGet(context.Background(), "minio.endpoint").String()
	accessKey := g.Cfg().MustGet(context.Background(), "minio.accessKey").String()
	secretKey := g.Cfg().MustGet(context.Background(), "minio.secretKey").String()
	bucketName := g.Cfg().MustGet(context.Background(), "minio.bucketName").String()
	useSSL := g.Cfg().MustGet(context.Background(), "minio.useSSL").Bool()
	domain := g.Cfg().MustGet(context.Background(), "minio.domain").String()

	// 创建MinIO客户端
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return err
	}

	// 检查存储桶是否存在
	exists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return err
	}

	// 如果存储桶不存在，则创建
	if !exists {
		err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}

		// 设置存储桶策略，允许公共读取
		policy := `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::` + bucketName + `/*"]
				}
			]
		}`
		err = client.SetBucketPolicy(context.Background(), bucketName, policy)
		if err != nil {
			return err
		}
	} else {
		// 即使存储桶已存在，也重新设置一次存储桶策略，确保公共读取权限
		policy := `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::` + bucketName + `/*"]
				}
			]
		}`
		err = client.SetBucketPolicy(context.Background(), bucketName, policy)
		if err != nil {
			g.Log().Warning(context.Background(), "设置存储桶策略失败: ", err)
		}
	}

	minioClient = &MinioClient{
		client:     client,
		endpoint:   endpoint,
		bucketName: bucketName,
		useSSL:     useSSL,
		domain:     domain,
	}

	return nil
}

// GetClient 获取MinIO客户端
func GetClient() *MinioClient {
	return minioClient
}

// GetRawClient 获取原始的MinIO客户端
func (m *MinioClient) GetRawClient() *minio.Client {
	return m.client
}

// GetBucketName 获取存储桶名称
func (m *MinioClient) GetBucketName() string {
	return m.bucketName
}

// BucketExists 检查存储桶是否存在
func (m *MinioClient) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	return m.client.BucketExists(ctx, bucketName)
}

// GetBucketPolicy 获取存储桶策略
func (m *MinioClient) GetBucketPolicy(ctx context.Context, bucketName string) (string, error) {
	return m.client.GetBucketPolicy(ctx, bucketName)
}

// UploadFile 上传文件
func (m *MinioClient) UploadFile(ctx context.Context, file *multipart.FileHeader, directory string) (string, error) {
	// 打开文件
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 生成唯一文件名
	ext := path.Ext(file.Filename)
	objectName := directory + "/" + uuid.New().String() + ext

	// 上传文件
	_, err = m.client.PutObject(ctx, m.bucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", err
	}

	g.Log().Debug(ctx, "文件上传成功，对象名称: ", objectName)

	// 尝试生成预签名URL，确保文件可访问
	presignedURL, err := m.client.PresignedGetObject(ctx, m.bucketName, objectName, time.Hour*24, nil)
	if err != nil {
		g.Log().Warning(ctx, "生成预签名URL失败: ", err)
	} else {
		g.Log().Debug(ctx, "预签名URL: ", presignedURL.String())
	}

	// 返回代理文件URL，而不是直接的MinIO URL
	serverURL := g.Cfg().MustGet(ctx, "server.url").String()
	if serverURL == "" {
		serverURL = "http://localhost:8000" // 默认服务器URL
	}

	return serverURL + "/api/v1/file/" + objectName, nil
}

// UploadFileFromReader 从Reader上传文件
func (m *MinioClient) UploadFileFromReader(ctx context.Context, reader io.Reader, size int64, contentType, directory, extension string) (string, error) {
	// 生成唯一文件名
	objectName := directory + "/" + uuid.New().String() + extension

	// 上传文件
	_, err := m.client.PutObject(ctx, m.bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	// 返回代理文件URL，而不是直接的MinIO URL
	serverURL := g.Cfg().MustGet(ctx, "server.url").String()
	if serverURL == "" {
		serverURL = "http://localhost:8000" // 默认服务器URL
	}

	return serverURL + "/api/v1/file/" + objectName, nil
}

// DeleteFile 删除文件
func (m *MinioClient) DeleteFile(ctx context.Context, objectName string) error {
	return m.client.RemoveObject(ctx, m.bucketName, objectName, minio.RemoveObjectOptions{})
}

// GeneratePresignedURL 生成预签名URL
func (m *MinioClient) GeneratePresignedURL(ctx context.Context, objectName string, expires time.Duration) (string, error) {
	// 检查对象是否存在
	_, err := m.client.StatObject(ctx, m.bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return "", gerror.New("文件不存在")
	}

	// 生成预签名URL
	presignedURL, err := m.client.PresignedGetObject(ctx, m.bucketName, objectName, expires, nil)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}
