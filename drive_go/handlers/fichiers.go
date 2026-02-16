package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"drive-go/db"
)

type Fichier struct {
	ID           int64     `json:"id"`
	OriginalName string    `json:"original_name"`
	MimeType     string    `json:"mime_type"`
	SizeBytes    int64     `json:"size_bytes"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
}

type movefichreq struct {
	File_id   int64  `json:"file_id"`
	Folder_id *int64 `json:"folder_id"`
}

type renamefichreq struct {
	FileID int64  `json:"file_id"`
	NvNom  string `json:"new_name"`
}

type delfichreq struct {
	FiledID int64 `json:"file_id"`
	UsersID int64 `json:"user_id"`
}

func Download(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeJSON(w, http.StatusMethodNotAllowed, erreurAPI{Error: "GET uniquement"})
		return
	}

	path := r.URL.Path
	parts := strings.Split(path, "/")
	// attendu: ["", "files", "{id}", "download"]
	if len(parts) < 4 || parts[1] != "files" || parts[3] != "download" {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "URL invalide (attendu: /files/{id}/download)"})
		return
	}

	idStr := parts[2]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "id invalide (int positif requis)"})
		return
	}

	var storageKey, originalName, mimeType string
	row := db.DB.QueryRow(`SELECT storage_key, original_name, mime_type FROM files WHERE id=? AND is_deleted=0`, id)
	if err := row.Scan(&storageKey, &originalName, &mimeType); err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, erreurAPI{Error: "fichier introuvable"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
		return
	}

	uploadDir := getEnv("UPLOAD_DIR", "storage/uploads")
	fullPath := filepath.Join(uploadDir, storageKey)

	file, err := os.Open(fullPath)
	if err != nil {
		writeJSON(w, http.StatusNotFound, erreurAPI{Error: "fichier absent sur le disque"})
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Content-Disposition", `attachment; filename="`+originalName+`"`)
	http.ServeContent(w, r, originalName, time.Now(), file)
}

func List_fichier(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeJSON(w, http.StatusMethodNotAllowed, erreurAPI{Error: "GET uniquement"})
		return
	}

	rows, err := db.DB.Query(`
		SELECT files.id, original_name, mime_type, size_bytes, files.created_at, full_name
		FROM files
		JOIN users ON users.id = user_id
		WHERE files.is_deleted = 0
	`)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
		return
	}
	defer rows.Close()

	var tout_les_fichiers []Fichier
	for rows.Next() {
		var f Fichier
		if err := rows.Scan(
			&f.ID,
			&f.OriginalName,
			&f.MimeType,
			&f.SizeBytes,
			&f.CreatedAt,
			&f.CreatedBy,
		); err != nil {
			writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur lecture base de données"})
			return
		}
		tout_les_fichiers = append(tout_les_fichiers, f)
	}
	if err := rows.Err(); err != nil {
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
		return
	}

	writeJSON(w, http.StatusOK, tout_les_fichiers)
}

func Deplacer_fichier(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
		w.Header().Set("Allow", http.MethodPut)
		writeJSON(w, http.StatusMethodNotAllowed, erreurAPI{Error: "uniquement la method PUT"})
		return
	}

	var req movefichreq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "JSON invalide"})
		return
	}

	if req.File_id <= 0 {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "file_id invalide"})
		return
	}

	if req.Folder_id != nil && *req.Folder_id <= 0 {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "folder_id invalide"})
		return
	}

	var ok int
	err := db.DB.QueryRow(`SELECT 1 FROM files WHERE id=? AND is_deleted=0 LIMIT 1`, req.File_id).Scan(&ok)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, erreurAPI{Error: "fichier introuvable"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
		return
	}

	if req.Folder_id != nil {
		err = db.DB.QueryRow(`SELECT 1 FROM folders WHERE id=? AND is_deleted=0 LIMIT 1`, *req.Folder_id).Scan(&ok)
		if err != nil {
			if err == sql.ErrNoRows {
				writeJSON(w, http.StatusNotFound, erreurAPI{Error: "dossier introuvable"})
				return
			}
			writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
			return
		}
	}

	// UPDATE folder_id = NULL si Folder_id == nil
	var folderAny any = nil
	if req.Folder_id != nil {
		folderAny = *req.Folder_id
	}

	res, err := db.DB.Exec(`UPDATE files SET folder_id=? WHERE id=? AND is_deleted=0`, folderAny, req.File_id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
		return
	}

	aff, _ := res.RowsAffected()
	if aff == 0 {
		writeJSON(w, http.StatusNotFound, erreurAPI{Error: "fichier introuvable ou aucune modification"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"ok":        true,
		"file_id":   req.File_id,
		"folder_id": req.Folder_id, // null si racine
	})
}

func Modifier_name_files(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
		w.Header().Set("Allow", http.MethodPut)
		writeJSON(w, http.StatusMethodNotAllowed, erreurAPI{Error: "uniquement la method PUT"})
		return
	}

	var req renamefichreq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "JSON invalide"})
		return
	}

	req.NvNom = strings.TrimSpace(req.NvNom)
	if req.FileID <= 0 {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "file_id invalide"})
		return
	}

	if req.NvNom == "" || len(req.NvNom) > 255 {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "nom invalide"})
		return
	}

	res, err := db.DB.Exec(
		`UPDATE files SET original_name=? WHERE id=? AND is_deleted=0`,
		req.NvNom,
		req.FileID,
	)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
		return
	}

	aff, _ := res.RowsAffected()
	if aff == 0 {
		writeJSON(w, http.StatusNotFound, erreurAPI{Error: "fichier introuvable"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"ok":       true,
		"file_id":  req.FileID,
		"new_name": req.NvNom,
	})
}

func Delete_files(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		w.Header().Set("Allow", http.MethodDelete)
		writeJSON(w, http.StatusMethodNotAllowed, erreurAPI{Error: "uniquement la method DELETE"})
		return
	}

	var req delfichreq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "JSON invalide"})
		return
	}

	if req.FiledID <= 0 || req.UsersID <= 0 {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "file_id ou user_id invalide"})
		return
	}

	res, err := db.DB.Exec(`
		UPDATE files
		SET is_deleted=1, deleted_at=NOW(), delete_by=?
		WHERE id=? AND is_deleted=0
	`, req.UsersID, req.FiledID)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
		return
	}

	aff, _ := res.RowsAffected()
	if aff == 0 {
		writeJSON(w, http.StatusNotFound, erreurAPI{Error: "fichier introuvable ou déjà supprimé"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"ok":      true,
		"file_id": req.FiledID,
		"user_id": req.UsersID,
	})
}
