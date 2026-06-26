package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"md-note/backend/internal/auth"
	"md-note/backend/internal/repository"
)

type UserHandler struct {
	users *repository.UserRepository
}

func NewUserHandler(users *repository.UserRepository) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) List(c *gin.Context) {
	users, err := h.users.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

type updateRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=admin penulis"`
}

func (h *UserHandler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	var req updateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role harus 'admin' atau 'penulis'"})
		return
	}
	if err := h.users.UpdateRole(id, req.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update role"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "role updated"})
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	if id == c.GetUint64(auth.ContextUserIDKey) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tidak bisa menghapus akun sendiri"})
		return
	}
	if err := h.users.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}
	c.Status(http.StatusNoContent)
}
