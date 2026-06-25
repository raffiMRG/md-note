package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"md-note/backend/internal/auth"
	"md-note/backend/internal/config"
	"md-note/backend/internal/handlers"
)

func New(
	cfg config.Config,
	cache *CORSCache,
	authHandler *handlers.AuthHandler,
	noteHandler *handlers.NoteHandler,
	tagHandler *handlers.TagHandler,
	corsHandler *handlers.CORSHandler,
) *gin.Engine {
	r := gin.Default()

	// Dynamic CORS — checks origin against cache on every request
	r.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" && cache.Allow(origin) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now()})
	})

	api := r.Group("/api")
	{
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)

		api.GET("/notes", noteHandler.List)
		api.GET("/notes/search", noteHandler.Search)
		api.GET("/notes/:id", noteHandler.Get)
		api.GET("/tags", tagHandler.List)

		// CORS origins list is readable by anyone (client may need to self-check)
		api.GET("/cors-origins", corsHandler.List)

		protected := api.Group("")
		protected.Use(auth.Middleware(cfg.JWTSecret))
		{
			protected.GET("/me", authHandler.Me)

			protected.POST("/notes", noteHandler.Create)
			protected.PUT("/notes/:id", noteHandler.Update)
			protected.DELETE("/notes/:id", noteHandler.Delete)

			protected.POST("/tags", tagHandler.Create)
			protected.PUT("/tags/:id", tagHandler.Update)
			protected.DELETE("/tags/:id", tagHandler.Delete)

			protected.POST("/cors-origins", corsHandler.Create)
			protected.DELETE("/cors-origins/:id", corsHandler.Delete)
		}
	}

	return r
}
