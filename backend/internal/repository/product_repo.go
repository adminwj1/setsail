package repository

import (
	"projecthub/internal/model"
	"projecthub/pkg/database"
)

type ProductRepository struct{}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

func (r *ProductRepository) FindByID(id uint) (*model.Product, error) {
	var product model.Product
	err := database.DB.Preload("Owner").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) FindAll(page, pageSize int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := database.DB.Model(&model.Product{})
	query.Count(&total)

	err := query.Preload("Owner").Offset((page-1)*pageSize).Limit(pageSize).Find(&products).Error
	return products, total, err
}

func (r *ProductRepository) Create(product *model.Product) error {
	return database.DB.Create(product).Error
}

func (r *ProductRepository) Update(product *model.Product) error {
	return database.DB.Save(product).Error
}

func (r *ProductRepository) Delete(id uint) error {
	return database.DB.Delete(&model.Product{}, id).Error
}
