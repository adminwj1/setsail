package service

import (
	"errors"
	"projecthub/internal/model"
	"projecthub/internal/repository"
)

type MenuService struct {
	menuRepo *repository.MenuRepository
	roleRepo *repository.RoleRepository
}

func NewMenuService(menuRepo *repository.MenuRepository, roleRepo *repository.RoleRepository) *MenuService {
	return &MenuService{
		menuRepo: menuRepo,
		roleRepo: roleRepo,
	}
}

type CreateMenuRequest struct {
	ParentID uint   `json:"parent_id"`
	Name     string `json:"name" binding:"required"`
	Code     string `json:"code"`
	Path     string `json:"path"`
	Icon     string `json:"icon"`
	Type     int8   `json:"type" binding:"required"`
	Sort     int    `json:"sort"`
}

type UpdateMenuRequest struct {
	ParentID uint   `json:"parent_id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Path     string `json:"path"`
	Icon     string `json:"icon"`
	Type     int8   `json:"type"`
	Sort     int    `json:"sort"`
	Status   int8   `json:"status"`
}

func (s *MenuService) GetAll() ([]model.Menu, error) {
	return s.menuRepo.FindAll()
}

func (s *MenuService) GetTree() ([]model.Menu, error) {
	menus, err := s.menuRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return s.menuRepo.BuildTree(menus), nil
}

func (s *MenuService) GetByID(id uint) (*model.Menu, error) {
	return s.menuRepo.FindByID(id)
}

func (s *MenuService) Create(req *CreateMenuRequest) error {
	menu := &model.Menu{
		ParentID: req.ParentID,
		Name:     req.Name,
		Code:     req.Code,
		Path:     req.Path,
		Icon:     req.Icon,
		Type:     req.Type,
		Sort:     req.Sort,
		Status:   1,
	}

	return s.menuRepo.Create(menu)
}

func (s *MenuService) Update(id uint, req *UpdateMenuRequest) error {
	menu, err := s.menuRepo.FindByID(id)
	if err != nil {
		return err
	}
	if menu == nil {
		return errors.New("菜单不存在")
	}

	menu.ParentID = req.ParentID
	menu.Name = req.Name
	menu.Code = req.Code
	menu.Path = req.Path
	menu.Icon = req.Icon
	menu.Type = req.Type
	menu.Sort = req.Sort
	menu.Status = req.Status

	return s.menuRepo.Update(menu)
}

func (s *MenuService) Delete(id uint) error {
	menu, err := s.menuRepo.FindByID(id)
	if err != nil {
		return err
	}
	if menu == nil {
		return errors.New("菜单不存在")
	}

	return s.menuRepo.Delete(id)
}

func (s *MenuService) GetUserMenus(userID uint, roleID uint) ([]model.Menu, error) {
	roleIDs := []uint{roleID}
	return s.menuRepo.FindMenusByRoleIDs(roleIDs)
}

func (s *MenuService) GetUserMenuTree(userID uint, roleID uint) ([]model.Menu, error) {
	menus, err := s.GetUserMenus(userID, roleID)
	if err != nil {
		return nil, err
	}
	return s.menuRepo.BuildTree(menus), nil
}
