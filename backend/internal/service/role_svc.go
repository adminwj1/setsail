package service

import (
	"errors"
	"projecthub/internal/model"
	"projecthub/internal/repository"
	"strconv"

	"github.com/casbin/casbin/v3"
)

type RoleService struct {
	roleRepo  *repository.RoleRepository
	menuRepo  *repository.MenuRepository
	enforcer  *casbin.Enforcer
}

func NewRoleService(roleRepo *repository.RoleRepository, menuRepo *repository.MenuRepository, enforcer *casbin.Enforcer) *RoleService {
	return &RoleService{
		roleRepo:  roleRepo,
		menuRepo:  menuRepo,
		enforcer:  enforcer,
	}
}

type CreateRoleRequest struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
	Desc string `json:"desc"`
	Sort int    `json:"sort"`
}

type UpdateRoleRequest struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
	Sort int    `json:"sort"`
}

type SetMenusRequest struct {
	MenuIDs []uint `json:"menu_ids"`
}

func (s *RoleService) Create(req *CreateRoleRequest) error {
	existing, err := s.roleRepo.FindByCode(req.Code)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("角色代码已存在")
	}

	role := &model.Role{
		Name: req.Name,
		Code: req.Code,
		Desc: req.Desc,
		Sort: req.Sort,
		Status: 1,
	}

	return s.roleRepo.Create(role)
}

func (s *RoleService) Update(id uint, req *UpdateRoleRequest) error {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("角色不存在")
	}

	role.Name = req.Name
	role.Desc = req.Desc
	role.Sort = req.Sort

	return s.roleRepo.Update(role)
}

func (s *RoleService) Delete(id uint) error {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("角色不存在")
	}

	return s.roleRepo.Delete(id)
}

func (s *RoleService) GetAll() ([]model.Role, error) {
	return s.roleRepo.FindAll()
}

func (s *RoleService) GetByID(id uint) (*model.Role, error) {
	return s.roleRepo.FindByID(id)
}

func (s *RoleService) SetMenus(id uint, req *SetMenusRequest) error {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("角色不存在")
	}

	// 更新数据库中的菜单关联
	if err := s.roleRepo.SetMenus(id, req.MenuIDs); err != nil {
		return err
	}

	// 同步 Casbin 策略
	if err := s.syncCasbinPolicies(id, req.MenuIDs); err != nil {
		return err
	}

	return nil
}

// syncCasbinPolicies 同步菜单到 Casbin 策略
func (s *RoleService) syncCasbinPolicies(roleID uint, menuIDs []uint) error {
	if s.enforcer == nil {
		return nil // 如果没有 enforcer，跳过同步
	}

	// 获取菜单的 API 路径
	menus, err := s.menuRepo.FindByIDs(menuIDs)
	if err != nil {
		return err
	}

	// 先移除该角色的所有现有策略
	_, _ = s.enforcer.RemoveFilteredPolicy(0, strconv.Itoa(int(roleID)))

	// 为每个菜单添加策略（GET, POST, PUT, DELETE）
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	for _, menu := range menus {
		if menu.ApiPath == "" {
			continue
		}
		for _, method := range methods {
			_, _ = s.enforcer.AddPolicy(strconv.Itoa(int(roleID)), menu.ApiPath, method)
		}
	}

	// 保存策略到数据库
	_ = s.enforcer.SavePolicy()

	return nil
}

func (s *RoleService) GetMenus(id uint) ([]uint, error) {
	return s.roleRepo.GetMenuIDs(id)
}
