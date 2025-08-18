package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

    "mytasks/internal/models"
)

type TaskRepository struct{ db *sql.DB }

func NewTaskRepository(db *sql.DB) *TaskRepository { return &TaskRepository{db: db} }

func (r *TaskRepository) Create(ctx context.Context, t *models.Task) error {
	q := `INSERT INTO tasks (user_id, category_id, parent_id, title, description, status, priority, due_date)
	      VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	      RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, q,
		t.UserID, t.CategoryID, t.ParentID, t.Title, t.Description, t.Status, t.Priority, t.DueDate,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *TaskRepository) GetByID(ctx context.Context, id int64, userID int64) (*models.Task, error) {
	q := `SELECT id, user_id, category_id, parent_id, title, description, status, priority, due_date, created_at, updated_at
	      FROM tasks
	      WHERE id=$1 AND user_id=$2 AND deleted_at IS NULL`
	var t models.Task
	err := r.db.QueryRowContext(ctx, q, id, userID).Scan(
		&t.ID, &t.UserID, &t.CategoryID, &t.ParentID, &t.Title, &t.Description,
		&t.Status, &t.Priority, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil { return nil, err }
	return &t, nil
}

type ListFilters struct {
	Status     string
	CategoryID *int64
	Search     string
	DueBefore  *time.Time
	DueAfter   *time.Time
	Page       int
	PageSize   int
}

func (r *TaskRepository) List(ctx context.Context, userID int64, f ListFilters) ([]models.Task, error) {
	var args []any
	var where []string
	args = append(args, userID)
	where = append(where, "user_id=$1", "deleted_at IS NULL")
	i := 2

	if f.Status != "" {
		where = append(where, fmt.Sprintf("status=$%d", i))
		args = append(args, f.Status); i++
	}
	if f.CategoryID != nil {
		where = append(where, fmt.Sprintf("category_id=$%d", i))
		args = append(args, *f.CategoryID); i++
	}
	if f.Search != "" {
		where = append(where, fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d)", i, i))
		args = append(args, "%"+f.Search+"%"); i++
	}
	if f.DueBefore != nil {
		where = append(where, fmt.Sprintf("due_date <= $%d", i))
		args = append(args, *f.DueBefore); i++
	}
	if f.DueAfter != nil {
		where = append(where, fmt.Sprintf("due_date >= $%d", i))
		args = append(args, *f.DueAfter); i++
	}

	if f.Page <= 0 { f.Page = 1 }
	if f.PageSize <= 0 || f.PageSize > 100 { f.PageSize = 20 }
	offset := (f.Page - 1) * f.PageSize

	q := fmt.Sprintf(`
	  SELECT id, user_id, category_id, parent_id, title, description, status, priority, due_date, created_at, updated_at
	  FROM tasks
	  WHERE %s
	  ORDER BY created_at DESC
	  LIMIT %d OFFSET %d`, strings.Join(where, " AND "), f.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil { return nil, err }
	defer rows.Close()

	var out []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(
			&t.ID, &t.UserID, &t.CategoryID, &t.ParentID, &t.Title, &t.Description,
			&t.Status, &t.Priority, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
		); err != nil { return nil, err }
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *TaskRepository) Update(ctx context.Context, t *models.Task) error {
	q := `UPDATE tasks
	      SET category_id=$1, parent_id=$2, title=$3, description=$4,
	          status=$5, priority=$6, due_date=$7, updated_at=now()
	      WHERE id=$8 AND user_id=$9 AND deleted_at IS NULL
	      RETURNING created_at, updated_at`
	return r.db.QueryRowContext(ctx, q,
		t.CategoryID, t.ParentID, t.Title, t.Description, t.Status, t.Priority, t.DueDate, t.ID, t.UserID,
	).Scan(&t.CreatedAt, &t.UpdatedAt)
}

func (r *TaskRepository) Delete(ctx context.Context, id int64, userID int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE tasks SET deleted_at=now() WHERE id=$1 AND user_id=$2 AND deleted_at IS NULL`,
		id, userID,
	)
	return err
}
