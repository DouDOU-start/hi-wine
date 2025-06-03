package file

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	v1 "backend/api/admin/v1"
	"backend/internal/utility/minio"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	miniogo "github.com/minio/minio-go/v7"
)

// FileProxy 文件代理，通过服务端访问MinIO中的文件
func (c *ControllerV1) FileProxy(ctx context.Context, req *v1.FileProxyReq) (res *ghttp.Response, err error) {
	// 获取MinIO客户端
	minioClient := minio.GetClient()
	if minioClient == nil {
		return nil, gerror.New("文件存储服务未初始化")
	}

	// 获取存储桶名称
	bucketName := minioClient.GetBucketName()

	// 清理和验证文件路径
	request := g.RequestFromCtx(ctx)
	path := strings.TrimPrefix(request.URL.Path, "/api/v1/file/")
	if path == "" {
		return nil, gerror.New("无效的文件路径")
	}

	g.Log().Debug(ctx, "文件代理请求路径: ", path, ", 存储桶: ", bucketName)

	// 检查文件是否存在
	rawClient := minioClient.GetRawClient()
	_, err = rawClient.StatObject(ctx, bucketName, path, miniogo.StatObjectOptions{})
	if err != nil {
		g.Log().Error(ctx, "文件不存在: ", err, ", 路径: ", path, ", 存储桶: ", bucketName)
		return nil, gerror.New("文件不存在或无法访问")
	}

	// 获取文件对象
	object, err := rawClient.GetObject(ctx, bucketName, path, miniogo.GetObjectOptions{})
	if err != nil {
		g.Log().Error(ctx, "获取文件失败: ", err)
		return nil, gerror.New("获取文件失败")
	}
	defer object.Close()

	// 获取文件信息
	stat, err := object.Stat()
	if err != nil {
		g.Log().Error(ctx, "获取文件信息失败: ", err)
		return nil, gerror.New("获取文件信息失败")
	}

	g.Log().Debug(ctx, "文件信息: ", stat.Key, ", 类型: ", stat.ContentType, ", 大小: ", stat.Size)

	// 获取当前的HTTP响应对象
	response := request.Response

	// 设置响应头
	response.Header().Set("Content-Type", stat.ContentType)
	response.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size))
	response.Header().Set("Content-Disposition", "inline; filename="+stat.Key)
	response.Header().Set("Cache-Control", "max-age=31536000") // 缓存1年
	response.Header().Set("Expires", time.Now().AddDate(1, 0, 0).Format(http.TimeFormat))

	// 将文件内容写入响应
	_, err = io.Copy(response.Writer, object)
	if err != nil {
		g.Log().Error(ctx, "写入响应失败: ", err)
		return nil, gerror.New("写入响应失败")
	}

	g.Log().Debug(ctx, "文件代理成功: ", path)

	// 返回nil，因为我们已经手动处理了响应
	return nil, nil
}
