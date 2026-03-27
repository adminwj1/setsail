package handler

import (
	"projecthub/internal/service"
	"projecthub/pkg/errcode"
	"projecthub/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	menuSvc *service.MenuService
}

func NewMenuHandler(menuSvc *service.MenuService) *MenuHandler {
	return &MenuHandler{
		menuSvc: menuSvc,
	}
}

// GetRouter
// @Summary 获取用户动态路由
// @Description 获取当前用户的菜单树，用于前端动态路由渲染
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]model.Menu}
// @Failure 500 {object} response.Response
// @Router /api/menus/router [get]
func (h *MenuHandler) GetRouter(c *gin.Context) {
	userID, _ := c.Get("user_id")
	roleID, _ := c.Get("role_id")
	menus, err := h.menuSvc.GetUserMenuTree(userID.(uint), roleID.(uint))
	if err != nil {
		response.Error(c, errcode.CodeServerError, errcode.ParseError(err))
		return
	}

	response.Success(c, menus)
}

// GetTree
// @Summary 获取菜单树
// @Description 获取所有菜单的树形结构
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]model.Menu}
// @Failure 500 {object} response.Response
// @Router /api/menus/tree [get]
func (h *MenuHandler) GetTree(c *gin.Context) {
	menus, err := h.menuSvc.GetTree()
	if err != nil {
		response.Error(c, errcode.CodeServerError, errcode.ParseError(err))
		return
	}

	response.Success(c, menus)
}

// GetAll
// @Summary 获取所有菜单
// @Description 获取所有菜单列表（扁平结构）
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]model.Menu}
// @Failure 500 {object} response.Response
// @Router /api/menus [get]
func (h *MenuHandler) GetAll(c *gin.Context) {
	menus, err := h.menuSvc.GetAll()
	if err != nil {
		response.Error(c, errcode.CodeServerError, errcode.ParseError(err))
		return
	}

	response.Success(c, menus)
}

// GetByID
// @Summary 获取菜单详情
// @Description 根据ID获取菜单详细信息
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "菜单ID"
// @Success 200 {object} response.Response{data=model.Menu}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/menus/{id} [get]
func (h *MenuHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.CodeBadRequest, "无效的ID")
		return
	}

	menu, err := h.menuSvc.GetByID(uint(id))
	if err != nil {
		response.Error(c, errcode.CodeServerError, errcode.ParseError(err))
		return
	}
	if menu == nil {
		response.Error(c, errcode.CodeNotFound, "菜单不存在")
		return
	}

	response.Success(c, menu)
}

// Create
// @Summary 创建菜单
// @Description 创建新菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateMenuRequest true "菜单信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/menus [post]
func (h *MenuHandler) Create(c *gin.Context) {
	var req service.CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.CodeBadRequest, "参数错误")
		return
	}

	if err := h.menuSvc.Create(&req); err != nil {
		response.Error(c, errcode.CodeBadRequest, errcode.ParseError(err))
		return
	}

	response.SuccessWithMessage(c, "创建成功", nil)
}

// Update
// @Summary 更新菜单
// @Description 根据ID更新菜单信息
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "菜单ID"
// @Param request body service.UpdateMenuRequest true "菜单信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/menus/{id} [put]
func (h *MenuHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.CodeBadRequest, "无效的ID")
		return
	}

	var req service.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.CodeBadRequest, "参数错误")
		return
	}

	if err := h.menuSvc.Update(uint(id), &req); err != nil {
		response.Error(c, errcode.CodeBadRequest, errcode.ParseError(err))
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

// Delete
// @Summary 删除菜单
// @Description 根据ID删除菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "菜单ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/menus/{id} [delete]
func (h *MenuHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.CodeBadRequest, "无效的ID")
		return
	}

	if err := h.menuSvc.Delete(uint(id)); err != nil {
		response.Error(c, errcode.CodeBadRequest, errcode.ParseError(err))
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}
