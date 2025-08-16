package services

import (
	"task-manager-backend/models"
	"task-manager-backend/repositories"
)

type TaskService struct {
	repo repositories.TaskRepository
}

func NewTaskService(repo repositories.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(task *models.Task) error {
	return s.repo.Create(task)
}

func (s *TaskService) GetTaskByID(id int) (*models.Task, error) {
	return s.repo.FindByID(id)
}

func (s *TaskService) UpdateTask(task *models.Task) error {
	return s.repo.Update(task)
}

func (s *TaskService) DeleteTask(id int) error {
	return s.repo.Delete(id)
}

func (s *TaskService) GetAllTasks() ([]models.Task, error) {
	return s.repo.FindAll()
}