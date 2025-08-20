// backend/internal/handlers/tasks.go
package handlers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type TasksHandler struct {
	DB *sql.DB
}

// Allowed statuses per guía
var allowedStatus = map[string]bool{
	"Sin Empezar": true,
	"Empezada":    true,
	"Finalizada":  true,
}

type taskCreateReq struct {
	Description       string  `json:"description" binding:"required"`
	TentativeDueDate  *string `json:"tentative_due_date"` // YYYY-MM-DD
	Status            *string `json:"status"`             // default "Sin Empezar"
	CategoryID        *int64  `json:"category_id"`        // opcional
}

type taskUpdateReq struct {
	Description       string  `json:"description" binding:"required"`
	TentativeDueDate  *string `json:"tentative_due_date"`
	Status            *string `json:"status"`
	CategoryID        *int64  `json:"category_id"`
}

type taskResp struct {
	ID                int64     `json:"id"`
	Description       string    `json:"description"`
	CreatedAt         time.Time `json:"created_at"`
	TentativeDueDate  *string   `json:"tentative_due_date,omitempty"`
	Status            string    `json:"status"`
	CategoryID        *int64    `json:"category_id,omitempty"`
	UserID            int64     `json:"user_id"`
}

func validateStatus(s *string) (string, error) {
	def := "Sin Empezar"
	if s == nil || strings.TrimSpace(*s) == "" {
		return def, nil
	}
	val := strings.TrimSpace(*s)
	if !allowedStatus[val] {
		return "", errors.New("estado inválido (permitidos: 'Sin Empezar', 'Empezada', 'Finalizada')")
	}
	return val, nil
}

// POST /tareas  (protegido)
func (h *TasksHandler) Create(c *gin.Context) {
	var req taskCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json inválido"})
		return
	}
	uidv, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return
	}
	userID := uidv.(int64)

	status, err := validateStatus(req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var id int64
	var createdAt time.Time
	q := `INSERT INTO tasks (description, tentative_due_date, status, category_id, user_id)
	      VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	err = h.DB.QueryRowContext(ctx, q,
		strings.TrimSpace(req.Description),
		req.TentativeDueDate, // nil -> NULL
		status,
		req.CategoryID, // nil -> NULL
		userID,
	).Scan(&id, &createdAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	resp := taskResp{
		ID:               id,
		Description:      strings.TrimSpace(req.Description),
		CreatedAt:        createdAt,
		TentativeDueDate: req.TentativeDueDate,
		Status:           status,
		CategoryID:       req.CategoryID,
		UserID:           userID,
	}
	c.JSON(http.StatusCreated, resp)
}

// PUT /tareas/:id  (protegido)
func (h *TasksHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	taskID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}
	var req taskUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json inválido"})
		return
	}
	uidv, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return
	}
	userID := uidv.(int64)

	status, err := validateStatus(req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// Verificar pertenencia
	var exists bool
	err = h.DB.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM tasks WHERE id=$1 AND user_id=$2)`,
		taskID, userID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "no existe o no te pertenece"})
		return
	}

	_, err = h.DB.ExecContext(ctx, `UPDATE tasks
		SET description=$1, tentative_due_date=$2, status=$3, category_id=$4
		WHERE id=$5 AND user_id=$6`,
		strings.TrimSpace(req.Description),
		req.TentativeDueDate,
		status,
		req.CategoryID,
		taskID,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.Status(http.StatusNoContent)
}

// DELETE /tareas/:id  (protegido)
func (h *TasksHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	taskID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}
	uidv, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return
	}
	userID := uidv.(int64)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	res, err := h.DB.ExecContext(ctx,
		`DELETE FROM tasks WHERE id=$1 AND user_id=$2`, taskID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no existe o no te pertenece"})
		return
	}
	c.Status(http.StatusNoContent)
}

// GET /tareas/usuario?categoria_id=&estado=  (protegido)
func (h *TasksHandler) ListByUser(c *gin.Context) {
	uidv, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return
	}
	userID := uidv.(int64)

	q := `SELECT id, description, created_at, tentative_due_date, status, category_id, user_id
	      FROM tasks WHERE user_id=$1`
	args := []interface{}{userID}
	i := 2

	if cat := strings.TrimSpace(c.Query("categoria_id")); cat != "" {
		if _, err := strconv.ParseInt(cat, 10, 64); err == nil {
			q += " AND category_id=$" + strconv.Itoa(i)
			args = append(args, cat)
			i++
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "categoria_id inválido"})
			return
		}
	}
	if st := strings.TrimSpace(c.Query("estado")); st != "" {
		if !allowedStatus[st] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "estado inválido"})
			return
		}
		q += " AND status=$" + strconv.Itoa(i)
		args = append(args, st)
		i++
	}
	q += " ORDER BY created_at DESC"

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	rows, err := h.DB.QueryContext(ctx, q, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	defer rows.Close()

	out := make([]taskResp, 0)
	for rows.Next() {
		var r taskResp
		var due *time.Time
		var catID sql.NullInt64
		if err := rows.Scan(&r.ID, &r.Description, &r.CreatedAt, &due, &r.Status, &catID, &r.UserID); err != nil {
			continue
		}
		if due != nil {
			s := due.Format("2006-01-02")
			r.TentativeDueDate = &s
		}
		if catID.Valid {
			v := catID.Int64
			r.CategoryID = &v
		}
		out = append(out, r)
	}
	c.JSON(http.StatusOK, out)
}

// GET /tareas/:id  (protegido)
func (h *TasksHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	taskID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}
	uidv, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return
	}
	userID := uidv.(int64)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var r taskResp
	var due *time.Time
	var catID sql.NullInt64
	err = h.DB.QueryRowContext(ctx, `SELECT id, description, created_at, tentative_due_date, status, category_id, user_id
	                                 FROM tasks WHERE id=$1 AND user_id=$2`,
		taskID, userID).Scan(&r.ID, &r.Description, &r.CreatedAt, &due, &r.Status, &catID, &r.UserID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "no existe o no te pertenece"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if due != nil {
		s := due.Format("2006-01-02")
		r.TentativeDueDate = &s
	}
	if catID.Valid {
		v := catID.Int64
		r.CategoryID = &v
	}
	c.JSON(http.StatusOK, r)
}
