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
	backupHandler *handlers.BackupHandler,
	userHandler *handlers.UserHandler,
	uploadHandler *handlers.UploadHandler,
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

		public := api.Group("", auth.OptionalMiddleware(cfg.JWTSecret))
		public.GET("/notes", noteHandler.List)
		public.GET("/notes/search", noteHandler.Search)
		public.GET("/notes/:id", noteHandler.Get)
		api.GET("/tags", tagHandler.List)

		uploads := api.Group("/uploads", func(c *gin.Context) {
			c.Header("X-Content-Type-Options", "nosniff")
		})
		uploads.Static("/", cfg.UploadDir)

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

			protected.POST("/uploads", uploadHandler.Upload)
		}

		admin := api.Group("")
		admin.Use(auth.Middleware(cfg.JWTSecret), auth.RequireAdmin())
		{
			admin.POST("/cors-origins", corsHandler.Create)
			admin.DELETE("/cors-origins/:id", corsHandler.Delete)

			admin.GET("/backup", backupHandler.Export)
			admin.POST("/restore", backupHandler.Import)

			admin.GET("/admin/users", userHandler.List)
			admin.PUT("/admin/users/:id/role", userHandler.UpdateRole)
			admin.DELETE("/admin/users/:id", userHandler.Delete)
		}
	}

	return r
}
