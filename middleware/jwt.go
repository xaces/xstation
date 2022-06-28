package middleware

import (
	"errors"
	"strings"
	"time"

	"github.com/xaces/xutils/ctx"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// UserClaims 用户Claims
type UserClaims struct {
	UserID   uint64 `json:"userId"`
	RoleID   uint64 `json:"roleId"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// JWT 定义
type JWT struct {
	SigningKey []byte
}

var (
	TokenMalformed                  error = errors.New("Not Event Token")
	TokenValidationErrorExpired     error = errors.New("Validation Error Expired")
	TokenValidationErrorNotValidYet error = errors.New("Validation Error Not Valid Yet")
	TokenInValidation               error = errors.New("Token Validation")
)

func signKey() string {
	return "wanguandong@126.com"
}

// NewJWT 创建JWT
func NewJWT() *JWT {
	return &JWT{
		[]byte(signKey()),
	}
}

// createToken 创建Token
func (j *JWT) createToken(claims UserClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析Token
func (j *JWT) ParseToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenValidationErrorExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenValidationErrorNotValidYet
			} else {
				return nil, TokenInValidation
			}
		}
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInValidation
}

// RefreshToken 刷新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.createToken(*claims)
	}
	return "", TokenInValidation
}

// GenerateToken 生成Token
func GenerateToken(userID, roleID uint64, name string) (string, error) {
	j := NewJWT()
	claims := UserClaims{
		userID,
		roleID,
		name,
		jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000), // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 3600), // 过期时间 一小时
			Issuer:    string(j.SigningKey),            // 签名的发行者
		},
	}
	return j.createToken(claims)
}

// JWTAuth 鉴权
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// jwt鉴权取头部信息 x-token
		// 登录时回返回token信息
		// 前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := c.Request.Header.Get("Authorization")
		if !strings.Contains(token, "Bearer ") {
			ctx.JSON(ctx.StatusLoginExpired).SetMsg("Invalid token").WriteTo(c)
			c.Abort()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		j := NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			// if err == middleware.TokenValidationErrorExpired {
			// 	c.Abort()
			// 	return
			// }
			ctx.JSON(ctx.StatusLoginExpired).SetMsg(err.Error()).WriteTo(c)
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

// GetUserOfToken 根据Token获取用户信息
func GetUserOfToken(c *gin.Context) (*UserClaims, error) {
	claims, noFound := c.Get("claims")
	if !noFound {
		return nil, TokenMalformed
	}
	userClaims, _ := claims.(*UserClaims)
	return userClaims, nil
}

// GetUserNameOfToken 根据Token获取用户名
func GetUserNameOfToken(c *gin.Context) string {
	claims, noFound := c.Get("claims")
	if !noFound {
		return ""
	}
	userClaims, _ := claims.(*UserClaims)
	return userClaims.Username
}
