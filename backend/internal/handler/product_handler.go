package handler

import (
	"projecthub/internal/service"
	"projecthub/pkg/errcode"
	"projecthub/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productSvc *service.ProductService
}

func NewProductHandler(productSvc *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productSvc: productSvc,
	}
}

// GetAll
// @Summary 获取产品列表
// @Description 分页获取所有产品列表
// @Tags 产品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=response.PageResult}
// @Failure 500 {object} response.Response
// @Router /api/products [get]
func (h *ProductHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	products, total, err := h.productSvc.GetAll(page, pageSize)
	if err != nil {
		response.Error(c, errcode.CodeServerError, errcode.ParseError(err))
		return
	}

	response.SuccessPage(c, products, total, page, pageSize)
}

// GetByID
// @Summary 获取产品详情
// @Description 根据ID获取产品详细信息
// @Tags 产品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "产品ID"
// @Success 200 {object} response.Response{data=model.Product}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/products/{id} [get]
func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.CodeBadRequest, "无效的ID")
		return
	}

	product, err := h.productSvc.GetByID(uint(id))
	if err != nil {
		response.Error(c, errcode.CodeServerError, errcode.ParseError(err))
		return
	}
	if product == nil {
		response.Error(c, errcode.CodeNotFound, "产品不存在")
		return
	}

	response.Success(c, product)
}

// Create
// @Summary 创建产品
// @Description 创建新产品
// @Tags 产品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateProductRequest true "产品信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.CodeBadRequest, "参数错误")
		return
	}

	if err := h.productSvc.Create(&req); err != nil {
		response.Error(c, errcode.CodeBadRequest, errcode.ParseError(err))
		return
	}

	response.SuccessWithMessage(c, "创建成功", nil)
}

// Update
// @Summary 更新产品
// @Description 根据ID更新产品信息
// @Tags 产品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "产品ID"
// @Param request body service.UpdateProductRequest true "产品信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/products/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.CodeBadRequest, "无效的ID")
		return
	}

	var req service.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.CodeBadRequest, "参数错误")
		return
	}

	if err := h.productSvc.Update(uint(id), &req); err != nil {
		response.Error(c, errcode.CodeBadRequest, errcode.ParseError(err))
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

// Delete
// @Summary 删除产品
// @Description 根据ID删除产品
// @Tags 产品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "产品ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.CodeBadRequest, "无效的ID")
		return
	}

	if err := h.productSvc.Delete(uint(id)); err != nil {
		response.Error(c, errcode.CodeBadRequest, errcode.ParseError(err))
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}
