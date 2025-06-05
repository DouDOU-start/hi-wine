package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/golang-jwt/jwt/v4"
)

// JWT自定义声明
type CustomClaims struct {
	Data g.Map `json:"data"`
	jwt.RegisteredClaims
}

// GetJWTSecret 从配置中获取JWT密钥
func GetJWTSecret() []byte {
	ctx := context.Background()
	secret := g.Cfg().MustGet(ctx, "jwt.secret").String()
	if secret == "" {
		// 如果配置中没有设置密钥，使用默认密钥
		return []byte("hi-wine-jwt-secret-key")
	}
	return []byte(secret)
}

// GetJWTExpire 从配置中获取JWT过期时间
func GetJWTExpire() time.Duration {
	ctx := context.Background()
	expire := g.Cfg().MustGet(ctx, "jwt.expire").Int64()
	if expire <= 0 {
		// 如果配置中没有设置过期时间，使用默认值（24小时）
		return 24 * time.Hour
	}
	return time.Duration(expire) * time.Second
}

// GenerateToken 生成JWT令牌
func GenerateToken(data g.Map) (string, error) {
	// 创建声明
	claims := CustomClaims{
		data,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(GetJWTExpire())),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	return token.SignedString(GetJWTSecret())
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证令牌
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// GenerateTestToken 生成测试用的JWT令牌并打印
func GenerateTestToken() {
	data := g.Map{
		"userId": 1,
		"openid": "test_openid",
	}
	token, err := GenerateToken(data)
	if err != nil {
		fmt.Println("Error generating token:", err)
		return
	}
	fmt.Println("Generated test token:")
	fmt.Println(token)
}
