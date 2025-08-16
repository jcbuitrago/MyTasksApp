package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"task-manager-backend/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id"
	err := r.db.QueryRow(query, user.Name, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}
	return nil
}

func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	query := "SELECT id, name, email, password FROM users WHERE id = $1"
	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		log.Printf("Error getting user by ID: %v", err)
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	query := "UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4"
	_, err := r.db.Exec(query, user.Name, user.Email, user.Password, user.ID)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return err
	}
	return nil
}

func (r *UserRepository) DeleteUser(id int) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return err
	}
	return nil
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	query := "SELECT id, name, email, password FROM users"
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("Error getting all users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
			log.Printf("Error scanning user: %v", err)
			continue
		}
		users = append(users, user)
	}
	return users, nil
}