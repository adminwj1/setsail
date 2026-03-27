package service

import (
	"errors"
	"projecthub/internal/repository"
	"projecthub/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
	jwtUtil  *jwt.JWT
}

func NewAuthService(userRepo *repository.UserRepository, jwtUtil *jwt.JWT) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwtUtil:  jwtUtil,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string      `json:"token"`
	UserInfo  *UserInfoVO `json:"user_info"`
}

type UserInfoVO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户名或密码错误")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if user.Status != 1 {
		return nil, errors.New("账号已被禁用")
	}

	roles, err := s.userRepo.FindRolesByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	var roleID uint
	if len(roles) > 0 {
		roleID = roles[0].ID
	}

	token, err := s.jwtUtil.GenerateToken(user.ID, user.Username, roleID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		UserInfo: &UserInfoVO{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Email:    user.Email,
		},
	}, nil
}

func (s *AuthService) GetUserInfo(userID uint) (*UserInfoVO, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	return &UserInfoVO{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Email:    user.Email,
	}, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
