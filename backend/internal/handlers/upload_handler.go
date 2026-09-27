package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const maxUploadSize = 5 << 20

// SVG sengaja tidak diizinkan: bisa berisi script (XSS)
var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

type UploadHandler struct {
	dir string
}

func NewUploadHandler(dir string) *UploadHandler {
	return &UploadHandler{dir: dir}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file wajib diisi"})
		return
	}
	if file.Size > maxUploadSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ukuran gambar maksimal 5 MB"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gagal membaca file"})
		return
	}
	head := make([]byte, 512)
	n, _ := f.Read(head)
	f.Close()

	ext, ok := allowedImageTypes[http.DetectContentType(head[:n])]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "hanya gambar jpg, png, gif, atau webp"})
		return
	}

	b := make([]byte, 16)
	rand.Read(b)
	name := hex.EncodeToString(b) + ext

	if err := c.SaveUploadedFile(file, filepath.Join(h.dir, name)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan gambar"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"url": "/api/uploads/" + name})
}
