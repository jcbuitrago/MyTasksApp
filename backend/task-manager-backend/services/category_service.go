package services

import (
	"task-manager-backend/models"
	"task-manager-backend/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) CreateCategory(category *models.Category) error {
	return s.repo.Create(category)
}

func (s *CategoryService) GetCategoryByID(id int) (*models.Category, error) {
	return s.repo.FindByID(id)
}

func (s *CategoryService) UpdateCategory(category *models.Category) error {
	return s.repo.Update(category)
}

func (s *CategoryService) DeleteCategory(id int) error {
	return s.repo.Delete(id)
}

func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	return s.repo.FindAll()
}