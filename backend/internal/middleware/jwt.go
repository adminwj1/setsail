package middleware

import (
	"net/http"
	"projecthub/pkg/jwt"
	"projecthub/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtUtil *jwt.JWT
}

func NewAuthMiddleware(jwtUtil *jwt.JWT) *AuthMiddleware {
	return &AuthMiddleware{
		jwtUtil: jwtUtil,
	}
}

func (m *AuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorWithStatus(c, http.StatusUnauthorized, 401, "未登录")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.ErrorWithStatus(c, http.StatusUnauthorized, 401, "token格式错误")
			c.Abort()
			return
		}

		claims, err := m.jwtUtil.ParseToken(parts[1])
		if err != nil {
			response.ErrorWithStatus(c, http.StatusUnauthorized, 401, "token已过期")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role_id", claims.RoleID)

		c.Next()
	}
}
