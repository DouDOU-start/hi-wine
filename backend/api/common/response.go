package common

// 响应码常量
const (
	CodeSuccess      = 200 // 成功
	CodeParamError   = 400 // 参数错误
	CodeUnauthorized = 401 // 未授权
	CodeForbidden    = 403 // 禁止访问
	CodeNotFound     = 404 // 资源不存在
	CodeServerError  = 500 // 服务器错误
)

// Response 统一响应结构体
// code: 业务码，message: 提示信息，data: 业务数据
// 所有接口响应均应嵌套此结构体
type Response[T any] struct {
	Code    int    `json:"code" example:"200" description:"业务码"`
	Message string `json:"message" example:"success" description:"提示信息"`
	Data    T      `json:"data" description:"业务数据"`
}
