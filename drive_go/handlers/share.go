package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"drive-go/db"
	"drive-go/utils"
)

func CreateShareLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSON(w, http.StatusMethodNotAllowed, erreurAPI{Error: "POST uniquement"})
		return
	}

	// URL attendue: /files/{id}/share
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 || parts[1] != "files" || parts[3] != "share" {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "URL invalide (attendu: /files/{id}/share)"})
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "id invalide"})
		return
	}

	var exists int
	err = db.DB.QueryRow(`SELECT 1 FROM files WHERE id=? AND is_deleted=0 LIMIT 1`, id).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, erreurAPI{Error: "fichier introuvable"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
		return
	}

	// Génération token (max 5 tentatives)
	var token string
	for i := 0; i < 5; i++ {
		t, err := utils.NewShareToken(24)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "erreur génération token"})
			return
		}
		token = t

		_, err = db.DB.Exec(`INSERT INTO share_links (file_id, token) VALUES (?, ?)`, id, token)
		if err == nil {
			break
		}

		if i == 4 {
			writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "impossible de créer un lien"})
			return
		}
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	url := scheme + "://" + r.Host + "/share/" + token

	writeJSON(w, http.StatusCreated, map[string]any{
		"token": token,
		"url":   url,
	})
}

func UseShareLink(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeJSON(w, http.StatusMethodNotAllowed, erreurAPI{Error: "GET uniquement"})
		return
	}

	// URL attendue: /share/{token}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[1] != "share" || parts[2] == "" {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "URL invalide (attendu: /share/{token})"})
		return
	}
	token := parts[2]

	var fileID int64
	var expiresAt sql.NullTime
	var isActive int
	err := db.DB.QueryRow(
		`SELECT file_id, expires_at, is_active FROM share_links WHERE token=? LIMIT 1`,
		token,
	).Scan(&fileID, &expiresAt, &isActive)

	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, erreurAPI{Error: "lien invalide"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
		return
	}

	if isActive != 1 {
		writeJSON(w, http.StatusGone, erreurAPI{Error: "lien désactivé"})
		return
	}

	if expiresAt.Valid && time.Now().After(expiresAt.Time) {
		writeJSON(w, http.StatusGone, erreurAPI{Error: "lien expiré"})
		return
	}

	var storageKey, originalName, mimeType string
	err = db.DB.QueryRow(
		`SELECT storage_key, original_name, mime_type FROM files WHERE id=? AND is_deleted=0 LIMIT 1`,
		fileID,
	).Scan(&storageKey, &originalName, &mimeType)

	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, erreurAPI{Error: "fichier introuvable"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
		return
	}

	uploadDir := getEnv("UPLOAD_DIR", "storage/uploads")
	fullPath := filepath.Join(uploadDir, storageKey)

	f, err := os.Open(fullPath)
	if err != nil {
		writeJSON(w, http.StatusNotFound, erreurAPI{Error: "fichier absent sur disque"})
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Content-Disposition", `attachment; filename="`+originalName+`"`)
	http.ServeContent(w, r, originalName, time.Now(), f)

	// Désactivation après usage (one-time link)
	_, err = db.DB.Exec(`UPDATE share_links SET is_active=0 WHERE token=? LIMIT 1`, token)
	if err != nil {
		log.Println("Erreur désactivation lien:", err)
	}
}
