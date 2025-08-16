package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"task-manager-backend/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(category *models.Category) error {
	query := "INSERT INTO categories (name) VALUES ($1) RETURNING id"
	err := r.db.QueryRow(query, category.Name).Scan(&category.ID)
	if err != nil {
		log.Printf("Error creating category: %v", err)
		return err
	}
	return nil
}

func (r *CategoryRepository) GetByID(id int) (*models.Category, error) {
	category := &models.Category{}
	query := "SELECT id, name FROM categories WHERE id = $1"
	err := r.db.QueryRow(query, id).Scan(&category.ID, &category.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("category not found")
		}
		log.Printf("Error fetching category: %v", err)
		return nil, err
	}
	return category, nil
}

func (r *CategoryRepository) Update(category *models.Category) error {
	query := "UPDATE categories SET name = $1 WHERE id = $2"
	_, err := r.db.Exec(query, category.Name, category.ID)
	if err != nil {
		log.Printf("Error updating category: %v", err)
		return err
	}
	return nil
}

func (r *CategoryRepository) Delete(id int) error {
	query := "DELETE FROM categories WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting category: %v", err)
		return err
	}
	return nil
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	query := "SELECT id, name FROM categories"
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			log.Printf("Error scanning category: %v", err)
			continue
		}
		categories = append(categories, category)
	}
	return categories, nil
}