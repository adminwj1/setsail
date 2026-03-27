package repository

import (
	"projecthub/internal/model"
	"projecthub/pkg/database"
)

type RequirementRepository struct{}

func NewRequirementRepository() *RequirementRepository {
	return &RequirementRepository{}
}

func (r *RequirementRepository) FindByID(id uint) (*model.Requirement, error) {
	var req model.Requirement
	err := database.DB.
		Preload("Product").
		Preload("Creator").
		Preload("Assignee").
		First(&req, id).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *RequirementRepository) FindByProductID(productID uint, page, pageSize int) ([]model.Requirement, int64, error) {
	var requirements []model.Requirement
	var total int64

	query := database.DB.Model(&model.Requirement{}).Where("product_id = ?", productID)
	query.Count(&total)

	err := query.
		Preload("Product").
		Preload("Creator").
		Preload("Assignee").
		Offset((page-1)*pageSize).
		Limit(pageSize).
		Find(&requirements).Error
	return requirements, total, err
}

func (r *RequirementRepository) FindAll(page, pageSize int, productID uint, reqType, status int) ([]model.Requirement, int64, error) {
	var requirements []model.Requirement
	var total int64

	query := database.DB.Model(&model.Requirement{})

	if productID > 0 {
		query = query.Where("product_id = ?", productID)
	}
	if reqType > 0 {
		query = query.Where("type = ?", reqType)
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	err := query.
		Preload("Product").
		Preload("Creator").
		Preload("Assignee").
		Offset((page-1)*pageSize).
		Limit(pageSize).
		Find(&requirements).Error
	return requirements, total, err
}

func (r *RequirementRepository) Create(req *model.Requirement) error {
	return database.DB.Create(req).Error
}

func (r *RequirementRepository) Update(req *model.Requirement) error {
	return database.DB.Save(req).Error
}

func (r *RequirementRepository) Delete(id uint) error {
	return database.DB.Delete(&model.Requirement{}, id).Error
}
