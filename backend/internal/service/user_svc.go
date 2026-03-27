package service

import (
	"errors"
	"projecthub/internal/model"
	"projecthub/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	RoleID   uint   `json:"role_id"`
}

type UpdateUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   int8   `json:"status"`
	RoleID   uint   `json:"role_id"`
}

type UserVO struct {
	ID       uint      `json:"id"`
	Username string    `json:"username"`
	Nickname string    `json:"nickname"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	Status   int8      `json:"status"`
	Roles    []model.Role `json:"roles"`
}

func (s *UserService) Create(req *CreateUserRequest) error {
	existing, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("用户名已存在")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   1,
	}

	if err := s.userRepo.Create(user); err != nil {
		return err
	}

	if req.RoleID > 0 {
		return s.userRepo.AssignRole(user.ID, req.RoleID)
	}

	return nil
}

func (s *UserService) Update(id uint, req *UpdateUserRequest) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	user.Nickname = req.Nickname
	user.Email = req.Email
	user.Phone = req.Phone
	user.Status = req.Status

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	if req.RoleID > 0 {
		if err := s.userRepo.RemoveRoles(id); err != nil {
			return err
		}
		return s.userRepo.AssignRole(id, req.RoleID)
	}

	return nil
}

func (s *UserService) Delete(id uint) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	if err := s.userRepo.RemoveRoles(id); err != nil {
		return err
	}

	return s.userRepo.Delete(id)
}

func (s *UserService) GetAll(page, pageSize int) ([]UserVO, int64, error) {
	users, total, err := s.userRepo.FindAll(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var result []UserVO
	for _, u := range users {
		roles, err := s.userRepo.FindRolesByUserID(u.ID)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, UserVO{
			ID:       u.ID,
			Username: u.Username,
			Nickname: u.Nickname,
			Email:    u.Email,
			Phone:    u.Phone,
			Status:   u.Status,
			Roles:    roles,
		})
	}

	return result, total, nil
}
