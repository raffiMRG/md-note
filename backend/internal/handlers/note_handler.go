package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"md-note/backend/internal/auth"
	"md-note/backend/internal/models"
	"md-note/backend/internal/repository"
)

type NoteHandler struct {
	notes *repository.NoteRepository
}

func NewNoteHandler(notes *repository.NoteRepository) *NoteHandler {
	return &NoteHandler{notes: notes}
}

type noteRequest struct {
	Title   string   `json:"title" binding:"required"`
	Content string   `json:"content" binding:"required"`
	TagIDs  []uint64 `json:"tag_ids"`
}

func pagination(c *gin.Context) (page, limit int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return
}

func (h *NoteHandler) List(c *gin.Context) {
	page, limit := pagination(c)
	tagSlug := c.Query("tag")

	notes, total, err := h.notes.List(tagSlug, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list notes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notes": notes, "total": total, "page": page, "limit": limit})
}

func (h *NoteHandler) Search(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query param 'q' is required"})
		return
	}
	page, limit := pagination(c)

	notes, total, err := h.notes.Search(q, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search notes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notes": notes, "total": total, "page": page, "limit": limit})
}

func (h *NoteHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note id"})
		return
	}

	note, err := h.notes.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"note": note})
}

func (h *NoteHandler) Create(c *gin.Context) {
	var req noteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint64(auth.ContextUserIDKey)
	note := models.Note{
		Title:     req.Title,
		Content:   req.Content,
		CreatedBy: &userID,
		UpdatedBy: &userID,
	}

	if err := h.notes.Create(&note, req.TagIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create note"})
		return
	}

	created, err := h.notes.FindByID(note.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "note created but failed to reload"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"note": created})
}

func (h *NoteHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note id"})
		return
	}

	var req noteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing, err := h.notes.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
		return
	}

	userID := c.GetUint64(auth.ContextUserIDKey)
	existing.Title = req.Title
	existing.Content = req.Content
	existing.UpdatedBy = &userID

	if err := h.notes.Update(existing, req.TagIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update note"})
		return
	}

	updated, err := h.notes.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "note updated but failed to reload"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"note": updated})
}

func (h *NoteHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note id"})
		return
	}

	if err := h.notes.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete note"})
		return
	}

	c.Status(http.StatusNoContent)
}
