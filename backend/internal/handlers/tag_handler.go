package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"md-note/backend/internal/repository"
)

type TagHandler struct {
	tags *repository.TagRepository
}

func NewTagHandler(tags *repository.TagRepository) *TagHandler {
	return &TagHandler{tags: tags}
}

type tagRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

func (h *TagHandler) List(c *gin.Context) {
	tags, err := h.tags.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tags"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tags": tags})
}

func (h *TagHandler) Create(c *gin.Context) {
	var req tagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag, err := h.tags.Create(req.Name)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "tag already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"tag": tag})
}

func (h *TagHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag id"})
		return
	}

	var req tagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag, err := h.tags.Update(id, req.Name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "tag name already exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tag": tag})
}

func (h *TagHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag id"})
		return
	}

	if err := h.tags.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete tag"})
		return
	}

	c.Status(http.StatusNoContent)
}
