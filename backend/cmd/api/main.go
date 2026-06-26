package main

import (
	"log"

	"md-note/backend/internal/auth"
	"md-note/backend/internal/config"
	"md-note/backend/internal/db"
	"md-note/backend/internal/handlers"
	"md-note/backend/internal/models"
	"md-note/backend/internal/repository"
	"md-note/backend/internal/router"
	"md-note/backend/migrations"
)

func main() {
	cfg := config.Load()

	gdb, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		log.Fatalf("get sql.DB failed: %v", err)
	}

	if err := db.RunMigrations(sqlDB, migrations.FS); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}
	log.Println("migrations applied")

	userRepo := repository.NewUserRepository(gdb)
	noteRepo := repository.NewNoteRepository(gdb)
	tagRepo := repository.NewTagRepository(gdb)
	corsRepo := repository.NewCORSRepository(gdb)

	// Build CORS cache: static origins from env + dynamic ones from DB
	corsCache := router.NewCORSCache(cfg.CORSOrigins)
	dbOrigins, err := corsRepo.AllOrigins()
	if err != nil {
		log.Printf("warn: could not load cors origins from db: %v", err)
	} else {
		corsCache.LoadDynamic(dbOrigins)
	}
	log.Printf("CORS allowed origins (static): %v", cfg.CORSOrigins)

	seedAdmin(userRepo)

	authHandler := handlers.NewAuthHandler(userRepo, cfg.JWTSecret)
	noteHandler := handlers.NewNoteHandler(noteRepo)
	tagHandler := handlers.NewTagHandler(tagRepo)
	corsHandler := handlers.NewCORSHandler(corsRepo, corsCache.Add, corsCache.Remove)
	backupHandler := handlers.NewBackupHandler(gdb)
	userHandler := handlers.NewUserHandler(userRepo)

	r := router.New(cfg, corsCache, authHandler, noteHandler, tagHandler, corsHandler, backupHandler, userHandler)

	log.Printf("listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func seedAdmin(repo *repository.UserRepository) {
	if _, err := repo.FindByEmail("mrg@raffimrg.my.id"); err == nil {
		return
	}
	hash, err := auth.HashPassword("123123123")
	if err != nil {
		log.Printf("warn: failed to hash admin password: %v", err)
		return
	}
	if err := repo.Create(&models.User{
		Username:     "mrg",
		Email:        "mrg@raffimrg.my.id",
		PasswordHash: hash,
		Role:         "admin",
	}); err != nil {
		log.Printf("warn: failed to seed admin: %v", err)
		return
	}
	log.Println("admin seed: mrg@raffimrg.my.id created")
}
