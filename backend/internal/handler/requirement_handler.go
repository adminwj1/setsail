package handler

import (
	"projecthub/internal/service"
	"projecthub/pkg/errcode"
	"projecthub/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RequirementHandler struct {
	reqSvc *service.RequirementService
}

func NewRequirementHandler(reqSvc *service.RequirementService) *RequirementHandler {
	return &RequirementHandler{
		reqSvc: reqSvc,
	}
}

// GetAll
// @Summary 获取需求列表
// @Description 分页获取需求列表，支持按产品和状态筛选
// @Tags 需求管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param product_id query int false "产品ID"
// @Param type query int false "需求类型: 1-业务需求, 2-用户需求, 3-研发需求"
// @Param status query int false "需求状态: 0-待评审, 1-已采纳, 2-开发中, 3-已完成, 4-已拒绝"
// @Success 200 {object} response.Response{data=response.PageResult}
// @Failure 500 {object} response.Response
// @Router /api/requirements [get]
func (h *RequirementHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	productID, _ := strconv.Atoi(c.Query("product_id"))
	reqType, _ := strconv.Atoi(c.Query("type"))
	status, _ := strconv.Atoi(c.Query("status"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	requirements, total, err := h.reqSvc.GetAll(page, pageSize, uint(productID), reqType, status)
	if err != nil {
		response.Error(c, errcode.CodeServerError, errcode.ParseError(err))
		return
	}

	response.SuccessPage(c, requirements, total, page, pageSize)
}

// GetByID
// @Summary 获取需求详情
// @Description 根据ID获取需求详细信息
// @Tags 需求管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "需求ID"
// @Success 200 {object} response.Response{data=model.Requirement}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/requirements/{id} [get]
func (h *RequirementHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.CodeBadRequest, "无效的ID")
		return
	}

	req, err := h.reqSvc.GetByID(uint(id))
	if err != nil {
		response.Error(c, errcode.CodeServerError, errcode.ParseError(err))
		return
	}
	if req == nil {
		response.Error(c, errcode.CodeNotFound, "需求不存在")
		return
	}

	response.Success(c, req)
}

// Create
// @Summary 创建需求
// @Description 创建新需求
// @Tags 需求管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateRequirementRequest true "需求信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/requirements [post]
func (h *RequirementHandler) Create(c *gin.Context) {
	var req service.CreateRequirementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.CodeBadRequest, "参数错误")
		return
	}

	userID, _ := c.Get("user_id")
	req.CreatorID = userID.(uint)

	if err := h.reqSvc.Create(&req); err != nil {
		response.Error(c, errcode.CodeBadRequest, errcode.ParseError(err))
		return
	}

	response.SuccessWithMessage(c, "创建成功", nil)
}

// Update
// @Summary 更新需求
// @Description 根据ID更新需求信息
// @Tags 需求管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "需求ID"
// @Param request body service.UpdateRequirementRequest true "需求信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/requirements/{id} [put]
func (h *RequirementHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.CodeBadRequest, "无效的ID")
		return
	}

	var req service.UpdateRequirementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.CodeBadRequest, "参数错误")
		return
	}

	if err := h.reqSvc.Update(uint(id), &req); err != nil {
		response.Error(c, errcode.CodeBadRequest, errcode.ParseError(err))
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

// Delete
// @Summary 删除需求
// @Description 根据ID删除需求
// @Tags 需求管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "需求ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/requirements/{id} [delete]
func (h *RequirementHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.CodeBadRequest, "无效的ID")
		return
	}

	if err := h.reqSvc.Delete(uint(id)); err != nil {
		response.Error(c, errcode.CodeBadRequest, errcode.ParseError(err))
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// GetByProductID
// @Summary 获取产品的需求列表
// @Description 获取指定产品的所有需求
// @Tags 需求管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "产品ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=response.PageResult}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/products/{id}/requirements [get]
func (h *RequirementHandler) GetByProductID(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.CodeBadRequest, "无效的产品ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	requirements, total, err := h.reqSvc.GetByProductID(uint(productID), page, pageSize)
	if err != nil {
		response.Error(c, errcode.CodeServerError, errcode.ParseError(err))
		return
	}

	response.SuccessPage(c, requirements, total, page, pageSize)
}
