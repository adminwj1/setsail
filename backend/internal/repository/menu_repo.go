package repository

import (
	"projecthub/internal/model"
	"projecthub/pkg/database"

	"gorm.io/gorm"
)

type MenuRepository struct{}

func NewMenuRepository() *MenuRepository {
	return &MenuRepository{}
}

func (r *MenuRepository) FindByID(id uint) (*model.Menu, error) {
	var menu model.Menu
	err := database.DB.First(&menu, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &menu, nil
}

func (r *MenuRepository) FindByIDs(ids []uint) ([]model.Menu, error) {
	var menus []model.Menu
	err := database.DB.Where("id IN ?", ids).Find(&menus).Error
	return menus, err
}

func (r *MenuRepository) FindAll() ([]model.Menu, error) {
	var menus []model.Menu
	err := database.DB.Order("sort asc").Find(&menus).Error
	return menus, err
}

func (r *MenuRepository) FindByParentID(parentID uint) ([]model.Menu, error) {
	var menus []model.Menu
	err := database.DB.Where("parent_id = ?", parentID).Order("sort asc").Find(&menus).Error
	return menus, err
}

func (r *MenuRepository) FindMenusByRoleIDs(roleIDs []uint) ([]model.Menu, error) {
	var menus []model.Menu
	err := database.DB.
		Distinct("sys_menu.*").
		Joins("JOIN sys_role_menu ON sys_role_menu.menu_id = sys_menu.id").
		Where("sys_role_menu.role_id IN ?", roleIDs).
		Where("sys_menu.status = 1").
		Order("sys_menu.sort asc").
		Find(&menus).Error
	return menus, err
}

func (r *MenuRepository) FindMenuIDsByRoleIDs(roleIDs []uint) ([]uint, error) {
	var menuIDs []uint
	err := database.DB.
		Model(&model.Menu{}).
		Distinct("sys_menu.id").
		Joins("JOIN sys_role_menu ON sys_role_menu.menu_id = sys_menu.id").
		Where("sys_role_menu.role_id IN ?", roleIDs).
		Where("sys_menu.status = 1").
		Pluck("sys_menu.id", &menuIDs).Error
	return menuIDs, err
}

func (r *MenuRepository) Create(menu *model.Menu) error {
	return database.DB.Create(menu).Error
}

func (r *MenuRepository) Update(menu *model.Menu) error {
	return database.DB.Save(menu).Error
}

func (r *MenuRepository) Delete(id uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("parent_id = ?", id).Delete(&model.Menu{}).Error; err != nil {
			return err
		}
		if err := tx.Where("menu_id = ?", id).Delete(&model.RoleMenu{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Menu{}, id).Error
	})
}

func (r *MenuRepository) BuildTree(menus []model.Menu) []model.Menu {
	menuMap := make(map[uint]*model.Menu)
	var roots []model.Menu

	for i := range menus {
		menus[i].Children = []model.Menu{}
		menuMap[menus[i].ID] = &menus[i]
	}

	for i := range menus {
		if menus[i].ParentID == 0 {
			roots = append(roots, menus[i])
		} else {
			if parent, ok := menuMap[menus[i].ParentID]; ok {
				parent.Children = append(parent.Children, menus[i])
			}
		}
	}

	return roots
}
