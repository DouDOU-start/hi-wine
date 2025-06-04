package common

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// NewError 创建一个新的错误
func NewError(code int, message string) error {
	return gerror.NewCode(gcode.New(code, message, nil), message)
}
