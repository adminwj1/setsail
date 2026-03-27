package service

import (
	"errors"
	"projecthub/internal/model"
	"projecthub/internal/repository"
)

type ProductService struct {
	productRepo *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

type CreateProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	Status      int8   `json:"status"`
	OwnerID     uint   `json:"owner_id"`
}

type UpdateProductRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Status      int8   `json:"status"`
	OwnerID     uint   `json:"owner_id"`
}

func (s *ProductService) Create(req *CreateProductRequest) error {
	product := &model.Product{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      req.Status,
		OwnerID:     req.OwnerID,
	}

	return s.productRepo.Create(product)
}

func (s *ProductService) Update(id uint, req *UpdateProductRequest) error {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("产品不存在")
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Code != "" {
		product.Code = req.Code
	}
	product.Description = req.Description
	product.Status = req.Status
	product.OwnerID = req.OwnerID

	return s.productRepo.Update(product)
}

func (s *ProductService) Delete(id uint) error {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("产品不存在")
	}

	return s.productRepo.Delete(id)
}

func (s *ProductService) GetByID(id uint) (*model.Product, error) {
	return s.productRepo.FindByID(id)
}

func (s *ProductService) GetAll(page, pageSize int) ([]model.Product, int64, error) {
	return s.productRepo.FindAll(page, pageSize)
}
