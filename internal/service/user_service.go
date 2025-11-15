// Package service 包含了应用的业务逻辑
package service

import (
	"errors"
	"questflow/internal/model"
	"questflow/internal/repository"
	"questflow/pkg/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// CustomClaims 定义了我们自己的 JWT 声明
type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     uint8  `json:"role"`
	jwt.RegisteredClaims
}

// UserService 定义了用户服务的接口
type UserService interface {
	Register(username, password string, email *string) (*model.User, error)
	Login(username, password string) (string, *model.User, error) // 返回 token 和用户信息
	UpdateProfile(userID uint, newUsername, newPassword string) (*model.User, error)
}

// userServiceImpl 是 UserService 的实现
type userServiceImpl struct {
	userRepo repository.UserRepository
}

// NewUserService 创建一个新的 UserService 实例
func NewUserService(repo repository.UserRepository) UserService {
	return &userServiceImpl{userRepo: repo}
}

// Register 处理用户注册的业务逻辑
func (s *userServiceImpl) Register(username, password string, email *string) (*model.User, error) {
	// 1. 检查用户名是否已存在
	existingUser, err := s.userRepo.FindByUsername(username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// 查询数据库时发生其他错误
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// 2. 创建新用户实例并加密密码
	newUser := &model.User{
		Username: username,
		Email:    email,
	}
	if err := newUser.SetPassword(password); err != nil {
		return nil, errors.New("failed to hash password")
	}

	// 3. 在数据库中创建用户
	if err := s.userRepo.Create(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

// Login 处理用户登录逻辑
func (s *userServiceImpl) Login(username, password string) (string, *model.User, error) {
	// 1. 根据用户名查找用户
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("invalid username or password")
		}
		return "", nil, err // 其他数据库错误
	}

	// 2. 验证密码
	if !user.CheckPassword(password) {
		return "", nil, errors.New("invalid username or password")
	}

	// 3. 密码验证成功，生成 JWT
	token, err := s.generateToken(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// UpdateProfile 处理更新用户信息的业务逻辑
func (s *userServiceImpl) UpdateProfile(userID uint, newUsername, newPassword string) (*model.User, error) {
	// 1. 查找当前用户
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 2. 如果提供了新用户名，则更新
	if newUsername != "" && newUsername != user.Username {
		// 检查新用户名是否已被占用
		existingUser, err := s.userRepo.FindByUsername(newUsername)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err // 数据库查询错误
		}
		if existingUser != nil {
			return nil, errors.New("new username is already taken")
		}
		user.Username = newUsername
	}

	// 3. 如果提供了新密码，则更新
	if newPassword != "" {
		if err := user.SetPassword(newPassword); err != nil {
			return nil, errors.New("failed to update password")
		}
	}

	// 4. 保存到数据库
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// generateToken 为指定用户生成 JWT
func (s *userServiceImpl) generateToken(user *model.User) (string, error) {
	expireTime := time.Now().Add(time.Hour * time.Duration(config.Cfg.JWT.ExpireHours))

	claims := CustomClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    config.Cfg.JWT.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 使用 HS256 签名方法
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用配置的 secret 签名并获取完整的编码后的字符串 token
	signedToken, err := token.SignedString([]byte(config.Cfg.JWT.Secret))
	if err != nil {
		return "", errors.New("failed to sign token")
	}

	return signedToken, nil
}
