package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"md-note/backend/internal/repository"
)

type CORSHandler struct {
	repo   *repository.CORSRepository
	onAdd  func(string)
	onDel  func(string)
}

func NewCORSHandler(repo *repository.CORSRepository, onAdd, onDel func(string)) *CORSHandler {
	return &CORSHandler{repo: repo, onAdd: onAdd, onDel: onDel}
}

func (h *CORSHandler) List(c *gin.Context) {
	rows, err := h.repo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"origins": rows})
}

func (h *CORSHandler) Create(c *gin.Context) {
	var req struct {
		Origin string `json:"origin" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	origin := strings.TrimRight(strings.TrimSpace(req.Origin), "/")
	row, err := h.repo.Create(origin)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "origin already exists"})
		return
	}
	h.onAdd(origin)
	c.JSON(http.StatusCreated, gin.H{"origin": row})
}

func (h *CORSHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	// fetch origin value before deleting so we can evict from cache
	rows, _ := h.repo.List()
	var target string
	for _, r := range rows {
		if r.ID == id {
			target = r.Origin
			break
		}
	}
	if err := h.repo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if target != "" {
		h.onDel(target)
	}
	c.Status(http.StatusNoContent)
}
