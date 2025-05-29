package cmd

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"backend/internal/controller"
	"backend/internal/controller/hello"
	"backend/internal/middleware"
	"backend/internal/service"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.ResponseWrapper)
				group.Bind(
					hello.NewV1(),
					controller.NewWechat(),
				)
			})

			// 需要登录的API组
			s.Group("/api", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.ResponseWrapper)
				group.Middleware(middleware.JwtAuth)
				group.Bind(
					controller.NewProduct(),
					controller.NewOrder(),
					controller.NewUser(),
				)
			})

			s.Group("/admin", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.ResponseWrapper)
				group.Bind(
					controller.NewAdmin(),
					controller.NewAdminUpload(),
					controller.NewAdminProduct(),
					controller.NewAdminCategory(),
				)
			})

			// 添加存储代理路由，用于访问MinIO中的文件
			s.Group("/storage", func(group *ghttp.RouterGroup) {
				// 不使用ResponseWrapper中间件，直接返回文件内容
				group.GET("/*path", func(r *ghttp.Request) {
					path := r.Get("path").String()
					if path == "" {
						r.Response.WriteStatus(http.StatusNotFound)
						return
					}

					// 去掉开头的斜杠
					if strings.HasPrefix(path, "/") {
						path = path[1:]
					}

					// 获取MinIO客户端
					client, err := service.Minio().GetClient(ctx)
					if err != nil {
						r.Response.WriteStatus(http.StatusInternalServerError)
						return
					}

					// 获取存储桶名称
					bucketName := service.Minio().GetBucketName(ctx)

					// 获取对象
					object, err := client.GetObject(ctx, bucketName, path, service.Minio().GetObjectOptions())
					if err != nil {
						r.Response.WriteStatus(http.StatusNotFound)
						return
					}
					defer object.Close()

					// 获取对象信息
					stat, err := object.Stat()
					if err != nil {
						r.Response.WriteStatus(http.StatusNotFound)
						return
					}

					// 设置响应头
					r.Response.Header().Set("Content-Type", stat.ContentType)
					r.Response.Header().Set("Content-Length", g.NewVar(stat.Size).String())
					r.Response.Header().Set("Cache-Control", "max-age=31536000") // 缓存1年

					// 复制文件内容到响应
					io.Copy(r.Response.Writer, object)
				})
			})

			s.Run()
			return nil
		},
	}
)
