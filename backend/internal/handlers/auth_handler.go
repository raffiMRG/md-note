package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"md-note/backend/internal/auth"
	"md-note/backend/internal/models"
	"md-note/backend/internal/repository"
)

type AuthHandler struct {
	users                *repository.UserRepository
	jwtSecret            string
	jwtExpiresIn         time.Duration
	jwtRememberExpiresIn time.Duration
}

func NewAuthHandler(users *repository.UserRepository, jwtSecret string, jwtExpiresIn, jwtRememberExpiresIn time.Duration) *AuthHandler {
	return &AuthHandler{
		users:                users,
		jwtSecret:            jwtSecret,
		jwtExpiresIn:         jwtExpiresIn,
		jwtRememberExpiresIn: jwtRememberExpiresIn,
	}
}

type registerRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RememberMe bool   `json:"remember_me"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process password"})
		return
	}

	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hash,
		Role:         "penulis",
	}
	if err := h.users.Create(&user); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username or email already taken"})
		return
	}

	token, err := auth.GenerateToken(h.jwtSecret, user.ID, user.Username, user.Role, h.jwtExpiresIn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token, "user": user})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.users.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to look up user"})
		return
	}

	if !auth.CheckPassword(user.PasswordHash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	ttl := h.jwtExpiresIn
	if req.RememberMe {
		ttl = h.jwtRememberExpiresIn
	}

	token, err := auth.GenerateToken(h.jwtSecret, user.ID, user.Username, user.Role, ttl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetUint64(auth.ContextUserIDKey)
	user, err := h.users.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}
