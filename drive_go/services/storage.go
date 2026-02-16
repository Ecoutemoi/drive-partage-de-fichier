package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type SavedFile struct {
	OriginalName string
	StorageKey   string
	MimeType     string
	SizeBytes    int64
	Path         string
}

func SaveUploadedFile(fileHeader *multipart.FileHeader, uploadDir string) (*SavedFile, error) {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, err
	}

	// Optionnel: limite
	const max = 25 << 20
	if fileHeader.Size > max {
		return nil, fmt.Errorf("fichier trop gros (max %d bytes)", max)
	}

	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if len(ext) > 10 {
		ext = ""
	}

	// Génère un storageKey et évite collision disque
	var storageKey, dstPath string
	for {
		key, err := randomHex(16)
		if err != nil {
			return nil, err
		}
		storageKey = key + ext
		dstPath = filepath.Join(uploadDir, storageKey)

		if _, err := os.Stat(dstPath); os.IsNotExist(err) {
			break
		}
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	size, err := io.Copy(dst, src)
	if err != nil {
		_ = os.Remove(dstPath)
		return nil, err
	}

	mimeType, err := detectMime(dstPath)
	if err != nil {
		_ = os.Remove(dstPath)
		return nil, err
	}

	return &SavedFile{
		OriginalName: fileHeader.Filename,
		StorageKey:   storageKey,
		MimeType:     mimeType,
		SizeBytes:    size,
		Path:         dstPath,
	}, nil
}

func detectMime(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, _ := f.Read(buf)
	return http.DetectContentType(buf[:n]), nil
}

func randomHex(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
