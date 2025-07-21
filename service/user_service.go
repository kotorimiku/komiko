package service

import (
	"fmt"
	"komiko/config"
	"komiko/dto"
	"komiko/model"
	"komiko/repo"
	"komiko/utils"
)

type UserService struct {
	baseService[model.User, *repo.UserRepo]
	repos *repo.Repo
}

func NewUserService(repos *repo.Repo) *UserService {
	return &UserService{baseService: baseService[model.User, *repo.UserRepo]{repo: repos.UserRepo}, repos: repos}
}

func (s *UserService) Create(user *model.User) error {
	exists, err := s.repo.Exists(&model.User{Username: user.Username})
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("用户名已存在")
	}
	if !utils.IsValidPassword(user.Password) {
		return fmt.Errorf("密码不符合要求")
	}

	user.Password, err = utils.EncryptPassword(user.Password)
	if err != nil {
		return err
	}
	err = s.repo.Create(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) Register(user *model.User) error {
	count, err := s.repo.Count()
	if err != nil {
		return err
	}
	if !config.GetConfig().AllowRegister && count > 0 {
		return fmt.Errorf("not allowed to register")
	}
	if count == 0 {
		user.Role = "admin"
	}
	return s.Create(user)
}

func (s *UserService) Login(user *model.User) (string, error) {
	dbUser, err := s.repo.GetByUsername(user.Username)
	if err != nil || dbUser == nil {
		return "", fmt.Errorf("invalid username or password")
	}
	if utils.PasswordMatches(user.Password, dbUser.Password) {
		return utils.GenerateJWT(fmt.Sprintf("%d", dbUser.ID))
	}
	return "", fmt.Errorf("invalid username or password")
}

func (s *UserService) CreateUser(currentUserID uint, user *model.User) error {
	if currentUserID == 1 {
		err := s.Create(user)
		return err
	}
	currentUser, err := s.repo.GetByID(currentUserID)
	if err != nil {
		return err
	}
	if currentUser.Role != "admin" {
		return fmt.Errorf("无权限创建用户")
	}

	err = s.Create(user)
	return err
}

func (s *UserService) DeleteUserByID(currentUserID uint, id uint) error {
	if currentUserID == 1 {
		return fmt.Errorf("禁止删除root用户")
	}
	if id == currentUserID || currentUserID == 1 {
		return s.repo.DeleteByID(id)
	}

	currentUser, err := s.repo.GetByID(currentUserID)
	if err != nil {
		return err
	}
	if currentUser.Role != "admin" {
		return fmt.Errorf("无权限删除用户")
	}

	return s.repo.DeleteByID(id)
}

func (s *UserService) UpdateUser(currentUserID uint, user *model.User) error {
	if !utils.IsValidPassword(user.Password) && user.Password != "" {
		return fmt.Errorf("密码不符合要求")
	}
	if currentUserID == 1 {
		return s.repo.Update(user)
	}

	if user.ID == 1 {
		return fmt.Errorf("禁止更新root用户")
	}

	currentUser, err := s.repo.GetByID(currentUserID)
	if err != nil {
		return err
	}

	if user.ID != currentUserID && currentUser.Role != "admin" {
		return fmt.Errorf("无权限更新用户")
	}

	if currentUser.Role != "admin" {
		user.Role = ""
	}

	return s.repo.Update(user)
}

func (s *UserService) GetAllUserDto() ([]*dto.User, error) {
	return s.repo.GetAllDto()
}

func (s *UserService) GetUserDtoByID(id uint) (*dto.User, error) {
	return s.repo.GetDtoByID(id)
}

func (s *UserService) AllowedRegister() bool {
	count, err := s.repo.Count()
	if err != nil {
		return false
	}
	return config.GetConfig().AllowRegister || count == 0
}
