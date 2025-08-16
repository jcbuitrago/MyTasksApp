package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"task-manager-backend/models"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *models.Task) error {
	query := "INSERT INTO tasks (title, description, category_id, user_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err := r.db.QueryRow(query, task.Title, task.Description, task.CategoryID, task.UserID).Scan(&task.ID)
	if err != nil {
		log.Printf("Error creating task: %v", err)
		return err
	}
	return nil
}

func (r *TaskRepository) GetByID(id int) (*models.Task, error) {
	query := "SELECT id, title, description, category_id, user_id FROM tasks WHERE id = $1"
	task := &models.Task{}
	err := r.db.QueryRow(query, id).Scan(&task.ID, &task.Title, &task.Description, &task.CategoryID, &task.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("task not found")
		}
		log.Printf("Error retrieving task: %v", err)
		return nil, err
	}
	return task, nil
}

func (r *TaskRepository) Update(task *models.Task) error {
	query := "UPDATE tasks SET title = $1, description = $2, category_id = $3, user_id = $4 WHERE id = $5"
	_, err := r.db.Exec(query, task.Title, task.Description, task.CategoryID, task.UserID, task.ID)
	if err != nil {
		log.Printf("Error updating task: %v", err)
		return err
	}
	return nil
}

func (r *TaskRepository) Delete(id int) error {
	query := "DELETE FROM tasks WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting task: %v", err)
		return err
	}
	return nil
}

func (r *TaskRepository) GetAll() ([]models.Task, error) {
	query := "SELECT id, title, description, category_id, user_id FROM tasks"
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("Error retrieving tasks: %v", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.CategoryID, &task.UserID); err != nil {
			log.Printf("Error scanning task: %v", err)
			continue
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}