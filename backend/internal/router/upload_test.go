package router

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"md-note/backend/internal/auth"
	"md-note/backend/internal/config"
	"md-note/backend/internal/handlers"
)

func TestUpload(t *testing.T) {
	cfg := config.Config{JWTSecret: "s", UploadDir: t.TempDir()}
	r := New(cfg, NewCORSCache(nil), nil, nil, nil, nil, nil, nil, handlers.NewUploadHandler(cfg.UploadDir))
	token, _ := auth.GenerateToken(cfg.JWTSecret, 1, "u", "user", time.Hour)

	upload := func(content []byte, withAuth bool) *httptest.ResponseRecorder {
		var body bytes.Buffer
		w := multipart.NewWriter(&body)
		fw, _ := w.CreateFormFile("file", "x")
		fw.Write(content)
		w.Close()
		req := httptest.NewRequest(http.MethodPost, "/api/uploads", &body)
		req.Header.Set("Content-Type", w.FormDataContentType())
		if withAuth {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		return rec
	}

	png := []byte("\x89PNG\r\n\x1a\n0000")
	if rec := upload(png, false); rec.Code != http.StatusUnauthorized {
		t.Fatalf("no auth: got %d", rec.Code)
	}
	if rec := upload([]byte(`<svg onload="alert(1)"/>`), true); rec.Code != http.StatusBadRequest {
		t.Fatalf("svg: got %d", rec.Code)
	}
	rec := upload(png, true)
	if rec.Code != http.StatusCreated {
		t.Fatalf("png: got %d %s", rec.Code, rec.Body)
	}
	var res struct{ URL string }
	json.Unmarshal(rec.Body.Bytes(), &res)

	get := httptest.NewRecorder()
	r.ServeHTTP(get, httptest.NewRequest(http.MethodGet, res.URL, nil))
	if get.Code != http.StatusOK || get.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("get %s: %d %v", res.URL, get.Code, get.Header())
	}
}
