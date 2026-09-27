package handlers

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"md-note/backend/internal/models"
)

type BackupHandler struct {
	db        *gorm.DB
	uploadDir string
}

func NewBackupHandler(db *gorm.DB, uploadDir string) *BackupHandler {
	return &BackupHandler{db: db, uploadDir: uploadDir}
}

// nama file hasil UploadHandler: 32 hex + ekstensi gambar
var uploadNamePattern = regexp.MustCompile(`^[0-9a-f]{32}\.(jpg|png|gif|webp)$`)

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
	IsPrivate bool      `json:"is_private"`
	CreatedBy *uint64   `json:"created_by"`
	UpdatedBy *uint64   `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RawNoteTag struct {
	NoteID uint64 `json:"note_id"`
	TagID  uint64 `json:"tag_id"`
}

// Data di-encode base64 otomatis oleh encoding/json
type RawUpload struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}

type BackupData struct {
	Version     int                 `json:"version"`
	ExportedAt  time.Time           `json:"exported_at"`
	Users       []RawUser           `json:"users"`
	Tags        []models.Tag        `json:"tags"`
	Notes       []RawNote           `json:"notes"`
	NoteTags    []RawNoteTag        `json:"note_tags"`
	CORSOrigins []models.CORSOrigin `json:"cors_origins"`
	Uploads     []RawUpload         `json:"uploads"`
}

func (h *BackupHandler) Export(c *gin.Context) {
	var data BackupData
	data.Version = 2
	data.ExportedAt = time.Now().UTC()

	queries := []struct {
		sql    string
		dest   interface{}
		errMsg string
	}{
		{"SELECT id, username, email, password_hash, created_at, updated_at FROM users", &data.Users, "gagal mengambil data users"},
		{"SELECT id, name, slug FROM tags", &data.Tags, "gagal mengambil data tags"},
		{"SELECT id, title, content, is_private, created_by, updated_by, created_at, updated_at FROM notes", &data.Notes, "gagal mengambil data notes"},
		{"SELECT note_id, tag_id FROM note_tags", &data.NoteTags, "gagal mengambil data note_tags"},
		{"SELECT id, origin, created_at FROM cors_origins", &data.CORSOrigins, "gagal mengambil data cors_origins"},
	}

	for _, q := range queries {
		if err := h.db.Raw(q.sql).Scan(q.dest).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": q.errMsg})
			return
		}
	}

	uploads, err := readUploads(h.uploadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membaca gambar: " + err.Error()})
		return
	}
	data.Uploads = uploads

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
				"INSERT IGNORE INTO notes (id, title, content, is_private, created_by, updated_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
				n.ID, n.Title, n.Content, n.IsPrivate, n.CreatedBy, n.UpdatedBy, n.CreatedAt, n.UpdatedAt,
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

	if err := writeUploads(h.uploadDir, data.Uploads); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "data berhasil, tapi gagal menyimpan gambar: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "restore berhasil"})
}

// ponytail: semua gambar dimuat ke memori; ganti ke zip streaming kalau total upload sudah ratusan MB
func readUploads(dir string) ([]RawUpload, error) {
	entries, err := os.ReadDir(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var uploads []RawUpload
	for _, e := range entries {
		if !uploadNamePattern.MatchString(e.Name()) {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		uploads = append(uploads, RawUpload{Name: e.Name(), Data: b})
	}
	return uploads, nil
}

// File yang sudah ada dilewati (sama seperti INSERT IGNORE untuk data DB)
func writeUploads(dir string, uploads []RawUpload) error {
	for _, u := range uploads {
		// cegah path traversal & file non-gambar dari backup yang dimodifikasi
		if !uploadNamePattern.MatchString(u.Name) ||
			allowedImageTypes[http.DetectContentType(u.Data)] != filepath.Ext(u.Name) {
			continue
		}
		path := filepath.Join(dir, u.Name)
		if _, err := os.Stat(path); err == nil {
			continue
		}
		if err := os.WriteFile(path, u.Data, 0o644); err != nil {
			return fmt.Errorf("%s: %w", u.Name, err)
		}
	}
	return nil
}
