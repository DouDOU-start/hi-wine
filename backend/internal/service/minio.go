package service

import (
	"context"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// IMinioService MinIO服务接口
type IMinioService interface {
	// Upload 上传文件
	Upload(ctx context.Context, fileReader io.Reader, fileSize int64, fileName, fileType string) (string, error)
	// GetURL 获取文件URL
	GetURL(ctx context.Context, objectName string) (string, error)
	// Delete 删除文件
	Delete(ctx context.Context, objectName string) error
	// GetClient 获取MinIO客户端
	GetClient(ctx context.Context) (*minio.Client, error)
	// GetBucketName 获取存储桶名称
	GetBucketName(ctx context.Context) string
	// GetObjectOptions 获取对象选项
	GetObjectOptions() minio.GetObjectOptions
}

var minioService = localMinioService{}

type localMinioService struct{}

func Minio() IMinioService {
	return &minioService
}

// getClient 获取MinIO客户端
func (s *localMinioService) getClient(ctx context.Context) (*minio.Client, error) {
	endpoint := g.Cfg().MustGet(ctx, "minio.endpoint").String()
	accessKey := g.Cfg().MustGet(ctx, "minio.accessKey").String()
	secretKey := g.Cfg().MustGet(ctx, "minio.secretKey").String()
	useSSL := g.Cfg().MustGet(ctx, "minio.useSSL").Bool()

	// 初始化MinIO客户端
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, gerror.Wrap(err, "MinIO客户端初始化失败")
	}

	return client, nil
}

// getBucketName 获取存储桶名称
func (s *localMinioService) getBucketName(ctx context.Context) string {
	return g.Cfg().MustGet(ctx, "minio.bucket").String()
}

// ensureBucketExists 确保存储桶存在
func (s *localMinioService) ensureBucketExists(ctx context.Context, client *minio.Client, bucketName string) error {
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return gerror.Wrap(err, "检查存储桶是否存在失败")
	}

	if !exists {
		region := g.Cfg().MustGet(ctx, "minio.region").String()
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
			Region: region,
		})
		if err != nil {
			return gerror.Wrap(err, "创建存储桶失败")
		}

		// 设置存储桶策略，允许公共读取
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::%s/*"]
				}
			]
		}`, bucketName)

		err = client.SetBucketPolicy(ctx, bucketName, policy)
		if err != nil {
			return gerror.Wrap(err, "设置存储桶策略失败")
		}
	}

	return nil
}

// Upload 上传文件
func (s *localMinioService) Upload(ctx context.Context, fileReader io.Reader, fileSize int64, fileName, fileType string) (string, error) {
	client, err := s.getClient(ctx)
	if err != nil {
		return "", err
	}

	bucketName := s.getBucketName(ctx)
	err = s.ensureBucketExists(ctx, client, bucketName)
	if err != nil {
		return "", err
	}

	// 生成文件名
	fileExt := path.Ext(fileName)
	objectName := fmt.Sprintf("uploads/%s/%s%s",
		gtime.Now().Format("Y/m/d"),
		guid.S(),
		fileExt)

	// 上传文件
	_, err = client.PutObject(ctx, bucketName, objectName, fileReader, fileSize, minio.PutObjectOptions{
		ContentType: fileType,
	})
	if err != nil {
		return "", gerror.Wrap(err, "上传文件失败")
	}

	return objectName, nil
}

// GetURL 获取文件URL
func (s *localMinioService) GetURL(ctx context.Context, objectName string) (string, error) {
	if objectName == "" {
		return "", nil
	}

	// 如果是完整URL，直接返回
	if strings.HasPrefix(objectName, "http://") || strings.HasPrefix(objectName, "https://") {
		return objectName, nil
	}

	domain := g.Cfg().MustGet(ctx, "minio.domain").String()
	bucketName := s.getBucketName(ctx)
	serverUrl := g.Cfg().MustGet(ctx, "server.url").String()

	// 如果配置了服务器URL，优先使用服务器URL作为代理
	if serverUrl != "" {
		// 构建URL，通过后端服务器代理访问MinIO
		url := fmt.Sprintf("%s/storage/%s", serverUrl, objectName)
		return url, nil
	}

	// 否则直接使用MinIO的域名
	url := fmt.Sprintf("%s/%s/%s", domain, bucketName, objectName)
	return url, nil
}

// Delete 删除文件
func (s *localMinioService) Delete(ctx context.Context, objectName string) error {
	// 如果是完整URL，提取对象名称
	if strings.HasPrefix(objectName, "http://") || strings.HasPrefix(objectName, "https://") {
		domain := g.Cfg().MustGet(ctx, "minio.domain").String()
		bucketName := s.getBucketName(ctx)
		prefix := fmt.Sprintf("%s/%s/", domain, bucketName)
		if strings.HasPrefix(objectName, prefix) {
			objectName = strings.TrimPrefix(objectName, prefix)
		} else {
			return gerror.New("无法从URL中提取对象名称")
		}
	}

	client, err := s.getClient(ctx)
	if err != nil {
		return err
	}

	bucketName := s.getBucketName(ctx)
	err = client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return gerror.Wrap(err, "删除文件失败")
	}

	return nil
}

// GetClient 获取MinIO客户端
func (s *localMinioService) GetClient(ctx context.Context) (*minio.Client, error) {
	return s.getClient(ctx)
}

// GetBucketName 获取存储桶名称
func (s *localMinioService) GetBucketName(ctx context.Context) string {
	return s.getBucketName(ctx)
}

// GetObjectOptions 获取对象选项
func (s *localMinioService) GetObjectOptions() minio.GetObjectOptions {
	return minio.GetObjectOptions{}
}
