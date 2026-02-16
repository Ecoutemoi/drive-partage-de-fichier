package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"drive-go/db"
)

type createFolderReq struct {
	ParentID *int64 `json:"parent_id"`
	Name     string `json:"name"`
}

type FolderItem struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	ParentID  *int64    `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
}

type FileItem struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	FolderID     *int64    `json:"folder_id"`
	OriginalName string    `json:"original_name"`
	StorageKey   string    `json:"storage_key"`
	MimeType     string    `json:"mime_type"`
	SizeBytes    int64     `json:"size_bytes"`
	CreatedAt    time.Time `json:"created_at"`
}

type renamefichdoss struct {
	DossID int64  `json:"doss_id"`
	NvNom  string `json:"new_name"`
}

type delFolderReq struct {
	FolderID int64 `json:"folder_id"`
	UsersID  int64 `json:"users_id"`
}

type movefoldreq struct {
	Fold_id   int64  `json:"folder_id"`
	Parent_id *int64 `json:"parent_id"`
}

func nullIfNil(p *int64) any {
	if p == nil {
		return nil
	}
	return *p
}

func CreateFolder(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSON(w, http.StatusMethodNotAllowed, erreurAPI{Error: "POST uniquement"})
		return
	}

	userID := int64(1) // temporaire


	var req createFolderReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "JSON invalide"})
		return
	}


	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" || len(req.Name) > 255 {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "Nom de dossier invalide"})
		return
	}

	
	if req.ParentID != nil {
		var ok int
		err := db.DB.QueryRow(`SELECT 1 FROM folders WHERE id=? LIMIT 1`, *req.ParentID).Scan(&ok)
		if err != nil {
			if err == sql.ErrNoRows {
				writeJSON(w, http.StatusNotFound, erreurAPI{Error: "Dossier parent introuvable"})
				return
			}
			writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
			return
		}
	}


	res, err := db.DB.Exec(
		`INSERT INTO folders (user_id, parent_id, name) VALUES (?, ?, ?)`,
		userID,
		nullIfNil(req.ParentID),
		req.Name,
	)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
		return
	}

	id, _ := res.LastInsertId()


	writeJSON(w, http.StatusCreated, map[string]any{
		"id":        id,
		"user_id":   userID,
		"parent_id": req.ParentID,
		"name":      req.Name,
	})
}

// GET /folders?parent_id=
// - parent_id absent/vide => root
// - parent_id=12 => contenu du dossier 12
func ListDossiers(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeJSON(w, http.StatusMethodNotAllowed, erreurAPI{Error: "GET uniquement"})
		return
	}

	pid := strings.TrimSpace(r.URL.Query().Get("parent_id"))

	// parentID = nil => racine
	var parentID *int64
	if pid != "" {
		v, err := strconv.ParseInt(pid, 10, 64)
		
		if err != nil || v <= 0 {
			writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "parent_id invalide"})
			return
		}
		parentID = &v

	
		
		var ok int
		err = db.DB.QueryRow(`SELECT 1 FROM folders WHERE id=? LIMIT 1`, v).Scan(&ok)
		if err != nil {
			if err == sql.ErrNoRows {
				writeJSON(w, http.StatusNotFound, erreurAPI{Error: "Dossier introuvable"})
				return
			}
			writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
			return
		}
	}


	rowsFolders, err := db.DB.Query(`
		SELECT id, user_id, name, parent_id, created_at
		FROM folders
		WHERE parent_id <=> ?
		ORDER BY name ASC
	`, nullIfNil(parentID))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
		return
	}
	defer rowsFolders.Close()

	var folders []FolderItem
	for rowsFolders.Next() {
		var f FolderItem
		if err := rowsFolders.Scan(&f.ID, &f.UserID, &f.Name, &f.ParentID, &f.CreatedAt); err != nil {
			writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
			return
		}
		folders = append(folders, f)
	}
	if err := rowsFolders.Err(); err != nil {
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
		return
	}


	rowsFiles, err := db.DB.Query(`
		SELECT id, user_id, folder_id, original_name, storage_key, mime_type, size_bytes, created_at
		FROM files
		WHERE folder_id <=> ? AND is_deleted = 0
		ORDER BY created_at DESC
	`, nullIfNil(parentID))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
		return
	}
	defer rowsFiles.Close()

	var files []FileItem
	for rowsFiles.Next() {
		var f FileItem
		if err := rowsFiles.Scan(
			&f.ID,
			&f.UserID,
			&f.FolderID,
			&f.OriginalName,
			&f.StorageKey,
			&f.MimeType,
			&f.SizeBytes,
			&f.CreatedAt,
		); err != nil {
			writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
			return
		}
		files = append(files, f)
	}

	if err := rowsFiles.Err(); err != nil {
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
		return
	}


	writeJSON(w, http.StatusOK, map[string]any{
		"parent_id": parentID, // null si racine
		"folders":   folders,
		"files":     files,
	})
}

func Modifier_name_dossiers(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
		w.Header().Set("Allow", http.MethodPut)
		writeJSON(w, http.StatusMethodNotAllowed, erreurAPI{Error: "uniquement la method PUT"})
		return
	}


	var req renamefichdoss
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "JSON invalide"})
		return
	}

	req.NvNom = strings.TrimSpace(req.NvNom)
	if req.DossID <= 0 {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "doss_id invalide"})
		return
	}

	if req.NvNom == "" || len(req.NvNom) > 255 {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "nom invalide"})
		return
	}

	res, err := db.DB.Exec(`UPDATE folders SET name=? WHERE id=?`, req.NvNom, req.DossID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
		return
	}

	aff, _ := res.RowsAffected()
	if aff == 0 {
		writeJSON(w, http.StatusNotFound, erreurAPI{Error: "dossier introuvable"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"ok":       true,
		"doss_id":  req.DossID,
		"new_name": req.NvNom,
	})
}

func Delete_folder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.Header().Set("Allow", http.MethodDelete)
		writeJSON(w, http.StatusMethodNotAllowed, erreurAPI{Error: "uniquement la method DELETE"})
		return
	}

	var req delFolderReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "JSON invalide"})
		return
	}

	if req.FolderID <= 0 || req.UsersID <= 0 {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "folder_id et users_id invalides"})
		return
	}

	var ok int
	err := db.DB.QueryRow(`SELECT 1 FROM folders WHERE id=? AND is_deleted=0 LIMIT 1`, req.FolderID).Scan(&ok)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, erreurAPI{Error: "folder introuvable ou déjà supprimé"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
		return
	}

	rows, err := db.DB.Query(`
		SELECT 1
		FROM folders
		LEFT JOIN files ON folders.id = files.folder_id AND files.is_deleted=0
		WHERE (folders.parent_id = ? AND folders.is_deleted=0) OR files.folder_id = ?
		LIMIT 1
	`, req.FolderID, req.FolderID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
		return
	}
	defer rows.Close()

	if rows.Next() {
		writeJSON(w, http.StatusConflict, erreurAPI{Error: "le dossier n'est pas vide"})
		return
	}

	res, err := db.DB.Exec(`
		UPDATE folders
		SET is_deleted=1, deleted_at=NOW(), delete_by=?
		WHERE id=? AND is_deleted=0
	`, req.UsersID, req.FolderID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
		return
	}

	aff, _ := res.RowsAffected()
	if aff == 0 {
		// (course condition possible)
		writeJSON(w, http.StatusNotFound, erreurAPI{Error: "folder introuvable ou déjà supprimé"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"ok":        true,
		"folder_id": req.FolderID,
		"users_id":  req.UsersID,
	})
}

func Deplacer_dossiers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.Header().Set("Allow", http.MethodPut)
		writeJSON(w, http.StatusMethodNotAllowed, erreurAPI{Error: "uniquement la method PUT"})
		return
	}

	var req movefoldreq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "JSON invalide"})
		return
	}

	if req.Fold_id <= 0 {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "folder_id invalide"})
		return
	}
	if req.Parent_id != nil && *req.Parent_id <= 0 {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "parent_id invalide"})
		return
	}

	var ok int
	err := db.DB.QueryRow(`SELECT 1 FROM folders WHERE id=? AND is_deleted=0 LIMIT 1`, req.Fold_id).Scan(&ok)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, erreurAPI{Error: "dossier introuvable"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
		return
	}

	if req.Parent_id != nil {
		err = db.DB.QueryRow(`SELECT 1 FROM folders WHERE id=? AND is_deleted=0 LIMIT 1`, *req.Parent_id).Scan(&ok)
		if err != nil {
			if err == sql.ErrNoRows {
				writeJSON(w, http.StatusNotFound, erreurAPI{Error: "dossier parent introuvable"})
				return
			}
			writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
			return
		}
	}

	if req.Parent_id != nil && *req.Parent_id == req.Fold_id {
		writeJSON(w, http.StatusConflict, erreurAPI{Error: "impossible de déplacer un dossier dans lui-même"})
		return
	}

	// UPDATE : parent_id = NULL si req.Parent_id == nil
	var parentAny any = nil
	if req.Parent_id != nil {
		parentAny = *req.Parent_id
	}

	res, err := db.DB.Exec(`UPDATE folders SET parent_id=? WHERE id=?`, parentAny, req.Fold_id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
		return
	}

	aff, _ := res.RowsAffected()
	if aff == 0 {
		// cas rare : même valeur, ou course condition
		writeJSON(w, http.StatusNotFound, erreurAPI{Error: "aucune modification"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"ok":        true,
		"folder_id": req.Fold_id,
		"parent_id": req.Parent_id, // null si racine
	})
}
