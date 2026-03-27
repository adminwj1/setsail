package handler

import (
	"projecthub/internal/service"
	"projecthub/pkg/errcode"
	"projecthub/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authSvc: authSvc,
	}
}

// Login
// @Summary 用户登录
// @Description 用户使用用户名密码登录，返回JWT token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "登录请求"
// @Success 200 {object} response.Response{data=service.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.CodeBadRequest, "参数错误")
		return
	}

	result, err := h.authSvc.Login(&req)
	if err != nil {
		response.Error(c, errcode.CodeUnauthorized, errcode.ParseError(err))
		return
	}

	response.Success(c, result)
}

// Logout
// @Summary 用户登出
// @Description 用户登出，注销当前会话
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	response.SuccessWithMessage(c, "登出成功", nil)
}

// GetUserInfo
// @Summary 获取当前用户信息
// @Description 获取已登录用户的详细信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=service.UserInfoVO}
// @Failure 401 {object} response.Response
// @Router /api/auth/userinfo [get]
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, errcode.CodeUnauthorized, "未登录")
		return
	}

	userInfo, err := h.authSvc.GetUserInfo(userID.(uint))
	if err != nil {
		response.Error(c, errcode.CodeBadRequest, errcode.ParseError(err))
		return
	}

	response.Success(c, userInfo)
}
