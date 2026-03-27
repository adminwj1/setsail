package service

import (
	"errors"
	"projecthub/internal/model"
	"projecthub/internal/repository"
)

type RequirementService struct {
	reqRepo *repository.RequirementRepository
}

func NewRequirementService(reqRepo *repository.RequirementRepository) *RequirementService {
	return &RequirementService{
		reqRepo: reqRepo,
	}
}

type CreateRequirementRequest struct {
	ProductID   uint   `json:"product_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Type        int8   `json:"type" binding:"required"`
	Priority    int8   `json:"priority"`
	Description string `json:"description"`
	CreatorID   uint   `json:"creator_id"`
	AssigneeID  uint   `json:"assignee_id"`
}

type UpdateRequirementRequest struct {
	Title       string `json:"title"`
	Type        int8   `json:"type"`
	Priority    int8   `json:"priority"`
	Status      int8   `json:"status"`
	Description string `json:"description"`
	AssigneeID  uint   `json:"assignee_id"`
}

func (s *RequirementService) Create(req *CreateRequirementRequest) error {
	r := &model.Requirement{
		ProductID:   req.ProductID,
		Title:       req.Title,
		Type:        req.Type,
		Priority:    req.Priority,
		Status:      model.ReqStatusPending,
		Description: req.Description,
		CreatorID:   req.CreatorID,
		AssigneeID:  req.AssigneeID,
	}

	if r.Priority == 0 {
		r.Priority = model.ReqPriorityMedium
	}

	return s.reqRepo.Create(r)
}

func (s *RequirementService) Update(id uint, req *UpdateRequirementRequest) error {
	r, err := s.reqRepo.FindByID(id)
	if err != nil {
		return err
	}
	if r == nil {
		return errors.New("需求不存在")
	}

	if req.Title != "" {
		r.Title = req.Title
	}
	if req.Type != 0 {
		r.Type = req.Type
	}
	if req.Priority != 0 {
		r.Priority = req.Priority
	}
	if req.Status != 0 {
		r.Status = req.Status
	}
	r.Description = req.Description
	r.AssigneeID = req.AssigneeID

	return s.reqRepo.Update(r)
}

func (s *RequirementService) Delete(id uint) error {
	r, err := s.reqRepo.FindByID(id)
	if err != nil {
		return err
	}
	if r == nil {
		return errors.New("需求不存在")
	}

	return s.reqRepo.Delete(id)
}

func (s *RequirementService) GetByID(id uint) (*model.Requirement, error) {
	return s.reqRepo.FindByID(id)
}

func (s *RequirementService) GetByProductID(productID uint, page, pageSize int) ([]model.Requirement, int64, error) {
	return s.reqRepo.FindByProductID(productID, page, pageSize)
}

func (s *RequirementService) GetAll(page, pageSize int, productID uint, reqType, status int) ([]model.Requirement, int64, error) {
	return s.reqRepo.FindAll(page, pageSize, productID, reqType, status)
}
