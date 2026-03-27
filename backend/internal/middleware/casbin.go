package middleware

import (
	"net/http"
	"strconv"

	"projecthub/pkg/response"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
)

type CasbinMiddleware struct {
	enforcer *casbin.Enforcer
}

func NewCasbinMiddleware(e *casbin.Enforcer) *CasbinMiddleware {
	return &CasbinMiddleware{
		enforcer: e,
	}
}

func (m *CasbinMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, exists := c.Get("role_id")
		if !exists {
			response.ErrorWithStatus(c, http.StatusUnauthorized, 401, "未分配角色")
			c.Abort()
			return
		}

		obj := c.Request.URL.Path
		act := c.Request.Method

		ok, err := m.enforcer.Enforce(strconv.Itoa(int(roleID.(uint))), obj, act)
		if err != nil {
			response.ErrorWithStatus(c, http.StatusInternalServerError, 500, "权限检查失败")
			c.Abort()
			return
		}

		if !ok {
			response.ErrorWithStatus(c, http.StatusForbidden, 403, "无权限访问")
			c.Abort()
			return
		}

		c.Next()
	}
}
