package repository

import (
	"projecthub/internal/model"
	"projecthub/pkg/database"

	"gorm.io/gorm"
)

type RoleRepository struct{}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{}
}

func (r *RoleRepository) FindByID(id uint) (*model.Role, error) {
	var role model.Role
	err := database.DB.First(&role, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) FindByCode(code string) (*model.Role, error) {
	var role model.Role
	err := database.DB.Where("code = ?", code).First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) FindAll() ([]model.Role, error) {
	var roles []model.Role
	err := database.DB.Order("sort asc").Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) Create(role *model.Role) error {
	return database.DB.Create(role).Error
}

func (r *RoleRepository) Update(role *model.Role) error {
	return database.DB.Save(role).Error
}

func (r *RoleRepository) Delete(id uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", id).Delete(&model.RoleMenu{}).Error; err != nil {
			return err
		}
		if err := tx.Where("role_id = ?", id).Delete(&model.UserRole{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Role{}, id).Error
	})
}

func (r *RoleRepository) GetMenuIDs(roleID uint) ([]uint, error) {
	var menuIDs []uint
	err := database.DB.
		Model(&model.RoleMenu{}).
		Where("role_id = ?", roleID).
		Pluck("menu_id", &menuIDs).Error
	return menuIDs, err
}

func (r *RoleRepository) SetMenus(roleID uint, menuIDs []uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error; err != nil {
			return err
		}
		for _, menuID := range menuIDs {
			rm := model.RoleMenu{RoleID: roleID, MenuID: menuID}
			if err := tx.Create(&rm).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
