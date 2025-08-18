package handlers

import (
	"net/http"
	"strconv"
	"time"

	"backend/models"
    "mytasks/internal/models"
    "mytasks/internal/repository"
)

type TaskHandler struct{ repo *repository.TaskRepository }

func NewTaskHandler(r *repository.TaskRepository) *TaskHandler { return &TaskHandler{repo: r} }

func userIDFromCtx(c *gin.Context) (int64, bool) {
	v, ok := c.Get("user_id")
	if !ok { return 0, false }
	switch x := v.(type) {
	case int64:   return x, true
	case int:     return int64(x), true
	case float64: return int64(x), true
	default:      return 0, false
	}
}

type createTaskPayload struct {
	CategoryID  *int64  `json:"category_id"`
	ParentID    *int64  `json:"parent_id"`
	Title       string  `json:"title" binding:"required,min=1,max=100"`
	Description string  `json:"description"`
	Status      string  `json:"status" binding:"omitempty,oneof=Backlog In_Progress Done"`
	Priority    *int    `json:"priority"`
	DueDate     *string `json:"due_date"` // YYYY-MM-DD
}

func (h *TaskHandler) Create(c *gin.Context) {
	uid, ok := userIDFromCtx(c)
	if !ok { c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"}); return }

	var p createTaskPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}

	var due *time.Time
	if p.DueDate != nil && *p.DueDate != "" {
		t, err := time.Parse("2006-01-02", *p.DueDate)
		if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid due_date YYYY-MM-DD"}); return }
		due = &t
	}

	status := p.Status
	if status == "" { status = "Backlog" }
	priority := 0
	if p.Priority != nil { priority = *p.Priority }

	t := &models.Task{
		UserID: uid, CategoryID: p.CategoryID, ParentID: p.ParentID,
		Title: p.Title, Description: p.Description, Status: status,
		Priority: priority, DueDate: due,
	}
	if err := h.repo.Create(c.Request.Context(), t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusCreated, t)
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	uid, ok := userIDFromCtx(c)
	if !ok { c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"}); return }

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}); return }

	t, err := h.repo.GetByID(c.Request.Context(), id, uid)
	if err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "task not found"}); return }
	c.JSON(http.StatusOK, t)
}

func (h *TaskHandler) List(c *gin.Context) {
	uid, ok := userIDFromCtx(c)
	if !ok { c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"}); return }

	var f repository.ListFilters
	f.Status = c.Query("status")
	if v := c.Query("category_id"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil { f.CategoryID = &id }
	}
	f.Search = c.Query("q")
	if v := c.Query("due_before"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil { f.DueBefore = &t }
	}
	if v := c.Query("due_after"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil { f.DueAfter = &t }
	}
	if v := c.Query("page"); v != "" {
		if i, err := strconv.Atoi(v); err == nil { f.Page = i }
	}
	if v := c.Query("page_size"); v != "" {
		if i, err := strconv.Atoi(v); err == nil { f.PageSize = i }
	}

	tasks, err := h.repo.List(c.Request.Context(), uid, f)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, tasks)
}

type updateTaskPayload struct {
	CategoryID  *int64  `json:"category_id"`
	ParentID    *int64  `json:"parent_id"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status" binding:"omitempty,oneof=Backlog In_Progress Done"`
	Priority    *int    `json:"priority"`
	DueDate     *string `json:"due_date"`
}

func (h *TaskHandler) Update(c *gin.Context) {
	uid, ok := userIDFromCtx(c)
	if !ok { c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"}); return }

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}); return }

	existing, err := h.repo.GetByID(c.Request.Context(), id, uid)
	if err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "task not found"}); return }

	var p updateTaskPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}

	if p.CategoryID != nil { existing.CategoryID = p.CategoryID }
	if p.ParentID   != nil { existing.ParentID   = p.ParentID }
	if p.Title      != nil { existing.Title      = *p.Title }
	if p.Description!= nil { existing.Description= *p.Description }
	if p.Status     != nil { existing.Status     = *p.Status }
	if p.Priority   != nil { existing.Priority   = *p.Priority }
	if p.DueDate    != nil {
		if *p.DueDate == "" {
			existing.DueDate = nil
		} else if t, err := time.Parse("2006-01-02", *p.DueDate); err == nil {
			existing.DueDate = &t
		} else { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid due_date"}); return }
	}

	if err := h.repo.Update(c.Request.Context(), existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, existing)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	uid, ok := userIDFromCtx(c)
	if !ok { c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"}); return }

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}); return }

	if err := h.repo.Delete(c.Request.Context(), id, uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.Status(http.StatusNoContent)
}

// handlers/tasks.go
func (h *TaskHandler) ListByUser(c *gin.Context) {
	uid, ok := userIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var f repository.ListFilters

	// estado => status
	f.Status = c.Query("estado")

	// categoria_id => CategoryID
	if v := c.Query("categoria_id"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.CategoryID = &id
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "categoria_id inválido"})
			return
		}
	}

	// (opcional) soporte a búsqueda por texto q=
	f.Search = c.Query("q")

	// (opcional) paginación page & page_size
	if v := c.Query("page"); v != "" {
		if i, err := strconv.Atoi(v); err == nil { f.Page = i }
	}
	if v := c.Query("page_size"); v != "" {
		if i, err := strconv.Atoi(v); err == nil { f.PageSize = i }
	}

	tasks, err := h.repo.List(c.Request.Context(), uid, f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}
