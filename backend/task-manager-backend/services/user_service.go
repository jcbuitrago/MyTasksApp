package services

import (
	"errors"
	"task-manager-backend/models"
	"task-manager-backend/repositories"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	return s.userRepo.Create(user)
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) UpdateUser(user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	return s.userRepo.Update(user)
}

func (s *UserService) DeleteUser(id int) error {
	return s.userRepo.Delete(id)
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.FindAll()
}