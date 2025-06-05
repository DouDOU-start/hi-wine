package cmd

import (
	"backend/api/common"
	v1 "backend/api/user/v1"
	"backend/internal/controller/admin"
	"backend/internal/controller/file"
	"backend/internal/controller/hello"
	"backend/internal/controller/order"
	"backend/internal/controller/print"
	"backend/internal/controller/product"
	"backend/internal/controller/qrcode"
	"backend/internal/controller/user"
	"backend/internal/middleware"
	"backend/internal/utility/minio"
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()

			// 初始化MinIO客户端
			if err := minio.Init(); err != nil {
				g.Log().Fatal(ctx, "MinIO初始化失败: ", err)
			}
			g.Log().Info(ctx, "MinIO初始化成功")

			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(
					middleware.ErrorHandler,
					middleware.ResponseWrapper,
				)

				// 无需认证的接口 - 仅微信登录
				group.Group("/auth", func(group *ghttp.RouterGroup) {
					// 直接绑定微信登录路由
					group.POST("/wechat-login", func(r *ghttp.Request) {
						var req *v1.WechatLoginReq
						if err := r.Parse(&req); err != nil {
							r.Response.WriteJson(g.Map{
								"code":    common.CodeParamError,
								"message": err.Error(),
							})
							r.Exit()
							return
						}

						authCtrl := user.NewAuth()
						res, err := authCtrl.WechatLogin(r.Context(), req)
						if err != nil {
							r.Response.WriteJson(g.Map{
								"code":    common.CodeServerError,
								"message": err.Error(),
							})
							r.Exit()
							return
						}

						r.Response.WriteJson(res)
					})
				})

				// 需要认证的接口
				group.Group("/api", func(group *ghttp.RouterGroup) {
					group.Middleware(middleware.JwtAuth)

					group.Bind(
						hello.NewV1(),
						user.NewV1(), // 用户相关接口（除了登录）
						product.NewV1(),
						order.NewV1(),
						print.NewV1(),
						qrcode.NewV1(),
						file.NewV1(),
					)
				})

				// 管理后台接口
				group.Group("/admin", func(group *ghttp.RouterGroup) {
					group.Middleware(middleware.JwtAuth)

					group.Bind(
						admin.NewV1(),
					)
				})
			})

			s.SetSwaggerUITemplate(ScalarUITemplate)

			s.Run()
			return nil
		},
	}
)

const (
	SwaggerUITemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	<meta name="description" content="SwaggerUI"/>
	<title>SwaggerUI</title>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.10.5/swagger-ui.min.css" />
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.10.5/swagger-ui-bundle.js" crossorigin></script>
<script>
	window.onload = () => {
		window.ui = SwaggerUIBundle({
			url:    '{SwaggerUIDocUrl}',
			dom_id: '#swagger-ui',
		});
	};
</script>
</body>
</html>
`

	OpenapiUITemplate = `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <title>openAPI UI</title>
  </head>
  <body>
    <div id="openapi-ui-container" spec-url="{SwaggerUIDocUrl}" theme="light"></div>
    <script src="https://cdn.jsdelivr.net/npm/openapi-ui-dist@latest/lib/openapi-ui.umd.js"></script>
  </body>
</html>
`

	ScalarUITemplate = `
<!doctype html>
<html>
  <head>
    <title>Scalar API Reference</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <!-- Need a Custom Header? Check out this example https://codepen.io/scalarorg/pen/VwOXqam -->
    <script
      id="api-reference"
      data-url="{SwaggerUIDocUrl}"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>
`
)
