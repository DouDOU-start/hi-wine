package cmd

import (
	"backend/api/admin"
	ctrlAdmin "backend/internal/controller/admin"
	"backend/internal/controller/file"
	"backend/internal/controller/order"
	"backend/internal/controller/print"
	"backend/internal/controller/product"
	"backend/internal/controller/qrcode"
	"backend/internal/controller/user"
	"backend/internal/middleware"
	"backend/internal/service"
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

			// 初始化定时任务
			if err := service.Cron().StartCronJobs(); err != nil {
				g.Log().Fatal(ctx, "定时任务初始化失败: ", err)
			}
			g.Log().Info(ctx, "定时任务初始化成功")

			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(
					middleware.ErrorHandler,
					middleware.ResponseWrapper,
				)

				// 公开API - 无需认证 - 商品查询和套餐查询
				group.Group("/api/public", func(group *ghttp.RouterGroup) {
					// 商品相关API
					group.Bind(
						product.NewV1(),
					)

					// 套餐查询相关API
					userController := user.NewV1()
					group.Bind(userController.UserPackageList)
					group.Bind(userController.UserPackageDetail)
				})

				// 文件代理API - 公开访问
				group.Group("/", func(group *ghttp.RouterGroup) {
					fileController := file.NewV1()
					group.Bind(fileController.FileProxyPublic)
				})

				// 用户登录接口 - 无需认证
				group.Group("/", func(group *ghttp.RouterGroup) {
					// 使用AuthController进行登录处理
					authCtrl := user.NewAuth()
					group.Bind(authCtrl.WechatLogin)
				})

				// 管理员登录接口 - 无需认证
				group.Group("/api/v1/admin", func(group *ghttp.RouterGroup) {
					// 使用AdminAuthController进行登录处理
					adminAuthCtrl := ctrlAdmin.NewAuth()
					group.Bind(adminAuthCtrl.AdminLogin)
				})

				// 需要认证的接口
				group.Group("/api/v1", func(group *ghttp.RouterGroup) {
					group.Middleware(middleware.JwtAuth)

					// 用户相关API (需要认证)
					group.Group("/user", func(group *ghttp.RouterGroup) {
						userController := user.NewV1()
						group.Bind(userController.UserBuyPackage)
						group.Bind(userController.UserProfile)
						group.Bind(userController.UpdateUserProfile)
						group.Bind(userController.UserOrderList)
						group.Bind(userController.UserOrderDetail)
						group.Bind(userController.UserMyPackages)
					})

					// 其他需要认证的API
					group.Bind(
						order.NewV1(),
					)

					// 管理后台接口 (需要认证)
					group.Group("/admin", func(group *ghttp.RouterGroup) {
						// 排除登录接口，因为它已经在前面定义了
						adminController := ctrlAdmin.NewV1()

						// 绑定所有其他管理员接口
						group.Bind(
							qrcode.NewV1(),
							print.NewV1(),
						)

						// 手动绑定管理员控制器，排除登录方法
						bindAdminControllerExceptLogin(group, adminController)
					})
				})
			})

			s.SetSwaggerUITemplate(ScalarUITemplate)

			s.Run()
			return nil
		},
	}
)

// bindAdminControllerExceptLogin 绑定管理员控制器除登录外的所有方法
func bindAdminControllerExceptLogin(group *ghttp.RouterGroup, controller admin.IAdminV1) {
	// 绑定分类相关接口
	group.Bind(controller.CategoryList)
	group.Bind(controller.CreateCategory)
	group.Bind(controller.UpdateCategory)
	group.Bind(controller.DeleteCategory)
	group.Bind(controller.CategoryDetail)

	// 绑定订单相关接口
	group.Bind(controller.AdminOrderList)
	group.Bind(controller.AdminOrderDetail)
	group.Bind(controller.AdminOrderUpdateStatus)

	// 绑定套餐相关接口
	group.Bind(controller.AdminPackageList)
	group.Bind(controller.AdminPackageCreate)
	group.Bind(controller.AdminPackageUpdate)
	group.Bind(controller.AdminPackageDelete)
	group.Bind(controller.AdminPackageAddProducts)
	group.Bind(controller.AdminPackageRemoveProduct)
	group.Bind(controller.AdminPackageProductList)

	// 绑定商品相关接口
	group.Bind(controller.AdminProductList)
	group.Bind(controller.AdminProductCreate)
	group.Bind(controller.AdminProductUpdate)
	group.Bind(controller.AdminProductDelete)
	group.Bind(controller.AdminProductDetail)

	// 绑定用户套餐相关接口
	group.Bind(controller.AdminUserPackageList)

	// 绑定上传文件接口
	group.Bind(controller.UploadFile)
}

const (
	SwaggerUITemplate = `
<!DOCTYPE HTML>
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
