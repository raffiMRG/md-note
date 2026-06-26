package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"md-note/backend/internal/models"
)

type BackupHandler struct {
	db *gorm.DB
}

func NewBackupHandler(db *gorm.DB) *BackupHandler {
	return &BackupHandler{db: db}
}

type RawUser struct {
	ID           uint64    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RawNote struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedBy *uint64   `json:"created_by"`
	UpdatedBy *uint64   `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RawNoteTag struct {
	NoteID uint64 `json:"note_id"`
	TagID  uint64 `json:"tag_id"`
}

type BackupData struct {
	Version     int                 `json:"version"`
	ExportedAt  time.Time           `json:"exported_at"`
	Users       []RawUser           `json:"users"`
	Tags        []models.Tag        `json:"tags"`
	Notes       []RawNote           `json:"notes"`
	NoteTags    []RawNoteTag        `json:"note_tags"`
	CORSOrigins []models.CORSOrigin `json:"cors_origins"`
}

func (h *BackupHandler) Export(c *gin.Context) {
	var data BackupData
	data.Version = 1
	data.ExportedAt = time.Now().UTC()

	queries := []struct {
		sql    string
		dest   interface{}
		errMsg string
	}{
		{"SELECT id, username, email, password_hash, created_at, updated_at FROM users", &data.Users, "gagal mengambil data users"},
		{"SELECT id, name, slug FROM tags", &data.Tags, "gagal mengambil data tags"},
		{"SELECT id, title, content, created_by, updated_by, created_at, updated_at FROM notes", &data.Notes, "gagal mengambil data notes"},
		{"SELECT note_id, tag_id FROM note_tags", &data.NoteTags, "gagal mengambil data note_tags"},
		{"SELECT id, origin, created_at FROM cors_origins", &data.CORSOrigins, "gagal mengambil data cors_origins"},
	}

	for _, q := range queries {
		if err := h.db.Raw(q.sql).Scan(q.dest).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": q.errMsg})
			return
		}
	}

	filename := fmt.Sprintf("md-note-backup-%s.json", time.Now().Format("2006-01-02"))
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.JSON(http.StatusOK, data)
}

func (h *BackupHandler) Import(c *gin.Context) {
	var data BackupData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format backup tidak valid"})
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		for _, u := range data.Users {
			if err := tx.Exec(
				"INSERT IGNORE INTO users (id, username, email, password_hash, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
				u.ID, u.Username, u.Email, u.PasswordHash, u.CreatedAt, u.UpdatedAt,
			).Error; err != nil {
				return err
			}
		}
		for _, t := range data.Tags {
			if err := tx.Exec(
				"INSERT IGNORE INTO tags (id, name, slug) VALUES (?, ?, ?)",
				t.ID, t.Name, t.Slug,
			).Error; err != nil {
				return err
			}
		}
		for _, n := range data.Notes {
			if err := tx.Exec(
				"INSERT IGNORE INTO notes (id, title, content, created_by, updated_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
				n.ID, n.Title, n.Content, n.CreatedBy, n.UpdatedBy, n.CreatedAt, n.UpdatedAt,
			).Error; err != nil {
				return err
			}
		}
		for _, nt := range data.NoteTags {
			if err := tx.Exec(
				"INSERT IGNORE INTO note_tags (note_id, tag_id) VALUES (?, ?)",
				nt.NoteID, nt.TagID,
			).Error; err != nil {
				return err
			}
		}
		for _, co := range data.CORSOrigins {
			if err := tx.Exec(
				"INSERT IGNORE INTO cors_origins (id, origin, created_at) VALUES (?, ?, ?)",
				co.ID, co.Origin, co.CreatedAt,
			).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengimpor backup: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "restore berhasil"})
}
