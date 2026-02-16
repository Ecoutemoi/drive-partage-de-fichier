package handlers

import (
	"database/sql"
	"drive-go/db"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	ID            int64     `json:"id"`
	Email         string    `json:"Email"`
	Password_hash string    `json:"-"`
	Full_name     string    `json:"Full_name"`
	CreatedAt     time.Time `json:"created_at"`
}

type signeupRequette struct {
	Email     string `json:"Email"`
	Password  string `json:"Password"`
	Full_name string `json:"Full_name"`
}

var emailRe = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

func Inscription(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSON(w, http.StatusMethodNotAllowed, erreurAPI{Error: "POST uniquement"})
		return
	}

	var req signeupRequette
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "JSON invalide"})
		return
	}

	req.Full_name = strings.TrimSpace(req.Full_name)
	req.Email = strings.TrimSpace(req.Email)

	if req.Email == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "email et password obligatoires"})
		return
	}

	if !emailRe.MatchString(req.Email) {
		writeJSON(w, http.StatusBadRequest, erreurAPI{Error: "Email invalide"})
		return
	}

	var existingID int64
	err := db.DB.QueryRow(`SELECT id FROM users WHERE email = ? LIMIT 1`, req.Email).Scan(&existingID)
	if err == nil {
		writeJSON(w, http.StatusConflict, erreurAPI{Error: "Email déjà utilisé"})
		return
	}

	if !errors.Is(err, sql.ErrNoRows) {
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur base de données"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur sécurité"})
		return
	}

	res, err := db.DB.Exec(
		`INSERT INTO users (email, password_hash, full_name) VALUES (?, ?, ?)`,
		req.Email,
		string(hash),
		nullIfEmpty(req.Full_name),
	)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, erreurAPI{Error: "Erreur insertion"})
		return
	}

	id, _ := res.LastInsertId()

	u := Users{
		ID:        id,
		Email:     req.Email,
		Full_name: req.Full_name,
		CreatedAt: time.Now(),
	}
	writeJSON(w, http.StatusCreated, u)
}

// helper : si full_name est vide, on insère NULL
func nullIfEmpty(s string) any {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}
