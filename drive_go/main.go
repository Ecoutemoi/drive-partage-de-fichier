package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"drive-go/db"
	"drive-go/handlers"
)

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Autorise uniquement les origines de dev (Flutter Web)
		if origin != "" && (strings.HasPrefix(origin, "http://localhost:") || strings.HasPrefix(origin, "http://127.0.0.1:")) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent) // 204
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	_ = godotenv.Load()

	_, err := db.Init()
	if err != nil {
		log.Fatal("DB ERROR:", err)
	}
	defer db.Close()

	// Routes simples
	http.HandleFunc("/upload", handlers.Upload)
	http.HandleFunc("/files", handlers.List_fichier)
	http.HandleFunc("/auth/Inscription", handlers.Inscription)

	http.HandleFunc("/Createfolders", handlers.CreateFolder)
	http.HandleFunc("/folders/list", handlers.ListDossiers)
	http.HandleFunc("/files/move", handlers.Deplacer_fichier)
	http.HandleFunc("/files/rename", handlers.Modifier_name_files)
	http.HandleFunc("/folders/rename", handlers.Modifier_name_dossiers)
	http.HandleFunc("/files/delete", handlers.Delete_files)
	http.HandleFunc("/folders/delete", handlers.Delete_folder)
	http.HandleFunc("/folders/move", handlers.Deplacer_dossiers)

	// ✅ IMPORTANT: net/http ne supporte pas /files/{id}/download en pattern
	// On utilise un handler "prefix" /files/ qui dispatch vers Download ou Share.
	http.HandleFunc("/files/", func(w http.ResponseWriter, r *http.Request) {
		// Exemples attendus:
		// /files/123/download
		// /files/123/share
		if strings.HasSuffix(r.URL.Path, "/download") {
			handlers.Download(w, r)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/share") {
			handlers.CreateShareLink(w, r)
			return
		}
		http.NotFound(w, r)
	})

	// /share/{token} (ton handler parse déjà le path)
	http.HandleFunc("/share/", handlers.UseShareLink)

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8010"
	}

	log.Println("Server started on :" + port)

	// Enveloppe avec CORS
	handler := withCORS(http.DefaultServeMux)

	log.Fatal(http.ListenAndServe(":"+port, handler))
}
