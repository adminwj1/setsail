package handler

import (
	"projecthub/internal/service"
	"projecthub/pkg/errcode"
	"projecthub/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleSvc *service.RoleService
}

func NewRoleHandler(roleSvc *service.RoleService) *RoleHandler {
	return &RoleHandler{
		roleSvc: roleSvc,
	}
}

// GetAll
// @Summary 获取角色列表
// @Description 获取所有角色列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]model.Role}
// @Failure 500 {object} response.Response
// @Router /api/roles [get]
func (h *RoleHandler) GetAll(c *gin.Context) {
	roles, err := h.roleSvc.GetAll()
	if err != nil {
		response.Error(c, errcode.CodeServerError, errcode.ParseError(err))
		return
	}

	response.Success(c, roles)
}

// GetByID
// @Summary 获取角色详情
// @Description 根据ID获取角色详细信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response{data=model.Role}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/roles/{id} [get]
func (h *RoleHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.CodeBadRequest, "无效的ID")
		return
	}

	role, err := h.roleSvc.GetByID(uint(id))
	if err != nil {
		response.Error(c, errcode.CodeServerError, errcode.ParseError(err))
		return
	}
	if role == nil {
		response.Error(c, errcode.CodeNotFound, "角色不存在")
		return
	}

	response.Success(c, role)
}

// Create
// @Summary 创建角色
// @Description 创建新角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateRoleRequest true "角色信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/roles [post]
func (h *RoleHandler) Create(c *gin.Context) {
	var req service.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.CodeBadRequest, "参数错误")
		return
	}

	if err := h.roleSvc.Create(&req); err != nil {
		response.Error(c, errcode.CodeBadRequest, errcode.ParseError(err))
		return
	}

	response.SuccessWithMessage(c, "创建成功", nil)
}

// Update
// @Summary 更新角色
// @Description 根据ID更新角色信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Param request body service.UpdateRoleRequest true "角色信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/roles/{id} [put]
func (h *RoleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.CodeBadRequest, "无效的ID")
		return
	}

	var req service.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.CodeBadRequest, "参数错误")
		return
	}

	if err := h.roleSvc.Update(uint(id), &req); err != nil {
		response.Error(c, errcode.CodeBadRequest, errcode.ParseError(err))
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

// Delete
// @Summary 删除角色
// @Description 根据ID删除角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/roles/{id} [delete]
func (h *RoleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.CodeBadRequest, "无效的ID")
		return
	}

	if err := h.roleSvc.Delete(uint(id)); err != nil {
		response.Error(c, errcode.CodeBadRequest, errcode.ParseError(err))
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// GetMenus
// @Summary 获取角色菜单权限
// @Description 获取指定角色关联的菜单ID列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response{data=[]uint}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/roles/{id}/menus [get]
func (h *RoleHandler) GetMenus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.CodeBadRequest, "无效的ID")
		return
	}

	menuIDs, err := h.roleSvc.GetMenus(uint(id))
	if err != nil {
		response.Error(c, errcode.CodeServerError, errcode.ParseError(err))
		return
	}

	response.Success(c, menuIDs)
}

// SetMenus
// @Summary 设置角色菜单权限
// @Description 为指定角色分配菜单权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Param request body service.SetMenusRequest true "菜单ID列表"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/roles/{id}/menus [post]
func (h *RoleHandler) SetMenus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.CodeBadRequest, "无效的ID")
		return
	}

	var req service.SetMenusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.CodeBadRequest, "参数错误")
		return
	}

	if err := h.roleSvc.SetMenus(uint(id), &req); err != nil {
		response.Error(c, errcode.CodeBadRequest, errcode.ParseError(err))
		return
	}

	response.SuccessWithMessage(c, "设置成功", nil)
}
