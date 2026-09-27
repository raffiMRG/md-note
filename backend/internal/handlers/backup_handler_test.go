package handlers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUploadsRoundTrip(t *testing.T) {
	src, dst := t.TempDir(), t.TempDir()
	png := []byte("\x89PNG\r\n\x1a\n0000")
	name := "0123456789abcdef0123456789abcdef.png"
	os.WriteFile(filepath.Join(src, name), png, 0o644)
	os.WriteFile(filepath.Join(src, "notes.txt"), []byte("x"), 0o644)

	uploads, err := readUploads(src)
	if err != nil || len(uploads) != 1 {
		t.Fatalf("read: %v %v", uploads, err)
	}

	uploads = append(uploads,
		RawUpload{Name: "../evil.png", Data: png},
		RawUpload{Name: "fedcba9876543210fedcba9876543210.png", Data: []byte("<svg/>")},
	)
	if err := writeUploads(dst, uploads); err != nil {
		t.Fatal(err)
	}
	entries, _ := os.ReadDir(dst)
	if len(entries) != 1 || entries[0].Name() != name {
		t.Fatalf("restored: %v", entries)
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(dst), "evil.png")); err == nil {
		t.Fatal("path traversal written")
	}

	if missing, err := readUploads(filepath.Join(src, "nope")); err != nil || missing != nil {
		t.Fatalf("missing dir: %v %v", missing, err)
	}
}
