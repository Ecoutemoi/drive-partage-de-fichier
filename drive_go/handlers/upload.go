package handlers

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"strings"

	"drive-go/db"
	"drive-go/services"
)

// Structure d'erreur API
type erreurAPI struct {
	Error string `json:"error"`
}

// POST /upload
func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSON(w, http.StatusMethodNotAllowed, erreurAPI{Error: "POST uniquement"})
		return
	}

	userID := int64(1) // temporaire

	// Limite 200MB -> 413 si trop gros (MaxBytesReader déclenche ça)
	r.Body = http.MaxBytesReader(w, r.Body, 200<<20)

	if err := r.ParseMultipartForm(200 << 20); err != nil {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "Formulaire invalide"})
		return
	}

	// Lire folder_id (optionnel)
	var folderID any = nil
	if v := strings.TrimSpace(r.FormValue("folder_id")); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil || id <= 0 {
			writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "folder_id invalide"})
			return
		}

		var ok int
		err = db.DB.QueryRow(
			`SELECT 1 FROM folders WHERE id=? AND user_id=? AND is_deleted=0 LIMIT 1`,
			id, userID,
		).Scan(&ok)
		if err != nil {
			if err == sql.ErrNoRows {
				writeJSON(w, http.StatusNotFound, erreurAPI{Error: "Dossier introuvable"})
				return
			}
			writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
			return
		}

		folderID = id
	}

	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "Fichier manquant"})
		return
	}
	fileHeader := files[0]

	uploadDir := getEnv("UPLOAD_DIR", "storage/uploads")
	saved, err := services.SaveUploadedFile(fileHeader, uploadDir)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur sauvegarde fichier"})
		return
	}

	res, err := db.DB.Exec(
		`INSERT INTO files (user_id, folder_id, original_name, storage_key, mime_type, size_bytes)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		userID,
		folderID,
		saved.OriginalName,
		saved.StorageKey,
		saved.MimeType,
		saved.SizeBytes,
	)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
		return
	}

	id, _ := res.LastInsertId()

	writeJSON(w, http.StatusCreated, map[string]any{
		"id":            id,
		"folder_id":     folderID,
		"original_name": saved.OriginalName,
		"storage_key":   saved.StorageKey,
		"mime_type":     saved.MimeType,
		"size_bytes":    saved.SizeBytes,
	})
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
